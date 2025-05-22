package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	redigo "github.com/redis/go-redis/v9"
	"github.com/scylladb/gocqlx/v3"
	"github.com/tech-inspire/backend/posts-service/internal/api/metrics"
	"github.com/tech-inspire/backend/posts-service/internal/api/rpc"
	"github.com/tech-inspire/backend/posts-service/internal/api/rpc/handlers"
	"github.com/tech-inspire/backend/posts-service/internal/clients"
	"github.com/tech-inspire/backend/posts-service/internal/config"
	avatarstorage "github.com/tech-inspire/backend/posts-service/internal/repository/avatar"
	"github.com/tech-inspire/backend/posts-service/internal/repository/redis"
	"github.com/tech-inspire/backend/posts-service/internal/service"
	"github.com/tech-inspire/backend/posts-service/migrations"
	"github.com/tech-inspire/backend/posts-service/pkg/generator"
	"github.com/tech-inspire/backend/posts-service/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func Run() {
	l := logger.New()

	options := []fx.Option{
		fx.Supply(l), // *zap.Logger
		fx.WithLogger(func(log *logger.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{Logger: log.Logger}
		}),

		fx.Provide(config.New),

		fx.Provide(
			clients.NewScyllaDBClient,
			gocqlx.NewSession,
		),
		fx.Invoke(migrations.ApplyMigrations),

		fx.Provide(
			fx.Annotate(clients.NewRedis, fx.As(new(redigo.UniversalClient))),
		),

		fx.Provide(
			fx.Annotate(redis.NewPendingImageUploadsRepository, fx.As(new(service.SessionRepository))),
		),

		fx.Provide(
			clients.NewS3Client,
			fx.Annotate(avatarstorage.New, fx.As(new(service.AvatarStorage))),
		),

		fx.Provide(
			fx.Annotate(service.NewAuthService, fx.As(new(handlers.AuthService))),
			fx.Annotate(service.NewAuthService),
			fx.Annotate(service.NewUserService, fx.As(new(handlers.UserService))),
			fx.Annotate(service.NewAvatarService, fx.As(new(handlers.AvatarService))),
		),

		//

		fx.Provide(
			// handlers.NewAuthHandler,
			handlers.PostsHandler{},
		),

		//

		//fx.Provide(jwt.NewSigner),
		//fx.Provide(func(signer *jwt.Signer) (*authjwt.Validator, error) {
		//	return signer.Validator()
		//}),

		fx.Provide(
			fx.Annotate(generator.New, fx.As(new(service.Generator))),
		),

		fx.Invoke(metrics.NewServer),
		fx.Invoke(metrics.RegisterCollectors),

		fx.Provide(rpc.NewServer),
		fx.Invoke(rpc.RegisterRoutes),
	}

	if err := fx.ValidateApp(options...); err != nil {
		l.Error("failed to validate fx app", zap.Error(err))
		return
	}

	app := fx.New(options...)

	err := app.Start(context.Background())
	if err != nil {
		l.Error("failed to start app", zap.Error(err))
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	<-ctx.Done()

	err = app.Stop(context.Background())
	if err != nil {
		l.Warn("failed to stop app", zap.Error(err))
	}
}
