package logger

import (
	"context"
	"gin-demo/internal/tracing"

	"go.uber.org/zap"
)

type contextKey string

const loggerContextKey contextKey = "contextLogger"

type ContextLogger struct {
	*zap.Logger
	requestID string
}

func NewContextLogger(ctx context.Context) *ContextLogger {
	requestID := tracing.FromContext(ctx)
	return &ContextLogger{
		Logger:    logger.With(zap.String("requestID", requestID)),
		requestID: requestID,
	}
}

func FromContext(ctx context.Context) *ContextLogger {
	if l, ok := ctx.Value(loggerContextKey).(*ContextLogger); ok {
		return l
	}
	return NewContextLogger(ctx)
}

func WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerContextKey, NewContextLogger(ctx))
}

func (l *ContextLogger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

func (l *ContextLogger) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

// 添加其他需要的日志级别方法...
