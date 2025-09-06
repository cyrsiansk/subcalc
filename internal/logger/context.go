package logger

import (
	"context"

	"go.uber.org/zap"
)

type ctxKeyType struct{}

var loggerCtxKey = ctxKeyType{}

func WithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, l)
}
