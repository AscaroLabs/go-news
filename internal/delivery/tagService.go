package delivery

import (
	"context"

	pb "github.com/AscaroLabs/go-news/internal/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type tagServiceServer struct {
	pb.UnimplementedTagServiceServer
}

func (tagServiceServer) Get(context.Context, *pb.EmptyRequest) (*pb.TagList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
