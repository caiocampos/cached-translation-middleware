package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestLogger returns a Gin middleware that logs each request using zap.
func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		logger.Info("request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.FullPath()),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("ip", c.ClientIP()),
		)
	}
}
