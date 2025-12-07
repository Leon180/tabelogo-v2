package http

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "restaurant_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "restaurant_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	restaurantOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "restaurant_operations_total",
			Help: "Total number of restaurant operations",
		},
		[]string{"operation", "status"},
	)

	favoriteOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "restaurant_favorite_operations_total",
			Help: "Total number of favorite operations",
		},
		[]string{"operation", "status"},
	)

	databaseQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "restaurant_database_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"operation", "status"},
	)
)

// MetricsMiddleware records HTTP request metrics
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			path,
			status,
		).Inc()

		httpRequestDuration.WithLabelValues(
			c.Request.Method,
			path,
		).Observe(duration)
	}
}

// RecordRestaurantOperation records a restaurant operation metric
func RecordRestaurantOperation(operation, status string) {
	restaurantOperationsTotal.WithLabelValues(operation, status).Inc()
}

// RecordFavoriteOperation records a favorite operation metric
func RecordFavoriteOperation(operation, status string) {
	favoriteOperationsTotal.WithLabelValues(operation, status).Inc()
}

// RecordDatabaseQuery records a database query metric
func RecordDatabaseQuery(operation, status string) {
	databaseQueriesTotal.WithLabelValues(operation, status).Inc()
}

// CORSMiddleware handles Cross-Origin Resource Sharing (CORS)
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-Cache-Status, X-Data-Source, X-Data-Age")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
