package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	redigo "github.com/redis/go-redis/v9"
	"github.com/tech-inspire/service/auth-service/internal/api/jwt"
	"github.com/tech-inspire/service/auth-service/internal/api/metrics"
	"github.com/tech-inspire/service/auth-service/internal/api/rpc"
	"github.com/tech-inspire/service/auth-service/internal/api/rpc/handlers"
	"github.com/tech-inspire/service/auth-service/internal/config"
	avatarstorage "github.com/tech-inspire/service/auth-service/internal/repository/avatar"
	"github.com/tech-inspire/service/auth-service/internal/repository/postgres"
	"github.com/tech-inspire/service/auth-service/internal/repository/postgres/sqlc"
	"github.com/tech-inspire/service/auth-service/internal/repository/redis"
	"github.com/tech-inspire/service/auth-service/internal/service"
	"github.com/tech-inspire/service/auth-service/pkg/clients"
	"github.com/tech-inspire/service/auth-service/pkg/generator"
	"github.com/tech-inspire/service/auth-service/pkg/logger"
	"github.com/tech-inspire/service/auth-service/pkg/mail"
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

		fx.Provide(clients.NewPostgres),
		fx.Provide(func(pool *pgxpool.Pool) *sqlc.Queries {
			return sqlc.New(pool)
		}),

		fx.Provide(
			fx.Annotate(clients.NewRedis, fx.As(new(redigo.UniversalClient))),
		),

		fx.Provide(
			fx.Annotate(mail.NewClient, fx.As(new(service.MailClient))),
		),

		fx.Provide(
			fx.Annotate(postgres.NewUserRepository, fx.As(new(service.UserRepository))),

			fx.Annotate(redis.NewSessionRepository, fx.As(new(service.SessionRepository))),
			fx.Annotate(redis.NewCodesRepository, fx.As(new(service.ConfirmationCodesRepository))),
			fx.Annotate(redis.NewResetCodesRepository, fx.As(new(service.ResetPasswordCodesRepository))),
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
			handlers.NewAuthHandler,
			handlers.NewUserHandler,
		),

		//

		fx.Provide(jwt.NewSigner),
		fx.Provide(
			fx.Annotate(generator.New, fx.As(new(service.Generator))),
		),

		fx.Invoke(metrics.NewServer),
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
