package metrics

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	middlewarestd "github.com/slok/go-http-metrics/middleware/std"
	"github.com/tech-inspire/backend/auth-service/internal/config"
	"github.com/tech-inspire/backend/auth-service/pkg/logger"
	"go.uber.org/fx"
)

func RecordMiddleware(next http.Handler) http.Handler {
	m := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	return middlewarestd.Handler("", m, next)
}

type Params struct {
	fx.In

	Config *config.Config
	Logger *logger.Logger
}
type Result struct {
	fx.Out
	Server *fiber.App `name:"metrics-api"`
}

func NewServer(lc fx.Lifecycle, in Params) error {
	r := chi.NewRouter()

	r.Get("/health", func(writer http.ResponseWriter, request *http.Request) {
		return
	})
	r.Mount("/debug", chimiddleware.Profiler())
	r.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", in.Config.Server.MetricsAddress)
			if err != nil {
				return err
			}

			slog.Info("metrics server started", slog.String("listening", in.Config.Server.MetricsAddress))

			go func() {
				if err := srv.Serve(ln); err != nil {
					slog.Error("failed to start metrics server", logger.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return nil
}
