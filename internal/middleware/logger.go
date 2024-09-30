package middleware

import (
	"bytes"
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

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ctx := tracing.NewContext(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)

		log := logger.FromContext(ctx)

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 记录请求信息
		log.Info("Request started",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("requestBody", string(requestBody)),
			zap.String("requestID", tracing.FromContext(ctx)))

		// 包装响应写入器以捕获响应体
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		// 记录响应信息
		duration := time.Since(start)
		log.Info("Request completed",
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.String("responseBody", blw.body.String()),
			zap.String("requestID", tracing.FromContext(ctx)))
	}
}
