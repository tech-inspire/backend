package logger

import (
	"log/slog"
	"os"
)

type Environment string

const (
	Dev   Environment = "dev"
	Stage Environment = "stage"
	Prod  Environment = "prod"
)

type Logger struct {
	*slog.Logger
	StackTrace  bool
	Environment Environment
}

func Error(err error) slog.Attr {
	return slog.String("error", err.Error())
}

func New() *Logger {
	var (
		disableStacktrace = os.Getenv("LOGGER_DISABLE_STACKTRACE") == "true"
		environment       = Environment(os.Getenv("ENVIRONMENT"))
	)

	var handler slog.Handler = slog.NewJSONHandler(os.Stdout, nil)
	if environment == Dev {
		handler = slog.NewTextHandler(os.Stdout, nil)
	}

	l := slog.New(handler)
	slog.SetDefault(l)

	return &Logger{
		StackTrace:  !disableStacktrace,
		Logger:      l,
		Environment: environment,
	}
}
