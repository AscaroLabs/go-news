package delivery

import (
	"context"

	"github.com/AscaroLabs/go-news/internal/config"
	pb "github.com/AscaroLabs/go-news/internal/proto"
	"github.com/AscaroLabs/go-news/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type newsServiceServer struct {
	cfg *config.Config
	pb.UnimplementedNewsServiceServer
}

func (ns *newsServiceServer) GetNews(ctx context.Context, r *pb.NewsRequestParams) (*pb.NewsList, error) {
	news_list, err := storage.GetNews(config.NewContext(ctx, ns.cfg), r)
	if err != nil {
		return nil, err
	}
	return &pb.NewsList{
		News:  news_list,
		Total: int32(len(news_list)),
	}, nil
}
func (newsServiceServer) GetOne(context.Context, *pb.ObjectId) (*pb.NewsObject, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOne not implemented")
}
func (newsServiceServer) GetOneBySlug(context.Context, *pb.ObjectSlug) (*pb.NewsObject, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOneBySlug not implemented")
}
func (newsServiceServer) Create(context.Context, *pb.RequestNewsObject) (*pb.BaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (newsServiceServer) Update(context.Context, *pb.RequestNewsObject) (*pb.BaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (newsServiceServer) Delete(context.Context, *pb.ObjectId) (*pb.BaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (newsServiceServer) GetFileLink(context.Context, *pb.FileId) (*pb.FileLink, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFileLink not implemented")
}
