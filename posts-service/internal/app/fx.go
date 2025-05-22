package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	redigo "github.com/redis/go-redis/v9"
	"github.com/scylladb/gocqlx/v3"
	"github.com/tech-inspire/backend/auth-service/pkg/jwt"
	"github.com/tech-inspire/backend/posts-service/internal/api/metrics"
	"github.com/tech-inspire/backend/posts-service/internal/api/rpc"
	"github.com/tech-inspire/backend/posts-service/internal/api/rpc/handlers"
	"github.com/tech-inspire/backend/posts-service/internal/clients"
	"github.com/tech-inspire/backend/posts-service/internal/config"
	"github.com/tech-inspire/backend/posts-service/internal/repository/cache"
	"github.com/tech-inspire/backend/posts-service/internal/repository/nats"
	"github.com/tech-inspire/backend/posts-service/internal/repository/redis"
	"github.com/tech-inspire/backend/posts-service/internal/repository/s3"
	"github.com/tech-inspire/backend/posts-service/internal/repository/scylla"
	"github.com/tech-inspire/backend/posts-service/internal/service"
	"github.com/tech-inspire/backend/posts-service/migrations"
	"github.com/tech-inspire/backend/posts-service/pkg/generator"
	"github.com/tech-inspire/backend/posts-service/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
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
		fx.Invoke(migrations.ApplyMigrationsFX),

		fx.Provide(
			fx.Annotate(scylla.NewPostsRepository),
		),

		fx.Provide(func(cfg *config.Config) (*jwt.Validator, error) {
			return jwt.NewValidatorFromURL(cfg.AuthJWKSPath)
		}),

		fx.Provide(
			fx.Annotate(clients.NewRedis, fx.As(new(redigo.UniversalClient))),
			redis.NewPostRepository,
			fx.Annotate(redis.NewPendingImageUploadsRepository, fx.As(new(service.PendingImagesRepository))),
		),

		fx.Provide(
			fx.Annotate(cache.NewPostsRepository, fx.As(new(service.PostsRepository))),
		),

		fx.Provide(
			fx.Annotate(nats.NewPostsEventDispatcher, fx.As(new(service.PostsEventDispatcher))),
		),

		fx.Provide(
			clients.NewS3Client,
			fx.Annotate(imagestorage.New, fx.As(new(service.ImageStorage))),
		),

		fx.Provide(

			fx.Annotate(service.NewPostsService, fx.As(new(handlers.PostsService))),
		),

		//

		fx.Provide(
			handlers.NewPostsHandler,
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
		l.Error("failed to validate fx app", logger.Error(err))
		return
	}

	app := fx.New(options...)

	err := app.Start(context.Background())
	if err != nil {
		l.Error("failed to start app", logger.Error(err))
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	<-ctx.Done()

	err = app.Stop(context.Background())
	if err != nil {
		l.Warn("failed to stop app", logger.Error(err))
	}
}
