package logger

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const loggerKey contextKey = "logger"

// WithContext adds logger to context
func WithContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext retrieves logger from context
// If logger is not found in context, returns a default logger
func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return GetLogger()
	}

	if logger, ok := ctx.Value(loggerKey).(*zap.Logger); ok && logger != nil {
		return logger
	}

	return GetLogger()
}

// WithFields adds fields to logger in context
func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	logger := FromContext(ctx)
	return WithContext(ctx, logger.With(fields...))
}
