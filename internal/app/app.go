// по конфигу добавляем части с gRPC и HTTP серверами
package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	pb "github.com/AscaroLabs/go-news/internal/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/AscaroLabs/go-news/internal/config"
	"github.com/AscaroLabs/go-news/internal/delivery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	cfg         *config.Config
	mux         *runtime.ServeMux
	grpc_server *grpc.Server
}

func NewApp(cfg *config.Config) (*App, error) {
	mux, err := delivery.NewWrapperMux(cfg)
	if err != nil {
		log.Fatalf("[REST] Can't create new mux: %v", err)
	}
	grpc_server, err := delivery.NewGRPCServer(cfg)
	if err != nil {
		log.Fatalf("[gRPC] Can't create new gRPC server: %v", err)
	}
	return &App{
		cfg:         cfg,
		mux:         mux,
		grpc_server: grpc_server,
	}, nil
}

func (application *App) Run() {
	go func() {
		log.Printf(
			"[REST] Server listening at %s:%s",
			application.cfg.GetRESTHost(),
			application.cfg.GetRESTPort(),
		)
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		pb.RegisterContentCheckServiceHandlerFromEndpoint(
			ctx,
			application.mux,
			fmt.Sprintf("%s:%s", application.cfg.GetGRPCHost(), application.cfg.GetGRPCPort()),
			opts,
		)
		if err := http.ListenAndServe(
			fmt.Sprintf(":%s", application.cfg.GetRESTPort()),
			application.mux,
		); err != nil {
			log.Fatalf("[REST] Can't start REST API server: %v", err)
		}
	}()
	go func() {
		lis, err := net.Listen(
			"tcp",
			fmt.Sprintf(":%s", application.cfg.GetGRPCPort()),
		)
		if err != nil {
			log.Fatalf("[gRPC] Failed to listen: %v", err)
		}
		log.Printf("[gRPC] Server listening at %v", lis.Addr())
		if err := application.grpc_server.Serve(lis); err != nil {
			log.Fatalf("[gRPC] Failed to serve: %v", err)
		}
	}()
}
