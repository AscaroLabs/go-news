package delivery

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/AscaroLabs/go-news/internal/auth"
	"github.com/AscaroLabs/go-news/internal/config"
	pb "github.com/AscaroLabs/go-news/internal/proto"
	"github.com/AscaroLabs/go-news/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type newsServiceServer struct {
	cfg *config.Config
	tm  *auth.TokenManager
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

func (ns *newsServiceServer) GetOne(ctx context.Context, r *pb.ObjectId) (*pb.NewsObject, error) {
	return storage.GetOne(config.NewContext(ctx, ns.cfg), r)
}

func (newsServiceServer) GetOneBySlug(context.Context, *pb.ObjectSlug) (*pb.NewsObject, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOneBySlug not implemented")
}

func (ns *newsServiceServer) Create(ctx context.Context, r *pb.RequestNewsObject) (*pb.BaseResponse, error) {

	log.Printf("creating news!")

	var tknDTO *storage.TokenDTO
	if md, ok := metadata.FromIncomingContext(ctx); ok {

		// log.Printf("get metadata from context %v", md)

		if auth_header_data, ok := md["authorization"]; ok {

			auth_str := auth_header_data[0]

			// log.Printf("get authorization header %v", auth_str)

			if len(auth_str) > 0 && strings.HasPrefix(auth_str, "Bearer") {
				token := strings.TrimSpace(strings.TrimPrefix(auth_str, "Bearer"))

				log.Printf("Token: %s", token)

				tokenDTO, err := ns.tm.ParseToken(token)

				if err != nil {
					log.Print(err.Error())
				} else {
					log.Printf("token parsed %s, %s, %s", tokenDTO.Name, tokenDTO.Role, tokenDTO.UserId)
				}

				if err != nil {
					_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "401"))
					return &pb.BaseResponse{
						Success: false,
						Message: fmt.Sprintf("Rly? %s", err.Error()),
					}, nil
				}
				if tokenDTO.Role != "dealer" {
					return &pb.BaseResponse{
						Success: false,
						Message: "Permission denied",
					}, nil
				}
				tknDTO = &tokenDTO
			} else {
				_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "401"))
				return &pb.BaseResponse{
					Success: false,
					Message: "Wrong Auth method",
				}, nil
			}
		} else {

			log.Printf("not any authorization header")

			_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "401"))
			return &pb.BaseResponse{
				Success: false,
				Message: "Unauthorized",
			}, nil
		}
	}

	log.Printf("now we have tokenDTO: %v", *tknDTO)

	ok, err := storage.CreateNewsTxn(config.NewContext(ctx, ns.cfg), r, tknDTO)
	if err != nil {
		return &pb.BaseResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	log.Printf("news created!")
	return &pb.BaseResponse{
		Success: ok,
		Message: "",
	}, nil
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
