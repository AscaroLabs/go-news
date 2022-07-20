package delivery

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/AscaroLabs/go-news/internal/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
)

func NewWrapperMux(cfg *config.Config) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(httpResponseModifier),
	)
	return mux, nil
}

func httpResponseModifier(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
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
