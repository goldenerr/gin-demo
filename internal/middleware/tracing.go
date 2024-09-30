package middleware

import (
	"context"

	"gin-demo/internal/constants"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func Tracing(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		ctx := context.WithValue(c.Request.Context(), constants.RequestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)

		logger.Debug("Request started",
			zap.String("request_id", requestID),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)

		c.Next()

		logger.Debug("Request completed",
			zap.String("request_id", requestID),
			zap.Int("status", c.Writer.Status()),
		)
	}
}

// GetRequestID 从上下文中获取 request_id
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(constants.RequestIDKey).(string); ok {
		return requestID
	}
	return "unknown"
}
