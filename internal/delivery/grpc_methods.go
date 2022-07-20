// Имплементируем методы gRPC севисов
package delivery

import (
	"context"
	"math/rand"

	pb "github.com/AscaroLabs/go-news/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s *contentCheckServiceServer) CheckHealth(ctx context.Context, in *pb.EmptyRequest) (*pb.HealthResponse, error) {
	if service_alive() {
		_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "200"))
		return &pb.HealthResponse{
			ServiceName:   "ContentCheckService",
			ServiceStatus: "200",
		}, nil
	} else {
		_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "500"))
		return &pb.HealthResponse{
			ServiceStatus: "500",
		}, nil
	}
}

func service_alive() bool {
	return (rand.Intn(100) < 50)
}
