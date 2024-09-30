package router

import (
	"bytes"
	"gin-demo/internal/handler"
	"gin-demo/internal/logger"
	"gin-demo/internal/tracing"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := tracing.NewContext(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)
		c.Header("X-Request-ID", tracing.FromContext(ctx))
		c.Next()
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := tracing.FromContext(c.Request.Context())

		// 记录请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 包装响应写入器以捕获响应体
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		// 记录响应体
		responseBody := blw.body.String()

		end := time.Now()
		latency := end.Sub(start)

		logger.Info("Request",
			zap.String("requestID", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("requestBody", string(requestBody)),
			zap.Int("status", c.Writer.Status()),
			zap.String("responseBody", responseBody),
			zap.Duration("latency", latency),
			zap.String("ip", c.ClientIP()),
		)
	}
}

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(TracingMiddleware())
	r.Use(LoggerMiddleware())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/users", handler.GetUsers)
	}

	return r
}
