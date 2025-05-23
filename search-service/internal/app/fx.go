package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/tech-inspire/backend/auth-service/pkg/jwt"
	"github.com/tech-inspire/backend/search-service/internal/api/metrics"
	"github.com/tech-inspire/backend/search-service/internal/api/rpc"
	"github.com/tech-inspire/backend/search-service/internal/api/rpc/handlers"
	"github.com/tech-inspire/backend/search-service/internal/clients"
	"github.com/tech-inspire/backend/search-service/internal/config"
	"github.com/tech-inspire/backend/search-service/internal/consumer"
	"github.com/tech-inspire/backend/search-service/internal/repository/nats"
	"github.com/tech-inspire/backend/search-service/internal/repository/postgres"
	"github.com/tech-inspire/backend/search-service/internal/service"
	"github.com/tech-inspire/backend/search-service/migrations"
	"github.com/tech-inspire/backend/search-service/pkg/logger"
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

		fx.Provide(clients.NewPostgres),
		fx.Invoke(migrations.ApplyMigrations),

		fx.Provide(
			fx.Annotate(postgres.NewSearchRepository, fx.As(new(service.SearchRepository))),
		),

		fx.Provide(
			clients.NewNatsJetstreamClient,
			fx.Annotate(consumer.NewImageEmbeddingsUpdatesConsumer),
			fx.Annotate(consumer.NewPostCreatedEventConsumer),
		),
		fx.Invoke(consumer.ImageEmbeddingsUpdatesConsumer.Start),
		fx.Invoke(consumer.PostCreatedEventConsumer.Start),

		fx.Provide(
			fx.Annotate(service.NewSearchService, fx.As(new(handlers.SearchService))),
			fx.Annotate(service.NewSearchService, fx.As(new(consumer.PostsEventProcessor))),
			fx.Annotate(service.NewSearchService, fx.As(new(consumer.ImageEmbeddingsUpdatesConsumerProcessor))),
		),

		fx.Provide(
			fx.Annotate(clients.NewEmbeddingServiceClient, fx.As(new(service.TextEmbeddingsGenerator))),
			fx.Annotate(nats.NewImageEmbeddingsEventDispatcher, fx.As(new(service.ImageEmbeddingsTaskManager))),
		),

		//

		fx.Provide(
			handlers.NewSearchHandler,
		),

		//

		fx.Provide(func(cfg *config.Config) (*jwt.Validator, error) {
			return jwt.NewValidatorFromURL(cfg.AuthJWKSPath)
		}),

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
