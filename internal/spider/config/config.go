package config

import "time"

// SpiderConfig holds all configuration for the Spider Service
type SpiderConfig struct {
	// Worker pool configuration
	WorkerCount int `env:"SPIDER_WORKER_COUNT" envDefault:"20"`

	// Cache configuration
	CacheTTL time.Duration `env:"SPIDER_CACHE_TTL" envDefault:"24h"`

	// Circuit breaker configuration
	CircuitBreaker CircuitBreakerConfig

	// Rate limiting configuration
	RateLimit RateLimitConfig
}

// CircuitBreakerConfig holds circuit breaker settings
type CircuitBreakerConfig struct {
	MaxRequests uint32        `env:"SPIDER_CB_MAX_REQUESTS" envDefault:"3"`
	Interval    time.Duration `env:"SPIDER_CB_INTERVAL" envDefault:"60s"`
	Timeout     time.Duration `env:"SPIDER_CB_TIMEOUT" envDefault:"30s"`
}

// RateLimitConfig holds rate limiting settings
type RateLimitConfig struct {
	RequestsPerMinute int           `env:"SPIDER_RATE_LIMIT_RPM" envDefault:"60"`
	BurstSize         int           `env:"SPIDER_RATE_LIMIT_BURST" envDefault:"10"`
	CleanupInterval   time.Duration `env:"SPIDER_RATE_LIMIT_CLEANUP" envDefault:"5m"`
}

// DefaultConfig returns default configuration
func DefaultConfig() *SpiderConfig {
	return &SpiderConfig{
		WorkerCount: 20,
		CacheTTL:    24 * time.Hour,
		CircuitBreaker: CircuitBreakerConfig{
			MaxRequests: 3,
			Interval:    60 * time.Second,
			Timeout:     30 * time.Second,
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: 60,
			BurstSize:         10,
			CleanupInterval:   5 * time.Minute,
		},
	}
}
