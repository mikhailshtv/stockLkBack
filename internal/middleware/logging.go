package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/mikhailshtv/stockLkBack/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func generateTraceID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		traceID := generateTraceID()
		c.Set("trace_id", traceID)

		c.Header("X-Trace-ID", traceID)

		log := logger.GetLogger()

		log.Info("request started",
			zap.String("trace_id", traceID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)

		c.Next()

		duration := time.Since(start)

		status := c.Writer.Status()

		var logFunc func(string, ...zap.Field)
		switch {
		case status >= 500:
			logFunc = log.Error
		case status >= 400:
			logFunc = log.Warn
		default:
			logFunc = log.Info
		}

		logFunc("request completed",
			zap.String("trace_id", traceID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", status),
			zap.Duration("duration", duration),
			zap.Int("size", c.Writer.Size()),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
