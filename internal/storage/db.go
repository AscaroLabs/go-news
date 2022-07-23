package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AscaroLabs/go-news/internal/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

func getUrl(cfg *config.Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.GetDBUsername(),
		cfg.GetDBPassword(),
		cfg.GetDBHost(),
		cfg.GetDBHostPort(),
		cfg.GetDBName(),
	)
}

// проверяем жизнеспособнось БД
func CheckHealth(cfg *config.Config) (bool, error) {
	pool, err := pgxpool.Connect(context.Background(), getUrl(cfg))

	log.Printf("pool created!\n")

	if err != nil {
		return false, err
	}
	defer pool.Close()

	var info string
	err = pool.QueryRow(context.Background(), "select 'ok'").Scan(&info)
	if err != nil {
		return false, err
	}
	return true, nil
}

// выполняет транзакцию, в которой заносит в БД информацию о новом пользователе
func MakeRegistrationTxn(cfg *config.Config, userDTO UserDTO) error {
	pool, err := pgxpool.Connect(context.Background(), getUrl(cfg))
	if err != nil {
		return err
	}
	defer pool.Close()
	tx, err := pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	q := `
		insert into Users (id,name,email,passwordHash,role,registeredFrom) 
		values ($1,$2,$3,$4,$5,$6)
	`
	_, err = tx.Exec(
		context.Background(),
		q,
		uuid.NewString(),
		userDTO.Name,
		userDTO.Email,
		userDTO.Password,
		userDTO.Role,
		time.Now().Unix(),
	)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func MakeNewSession(cfg *config.Config, userId string, refreshToken string) error {
	pool, err := pgxpool.Connect(context.Background(), getUrl(cfg))
	if err != nil {
		return err
	}
	defer pool.Close()
	tx, err := pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	q := `
		insert into RefreshSessions (userId,refreshToken,expiresIn,createdAt) 
		values ($1,$2,$3,$4)
	`

	createdAt := time.Now()
	refresh_ttl, _ := time.ParseDuration(cfg.GetRefreshTTL())

	_, err = tx.Exec(
		context.Background(),
		q,
		userId,
		refreshToken,
		createdAt.Add(refresh_ttl).Unix(),
		createdAt.Unix(),
	)

	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

// по данным о юзере находим информацию необходимую для создания JWT токена
func GetTokenDTOFromUserDTO(cfg *config.Config, userDTO *UserDTO) (*TokenDTO, error) {

	log.Printf("userDTO for user %s\n, connection string: %s", userDTO.Name, getUrl(cfg))
	pool, err := pgxpool.Connect(context.Background(), getUrl(cfg))

	log.Printf("uh no pool problems...\n")

	if err != nil {
		return nil, err
	}
	defer pool.Close()

	var tokenDTO TokenDTO

	log.Printf("going to the DB and scan data for tokenDTO\n")

	err = pool.QueryRow(
		context.Background(),
		"select id, role from Users where email=$1",
		userDTO.Email,
	).Scan(&tokenDTO.UserId, &tokenDTO.Role)
	if err != nil {
		return nil, err
	}

	log.Printf("nice.....\n")

	return &tokenDTO, nil
}

func GetUserDTObyEmail(cfg *config.Config, email string) (*UserDTO, error) {

	log.Printf("userDTO for email %s\n, connection string: %s", email, getUrl(cfg))
	pool, err := pgxpool.Connect(context.Background(), getUrl(cfg))

	log.Printf("uh no pool problems...\n")

	if err != nil {
		return nil, err
	}
	defer pool.Close()

	var userDTO UserDTO

	log.Printf("going to the DB and scan data for tokenDTO\n")

	err = pool.QueryRow(
		context.Background(),
		"select name,email,passwordHash,role from Users where email=$1",
		email,
	).Scan(&userDTO.Name, &userDTO.Email, &userDTO.Password, &userDTO.Role)
	if err != nil {
		return nil, err
	}

	log.Printf("nice.....\n")

	return &userDTO, nil
}
