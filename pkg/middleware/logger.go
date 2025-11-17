package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger returns a gin middleware for logging HTTP requests
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log after request is processed
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		fields := []zap.Field{
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		if len(c.Errors) > 0 {
			// Log errors if any
			for _, e := range c.Errors {
				logger.Error("Request error", append(fields, zap.Error(e))...)
			}
		} else {
			// Log normal request
			if statusCode >= 500 {
				logger.Error("Server error", fields...)
			} else if statusCode >= 400 {
				logger.Warn("Client error", fields...)
			} else {
				logger.Info("Request completed", fields...)
			}
		}
	}
}
