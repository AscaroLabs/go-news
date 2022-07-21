package delivery

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/AscaroLabs/go-news/internal/config"
	pb "github.com/AscaroLabs/go-news/internal/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

type RESTServer struct {
	mux *runtime.ServeMux
}

func NewRESTServer(cfg *config.Config) (*RESTServer, error) {
	mux, err := NewWrapperMux(cfg)
	if err != nil {
		return nil, err
	}
	return &RESTServer{
		mux: mux,
	}, nil
}

func (restServer *RESTServer) Run(cfg *config.Config) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	pb.RegisterContentCheckServiceHandlerFromEndpoint(
		ctx,
		restServer.mux,
		fmt.Sprintf("%s:%s", cfg.GetGRPCHost(), cfg.GetGRPCPort()),
		opts,
	)
	if err := http.ListenAndServe(
		fmt.Sprintf(":%s", cfg.GetRESTPort()),
		restServer.mux,
	); err != nil {
		return err
	}
	return nil
}

func NewWrapperMux(cfg *config.Config) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(httpResponseStatusCodeModifier),
	)
	return mux, nil
}

func httpResponseStatusCodeModifier(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		log.Fatal("[REST] Can't get metadata from context")
	}
	if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			return err
		}
		delete(md.HeaderMD, "x-http-code")
		delete(w.Header(), "Grpc-Metadata-X-Http-Code")
		w.WriteHeader(code)
	}
	return nil
}
