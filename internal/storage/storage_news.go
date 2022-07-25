package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/AscaroLabs/go-news/internal/config"
	pb "github.com/AscaroLabs/go-news/internal/proto"
	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Offset int64 `protobuf:"varint,1,opt,name=offset,proto3" json:"offset,omitempty"`
// 	// Количество возвращаемых объектов на странице
// 	Limit int64 `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
// 	// Поле для сортировки (active_from, date_create)
// 	Sort string `protobuf:"bytes,3,opt,name=sort,proto3" json:"sort,omitempty"`
// 	// Направление сортировки (desc, asc)
// 	Order string `protobuf:"bytes,4,opt,name=order,proto3" json:"order,omitempty"`
// 	// Поиск по строке
// 	Query string `protobuf:"bytes,5,opt,name=query,proto3" json:"query,omitempty"`
// 	// Параметры фильтрации
// 	Filter *ListRequestFilter

//  create table News (
//     id text primary key,
//     title text not null,
//     author text not null,
//     active boolean not null,
//     activeFrom integer not null,
//     text text not null,
//     textJSON text not null,
//     userId text references Users(id),
//     isImportant boolean not null
// );

// select * from News where ....(fiter string)....
// order by ...(Sort)... ...(Order)... limit ...(Limit)... offset ...(Offset)...

func GetNews(ctx context.Context, r *pb.NewsRequestParams) ([]*pb.NewsObject, error) {

	cfg, ok := config.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("can't read config")
	}
	pool, err := getPool(cfg)
	if err != nil {
		return nil, err
	}
	defer pool.Close()

	log.Print("pool created")

	// Остро стоит проблема инъекций, нужно подрубать ручную валидацию
	var filterQuery string
	if r.GetFilter() != nil {
		filterQuery = func(f *pb.ListRequestFilter) string {
			var active string
			switch f.Mode {
			case "active":
				active = "true"
			case "inactive":
				active = "false"
			default:
				active = ""
			}
			var res string
			switch {
			case len(f.UserId) > 0 && len(f.Mode) > 0:
				res = fmt.Sprintf("where userId=%s and active=%s", f.UserId, active)
			case len(f.UserId) > 0:
				res = fmt.Sprintf("where userId=%s", f.UserId)
			case len(f.Mode) > 0:
				res = fmt.Sprintf("where active=%s", active)
			default:
				res = ""
			}
			return res
		}(r.GetFilter())
	} else {
		filterQuery = ""
	}

	log.Printf("filter query: %s", filterQuery)

	var sortColumn string
	switch r.Sort {
	case "active_from":
		sortColumn = "activeFrom"
	case "date_create":
		sortColumn = ""
	default:
		sortColumn = ""
	}

	q := fmt.Sprintf(
		`
		select id,title,author,active,activeFrom,text,textJSON,userId,isImportant,tags,files from news
		%s order by %s %s limit %d offset %d`,
		filterQuery,
		sortColumn,
		r.GetOrder(),
		r.GetLimit(),
		r.GetOffset(),
	)

	log.Printf("query: %s", q)

	rows, err := pool.Query(context.Background(), q)
	if err != nil {
		log.Printf("ooops")
		return nil, err
	}
	defer rows.Close()

	res := make([]*pb.NewsObject, 0)
	for rows.Next() {
		log.Print("ok, hanle rows")
		n, err := newsObjectFromRow(pool, rows)
		if err != nil {
			return nil, err
		}
		// var n pb.NewsObject
		// var tags []*pb.Tag
		// var files []*pb.File
		// tags_ids := make([]string, 0)
		// files_ids := make([]string, 0)
		// if err := rows.Scan(
		// 	&n.Id,
		// 	&n.Title,
		// 	&n.Author,
		// 	&n.Active,
		// 	&n.ActiveFrom,
		// 	&n.Text,
		// 	&n.TextJson,
		// 	&n.UserId,
		// 	&n.IsImportant,
		// 	&tags_ids,
		// 	&files_ids,
		// ); err != nil {
		// 	return nil, err
		// }
		// for _, tag_id := range tags_ids {
		// 	var tag pb.Tag
		// 	if err := pool.QueryRow(
		// 		context.Background(),
		// 		`select id,name from Tags where id=$1`, tag_id,
		// 	).Scan(&tag.Id, &tag.Name); err != nil {
		// 		return nil, err
		// 	}
		// 	tags = append(tags, &tag)
		// }
		// for _, file_id := range files_ids {
		// 	var file pb.File
		// 	if err := pool.QueryRow(
		// 		context.Background(),
		// 		`select id,name,ext,base64,dateCreate,userId from Files where id=$1`, file_id,
		// 	).Scan(&file.Id, &file.Name, &file.Ext, &file.Base64, &file.DateCreate, &file.UserId); err != nil {
		// 		return nil, err
		// 	}
		// 	files = append(files, &file)
		// }
		// n.Tags = append(n.Tags, tags...)
		// n.FilesInfo = append(n.FilesInfo, files...)
		res = append(res, n)
	}

	log.Print("nice, let's return result")

	return res, nil
}

