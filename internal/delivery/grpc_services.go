package delivery

import (
	"github.com/AscaroLabs/go-news/internal/config"
	pb "github.com/AscaroLabs/go-news/internal/proto"
	"google.golang.org/grpc"
)

type contentCheckServiceServer struct {
	pb.UnimplementedContentCheckServiceServer
}

type newsServiceServer struct {
	pb.UnimplementedNewsServiceServer
}

type tagServiceServer struct {
	pb.UnimplementedTagServiceServer
}

func NewGRPCServer(cfg *config.Config) (*grpc.Server, error) {
	grpc_server := grpc.NewServer()
	pb.RegisterContentCheckServiceServer(grpc_server, &contentCheckServiceServer{})
	pb.RegisterNewsServiceServer(grpc_server, &newsServiceServer{})
	pb.RegisterTagServiceServer(grpc_server, &tagServiceServer{})
	return grpc_server, nil
}
