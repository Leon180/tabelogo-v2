package middleware

import (
	"strconv"
	"time"

	"github.com/Leon180/tabelogo-v2/pkg/metrics"
	"github.com/gin-gonic/gin"
)

// MetricsMiddleware records HTTP metrics
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Increment active requests
		metrics.ActiveRequests.Inc()
		defer metrics.ActiveRequests.Dec()

		// Process request
		c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		metrics.HTTPRequestsTotal.WithLabelValues(
			c.Request.Method,
			path,
			status,
		).Inc()

		metrics.HTTPRequestDuration.WithLabelValues(
			c.Request.Method,
			path,
		).Observe(duration)
	}
}