func GetOne(ctx context.Context, r *pb.ObjectId) (*pb.NewsObject, error) {
	cfg, ok := config.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("can't read config")
	}
	pool, err := getPool(cfg)
	if err != nil {
		return nil, err
	}
	defer pool.Close()

	q := `
		select id,title,author,active,activeFrom,text,textJSON,userId,isImportant,tags,files 
		from news where id=$1
		 `

	row := pool.QueryRow(context.Background(), q, r.Id)
	return newsObjectFromRow(pool, row)
}

func CreateNewsTxn(ctx context.Context, r *pb.RequestNewsObject, tknDTO *TokenDTO) (bool, error) {
	cfg, ok := config.FromContext(ctx)
	if !ok {
		return ok, fmt.Errorf("can't read config")
	}
	pool, err := getPool(cfg)
	if err != nil {
		return false, err
	}
	defer pool.Close()

	log.Print("start txn!")
	tx, err := pool.Begin(context.Background())
	if err != nil {
		return false, err
	}
	defer tx.Rollback(context.Background())

	// Сначала создадим теги если нужны и файлы
	q := `
		insert into News (id,title,author,active,activeFrom,text,textJSON,userId,isImportant,tags,files) 
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`

	tag_ids := make([]string, 0)
	for _, tag := range r.Tags {
		tag_ids = append(tag_ids, tag.Id)
	}
	files_ids := make([]string, 0)
	for _, file := range r.Files {
		tag_ids = append(files_ids, file.Id)
	}

	_, err = tx.Exec(
		context.Background(),
		q,
		uuid.NewString(),
		r.Title,
		tknDTO.Name,
		r.Active,
		r.ActiveFrom,
		r.Text,
		r.TextJson,
		tknDTO.UserId,
		r.IsImportant,
		tag_ids,
		files_ids,
	)
	if err != nil {
		return false, err
	}

	log.Printf("Query %s DONE", q)

	err = tx.Commit(context.Background())
	if err != nil {
		return false, err
	}

	return true, nil
}

func newsObjectFromRow(pool *pgxpool.Pool, row pgx.Row) (*pb.NewsObject, error) {
	log.Print("takes object from row")
	var n pb.NewsObject
	var tags []*pb.Tag
	var files []*pb.File
	tags_ids := make([]string, 0)
	files_ids := make([]string, 0)
	if err := row.Scan(
		&n.Id,
		&n.Title,
		&n.Author,
		&n.Active,
		&n.ActiveFrom,
		&n.Text,
		&n.TextJson,
		&n.UserId,
		&n.IsImportant,
		&tags_ids,
		&files_ids,
	); err != nil {
		log.Printf("little problem")
		return nil, err
	}

	for i, tag_id := range tags_ids {
		log.Printf("%d : %s", i, tag_id)
	}

	for _, tag_id := range tags_ids {
		var tag pb.Tag
		if err := pool.QueryRow(
			context.Background(),
			`select id,name from Tags where id=$1`, tag_id,
		).Scan(&tag.Id, &tag.Name); err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}
	log.Printf("takes files: %v", files_ids)
	for _, file_id := range files_ids {
		var file pb.File
		if err := pool.QueryRow(
			context.Background(),
			`select id,name,ext,base64,dateCreate,userId from Files where id=$1`, file_id,
		).Scan(&file.Id, &file.Name, &file.Ext, &file.Base64, &file.DateCreate, &file.UserId); err != nil {
			return nil, err
		}
		files = append(files, &file)
	}
	n.Tags = append(n.Tags, tags...)
	n.FilesInfo = append(n.FilesInfo, files...)
	return &n, nil
}
