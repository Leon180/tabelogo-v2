package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// SpiderMetrics holds all Prometheus metrics for the spider service
type SpiderMetrics struct {
	// Scraping metrics
	ScrapeRequestsTotal *prometheus.CounterVec
	ScrapeDuration      *prometheus.HistogramVec
	ScrapeErrorsTotal   *prometheus.CounterVec
	RestaurantsScraped  *prometheus.CounterVec

	// Job processing metrics
	JobsTotal      *prometheus.CounterVec
	JobDuration    *prometheus.HistogramVec
	WorkerPoolSize prometheus.Gauge
	JobQueueLength prometheus.Gauge

	// Cache metrics
	CacheHitsTotal   *prometheus.CounterVec
	CacheMissesTotal *prometheus.CounterVec
	CacheSizeBytes   *prometheus.GaugeVec

	// Circuit breaker metrics
	CircuitBreakerState    *prometheus.GaugeVec
	CircuitBreakerFailures *prometheus.CounterVec
}

// NewSpiderMetrics creates and registers all spider service metrics
func NewSpiderMetrics() *SpiderMetrics {
	return &SpiderMetrics{
		// Scraping metrics
		ScrapeRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "spider_scrape_requests_total",
				Help: "Total number of scrape requests",
			},
			[]string{"status"}, // success, failed, cached
		),
		ScrapeDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "spider_scrape_duration_seconds",
				Help:    "Duration of scrape operations in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation", "status"}, // operation: search, details; status: success, failure
		),
		ScrapeErrorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "spider_scrape_errors_total",
				Help: "Total number of scrape errors by type",
			},
			[]string{"error_type"}, // network, parse, not_found, circuit_breaker
		),
		RestaurantsScraped: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "spider_restaurants_scraped_total",
				Help: "Total number of restaurants scraped by status",
			},
			[]string{"status"}, // success, failure
		),

		// Job processing metrics
		JobsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "spider_jobs_total",
				Help: "Total number of jobs by status",
			},
			[]string{"status"}, // pending, running, completed, failed
		),
		JobDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "spider_job_duration_seconds",
				Help:    "Duration of job processing in seconds",
				Buckets: []float64{0.5, 1, 2, 5, 10, 30, 60},
			},
			[]string{"status"}, // completed, failed
		),
		WorkerPoolSize: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "spider_worker_pool_size",
				Help: "Number of active workers in the pool",
			},
		),
		JobQueueLength: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "spider_job_queue_length",
				Help: "Number of pending jobs in the queue",
			},
		),

		// Cache metrics
		CacheHitsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "spider_cache_hits_total",
				Help: "Total number of cache hits by cache type",
			},
			[]string{"cache_type"}, // result, job
		),
		CacheMissesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "spider_cache_misses_total",
				Help: "Total number of cache misses by cache type",
			},
			[]string{"cache_type"}, // result, job
		),
		CacheSizeBytes: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "spider_cache_size_bytes",
				Help: "Approximate size of cache in bytes by cache type",
			},
			[]string{"cache_type"}, // result, job
		),

		// Circuit breaker metrics
		CircuitBreakerState: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "spider_circuit_breaker_state",
				Help: "Circuit breaker state (0=closed, 1=open, 2=half-open)",
			},
			[]string{"circuit"},
		),
		CircuitBreakerFailures: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "spider_circuit_breaker_failures_total",
				Help: "Total number of circuit breaker failures",
			},
			[]string{"circuit"},
		),
	}
}

// RecordScrapeRequest records a scrape request with its status
func (m *SpiderMetrics) RecordScrapeRequest(status string) {
	m.ScrapeRequestsTotal.WithLabelValues(status).Inc()
}

// RecordScrapeDuration records the duration of a scrape operation
func (m *SpiderMetrics) RecordScrapeDuration(operation, status string, duration float64) {
	m.ScrapeDuration.WithLabelValues(operation, status).Observe(duration)
}

// RecordScrapeError records a scrape error by type
func (m *SpiderMetrics) RecordScrapeError(errorType string) {
	m.ScrapeErrorsTotal.WithLabelValues(errorType).Inc()
}

// RecordRestaurantsScraped increments the total restaurants scraped counter by status
func (m *SpiderMetrics) RecordRestaurantsScraped(status string, count int) {
	m.RestaurantsScraped.WithLabelValues(status).Add(float64(count))
}

// RecordJob records a job by status
func (m *SpiderMetrics) RecordJob(status string) {
	m.JobsTotal.WithLabelValues(status).Inc()
}

// RecordJobDuration records the duration of a job
func (m *SpiderMetrics) RecordJobDuration(status string, duration float64) {
	m.JobDuration.WithLabelValues(status).Observe(duration)
}

// SetWorkerPoolSize sets the current worker pool size
func (m *SpiderMetrics) SetWorkerPoolSize(size int) {
	m.WorkerPoolSize.Set(float64(size))
}

// SetJobQueueLength sets the current job queue length
func (m *SpiderMetrics) SetJobQueueLength(length int) {
	m.JobQueueLength.Set(float64(length))
}

// RecordCacheHit increments the cache hit counter for a specific cache type
func (m *SpiderMetrics) RecordCacheHit(cacheType string) {
	m.CacheHitsTotal.WithLabelValues(cacheType).Inc()
}

// RecordCacheMiss increments the cache miss counter for a specific cache type
func (m *SpiderMetrics) RecordCacheMiss(cacheType string) {
	m.CacheMissesTotal.WithLabelValues(cacheType).Inc()
}

// SetCacheSize sets the approximate cache size in bytes for a specific cache type
func (m *SpiderMetrics) SetCacheSize(cacheType string, bytes int64) {
	m.CacheSizeBytes.WithLabelValues(cacheType).Set(float64(bytes))
}

// SetCircuitBreakerState sets the circuit breaker state
// 0 = closed, 1 = open, 2 = half-open
func (m *SpiderMetrics) SetCircuitBreakerState(circuit string, state float64) {
	m.CircuitBreakerState.WithLabelValues(circuit).Set(state)
}

// RecordCircuitBreakerFailure records a circuit breaker failure
func (m *SpiderMetrics) RecordCircuitBreakerFailure(circuit string) {
	m.CircuitBreakerFailures.WithLabelValues(circuit).Inc()
}
