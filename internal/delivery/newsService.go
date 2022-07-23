package delivery

import (
	pb "github.com/AscaroLabs/go-news/internal/proto"
)

type newsServiceServer struct {
	pb.UnimplementedNewsServiceServer
}
