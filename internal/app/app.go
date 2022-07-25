// по конфигу добавляем части с gRPC и HTTP серверами
package app

import (
	"fmt"
	"log"
	"net"

	"github.com/AscaroLabs/go-news/internal/auth"
	"github.com/AscaroLabs/go-news/internal/config"
	"github.com/AscaroLabs/go-news/internal/delivery"
	"google.golang.org/grpc"
)

// Структура отвечающая за все приложение
type App struct {
	cfg        *config.Config
	restServer *delivery.RESTServer
	grpcServer *grpc.Server
}

func NewApp(cfg *config.Config) (*App, error) {
	tm, err := auth.NewTokenManager(cfg)
	if err != nil {
		log.Fatalf("[Auth] Can't create new token manager: %v", err)
	}
	restServer, err := delivery.NewRESTServer(cfg, tm)
	if err != nil {
		log.Fatalf("[REST] Can't create new REST server: %v", err)
	}
	grpcServer, err := delivery.NewGRPCServer(cfg, tm)
	if err != nil {
		log.Fatalf("[gRPC] Can't create new gRPC server: %v", err)
	}
	return &App{
		cfg:        cfg,
		restServer: restServer,
		grpcServer: grpcServer,
	}, nil
}

// в отдельных горутинах запускаем gRPC сервер и REST шлюз
func (application *App) Run() {
	go func() {
		log.Printf(
			"[REST] Server listening at %s:%s",
			application.cfg.GetRESTHost(),
			application.cfg.GetRESTPort(),
		)
		application.restServer.Run(application.cfg)
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
		if err := application.grpcServer.Serve(lis); err != nil {
			log.Fatalf("[gRPC] Failed to serve: %v", err)
		}
	}()
}
