package http

import (
	"strconv"
	"time"

	"github.com/Leon180/tabelogo-v2/pkg/metrics"
	"github.com/gin-gonic/gin"
)

// MetricsMiddleware records HTTP request metrics
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Increment active requests
		metrics.ActiveRequests.Inc()
		defer metrics.ActiveRequests.Dec()

		c.Next()

		// Record metrics after request completes
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// Record HTTP request count
		metrics.HTTPRequestsTotal.WithLabelValues(
			c.Request.Method,
			path,
			status,
		).Inc()

		// Record HTTP request duration
		metrics.HTTPRequestDuration.WithLabelValues(
			c.Request.Method,
			path,
		).Observe(duration)
	}
}

// RecordSearchOperation records a search operation metric
func RecordSearchOperation(operation, status string) {
	// This can be used by handlers to record specific operations
	// For now, we rely on the HTTP metrics middleware
}
