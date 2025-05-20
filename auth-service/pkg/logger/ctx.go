package logger

import (
	"context"
)

type loggerCtxKey struct{}

func (l *Logger) Ctx(ctx context.Context) *Logger {
	logger := ctx.Value(loggerCtxKey{})
	if logger != nil {
		return logger.(*Logger)
	}

	return l
}

func (l *Logger) With(fields ...any) *Logger {
	return &Logger{
		Logger:      l.Logger.With(fields),
		Environment: l.Environment,
	}
}

func AddToCtx(logger *Logger) (key any, value *Logger) {
	return loggerCtxKey{}, logger
}
