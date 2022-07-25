package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/AscaroLabs/go-news/internal/auth"
	"github.com/AscaroLabs/go-news/internal/config"
	pb "github.com/AscaroLabs/go-news/internal/proto"
	"github.com/AscaroLabs/go-news/internal/storage"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

type RESTServer struct {
	mux *runtime.ServeMux
}

func NewRESTServer(cfg *config.Config, tm *auth.TokenManager) (*RESTServer, error) {
	mux, err := NewWrapperMux(cfg, tm)
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
	pb.RegisterNewsServiceHandlerFromEndpoint(
		ctx,
		restServer.mux,
		fmt.Sprintf("%s:%s", cfg.GetGRPCHost(), cfg.GetGRPCPort()),
		opts,
	)
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

func NewWrapperMux(cfg *config.Config, tm *auth.TokenManager) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(httpResponseStatusCodeModifier),
		runtime.WithIncomingHeaderMatcher(CustomMatcher),
	)
	mux.HandlePath("POST", "/signup", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var userDTO storage.UserDTO
		userDTO_json, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Can't read request body: %v", err)))
			return
		}
		err = json.Unmarshal(userDTO_json, &userDTO)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Can't unmarshal body: %v", err)))
			return
		}
		tokens, err := auth.RegisterUser(cfg, tm, &userDTO)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Can't register user: %v", err)))
			return
		}
		tokensJSON, _ := json.Marshal(tokens)
		refresh_ttl, _ := time.ParseDuration(cfg.GetRefreshTTL())
		refreshCookie := &http.Cookie{
			Name:     "refresh_token",
			Value:    tokens.RefreshToken,
			MaxAge:   int(refresh_ttl.Seconds()),
			HttpOnly: true,
		}
		http.SetCookie(w, refreshCookie)
		w.Write(tokensJSON)
	})
	mux.HandlePath("GET", "/signin", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		email, password, ok := r.BasicAuth()
		if !ok {

			log.Printf("basic auth went wrong!")

			// Я в шоке, если эти две штуки поменять местами, то все сломается
			w.Header().Set("WWW-Authenticate", "Basic")
			w.WriteHeader(http.StatusUnauthorized)

			w.Write([]byte("Auth went wrong!"))
			return
		}
		tokens, err := auth.SignIn(cfg, tm, &storage.SignInDTO{
			Email:    email,
			Password: password,
		})
		if err != nil {
			w.Header().Set("WWW-Authenticate", "Basic")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("Auth went wrong: %v", err)))
			return
		}
		tokensJSON, _ := json.Marshal(tokens)
		refresh_ttl, _ := time.ParseDuration(cfg.GetRefreshTTL())
		refreshCookie := &http.Cookie{
			Name:   "refresh_token",
			Value:  tokens.RefreshToken,
			MaxAge: int(refresh_ttl.Seconds()),
		}
		http.SetCookie(w, refreshCookie)
		w.Write(tokensJSON)
	})
	return mux, nil
}

func CustomMatcher(key string) (string, bool) {
	switch key {
	case "Authorization":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

// функция для передачи в mux
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

		switch code {
		case 401:
			w.Header().Set("WWW-Authenticate", "Bearer")
			w.WriteHeader(http.StatusUnauthorized)
			return nil
		default:
			w.WriteHeader(code)
		}

	}
	return nil
}
