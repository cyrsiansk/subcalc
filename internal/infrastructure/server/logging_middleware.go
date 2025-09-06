package server

import (
	"time"

	loggerpkg "subcalc/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func ZapRequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		reqID := c.GetHeader("X-Request-Id")
		if reqID == "" {
			reqID = uuid.New().String()
			c.Request.Header.Set("X-Request-Id", reqID)
		}

		fullPath := c.FullPath()
		if fullPath == "" {
			fullPath = c.Request.URL.Path
		}

		reqLogger := logger.With(
			zap.String("request_id", reqID),
			zap.String("http_method", c.Request.Method),
			zap.String("http_path", fullPath),
			zap.String("remote_addr", c.ClientIP()),
		)

		c.Set("logger", reqLogger)
		c.Request = c.Request.WithContext(loggerpkg.WithLogger(c.Request.Context(), reqLogger))

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		size := c.Writer.Size()
		if size < 0 {
			size = 0
		}

		reqLogger.Info("request finished",
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.Int("response_bytes", size),
			zap.String("user_agent", c.Request.UserAgent()),
		)
	}
}
