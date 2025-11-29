package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP Metrics
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "map_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "map_http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Google API Metrics
	GoogleAPICallsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "map_google_api_calls_total",
			Help: "Total number of Google API calls",
		},
		[]string{"api_type", "status"},
	)

	GoogleAPIDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "map_google_api_duration_seconds",
			Help:    "Google API call latency in seconds",
			Buckets: []float64{0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"api_type"},
	)

	// Cache Metrics
	CacheHitsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "map_cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	CacheMissesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "map_cache_misses_total",
			Help: "Total number of cache misses",
		},
	)

	// Error Metrics
	ErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "map_errors_total",
			Help: "Total number of errors",
		},
		[]string{"type"},
	)

	// Active Requests
	ActiveRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "map_active_requests",
			Help: "Number of active requests",
		},
	)
)
