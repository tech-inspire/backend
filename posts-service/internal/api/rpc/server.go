package rpc

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"

	"connectrpc.com/authn"
	"connectrpc.com/connect"
	connectcors "connectrpc.com/cors"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/validate"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"github.com/tech-inspire/api-contracts/api/gen/go/posts/v1/postsv1connect"
	authjwt "github.com/tech-inspire/backend/auth-service/pkg/jwt"
	"github.com/tech-inspire/backend/posts-service/internal/api/metrics"
	"github.com/tech-inspire/backend/posts-service/internal/api/rpc/handlers"
	"github.com/tech-inspire/backend/posts-service/internal/api/rpc/middleware"
	"github.com/tech-inspire/backend/posts-service/internal/config"
	authmiddleware "github.com/tech-inspire/backend/posts-service/pkg/jwt/middleware"
	"github.com/tech-inspire/backend/posts-service/pkg/logger"
	"go.uber.org/fx"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func CORSMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: cfg.Server.CORSAllowedOrigins,
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: connectcors.AllowedHeaders(),
		ExposedHeaders: connectcors.ExposedHeaders(),
	})
	return corsMiddleware.Handler
}

type Params struct {
	fx.In

	Logger *logger.Logger

	JwtValidator *authjwt.Validator

	PostsHandler *handlers.PostsHandler
}

func RegisterRoutes(params Params, r *chi.Mux) error {
	validateInterceptor, err := validate.NewInterceptor()
	if err != nil {
		return fmt.Errorf("protovalidate: create interceptor: %w", err)
	}

	type postsService struct {
		*handlers.PostsHandler
	}

	authServicePath, authServiceHandler := postsv1connect.NewPostsServiceHandler(
		postsService{
			params.PostsHandler,
		},
		connect.WithInterceptors(
			middleware.ErrorInterceptor(params.Logger, postsv1connect.PostsServiceName),
			validateInterceptor,
		),
	)

	// without auth
	noAuthenticationProcedures := []string{
		postsv1connect.PostsServiceGetPostByIDProcedure,
		postsv1connect.PostsServiceGetPostsProcedure,
	}

	authMiddleware := authn.NewMiddleware(
		authmiddleware.New(params.JwtValidator, noAuthenticationProcedures),
	)

	reflector := grpcreflect.NewStaticReflector(postsv1connect.PostsServiceName)
	r.Mount(grpcreflect.NewHandlerV1(reflector))
	r.Mount(grpcreflect.NewHandlerV1Alpha(reflector))

	r.Mount(authServicePath, authMiddleware.Wrap(authServiceHandler))

	return nil
}

func NewServer(lc fx.Lifecycle, cfg *config.Config) (*chi.Mux, error) {
	r := chi.NewRouter()

	r.Use(metrics.RecordMiddleware)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(CORSMiddleware(cfg))

	srv := &http.Server{
		Handler: h2c.NewHandler(r, new(http2.Server)),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", cfg.Server.Address)
			if err != nil {
				return err
			}

			slog.Info("server started", slog.String("listening", cfg.Server.Address))

			go func() {
				if err := srv.Serve(ln); err != nil {
					log.Fatal(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return r, nil
}
