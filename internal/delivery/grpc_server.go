package delivery

import (
	"github.com/AscaroLabs/go-news/internal/auth"
	"github.com/AscaroLabs/go-news/internal/config"
	pb "github.com/AscaroLabs/go-news/internal/proto"
	"google.golang.org/grpc"
)

// функция создает новый gRPC сервер
func NewGRPCServer(cfg *config.Config, tm *auth.TokenManager) (*grpc.Server, error) {
	grpc_server := grpc.NewServer()
	pb.RegisterContentCheckServiceServer(grpc_server, &contentCheckServiceServer{})
	pb.RegisterNewsServiceServer(grpc_server, &newsServiceServer{cfg: cfg, tm: tm})
	pb.RegisterTagServiceServer(grpc_server, &tagServiceServer{})
	return grpc_server, nil
}
