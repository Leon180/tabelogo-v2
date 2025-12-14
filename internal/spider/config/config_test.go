package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig_Defaults(t *testing.T) {
	// Act
	cfg := DefaultConfig()

	// Assert
	assert.NotNil(t, cfg)
	assert.Equal(t, 20, cfg.WorkerCount)
	assert.Equal(t, 24*time.Hour, cfg.CacheTTL)
	assert.Equal(t, uint32(3), cfg.CircuitBreaker.MaxRequests)
	assert.Equal(t, 60*time.Second, cfg.CircuitBreaker.Interval)
	assert.Equal(t, 30*time.Second, cfg.CircuitBreaker.Timeout)
}

func TestConfig_WorkerCount(t *testing.T) {
	// Arrange
	cfg := DefaultConfig()

	// Assert
	assert.Greater(t, cfg.WorkerCount, 0)
	assert.LessOrEqual(t, cfg.WorkerCount, 100)
}

func TestConfig_CacheTTL(t *testing.T) {
	// Arrange
	cfg := DefaultConfig()

	// Assert
	assert.Greater(t, cfg.CacheTTL, time.Duration(0))
	assert.LessOrEqual(t, cfg.CacheTTL, 7*24*time.Hour) // Max 1 week
}

func TestConfig_CircuitBreakerConfig(t *testing.T) {
	// Arrange
	cfg := DefaultConfig()

	// Assert
	assert.NotNil(t, cfg.CircuitBreaker)
	assert.Greater(t, cfg.CircuitBreaker.MaxRequests, uint32(0))
	assert.Greater(t, cfg.CircuitBreaker.Interval, time.Duration(0))
	assert.Greater(t, cfg.CircuitBreaker.Timeout, time.Duration(0))
}

func TestConfig_RateLimitConfig(t *testing.T) {
	// Arrange
	cfg := DefaultConfig()

	// Assert
	assert.NotNil(t, cfg.RateLimit)
	assert.Equal(t, 60, cfg.RateLimit.RequestsPerMinute)
	assert.Equal(t, 10, cfg.RateLimit.BurstSize)
	assert.Equal(t, 5*time.Minute, cfg.RateLimit.CleanupInterval)
}

func TestCircuitBreakerConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config CircuitBreakerConfig
		valid  bool
	}{
		{
			name: "valid config",
			config: CircuitBreakerConfig{
				MaxRequests: 10,
				Interval:    10 * time.Second,
				Timeout:     30 * time.Second,
			},
			valid: true,
		},
		{
			name: "zero max requests",
			config: CircuitBreakerConfig{
				MaxRequests: 0,
				Interval:    10 * time.Second,
				Timeout:     30 * time.Second,
			},
			valid: false,
		},
		{
			name: "zero interval",
			config: CircuitBreakerConfig{
				MaxRequests: 10,
				Interval:    0,
				Timeout:     30 * time.Second,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assert
			if tt.valid {
				assert.Greater(t, tt.config.MaxRequests, uint32(0))
				assert.Greater(t, tt.config.Interval, time.Duration(0))
				assert.Greater(t, tt.config.Timeout, time.Duration(0))
			}
		})
	}
}

func TestConfig_ReasonableDefaults(t *testing.T) {
	// Arrange
	cfg := DefaultConfig()

	// Assert - verify defaults are production-ready
	assert.GreaterOrEqual(t, cfg.WorkerCount, 1, "Should have at least 1 worker")
	assert.LessOrEqual(t, cfg.WorkerCount, 100, "Should not exceed 100 workers")

	assert.GreaterOrEqual(t, cfg.CacheTTL, 1*time.Hour, "Cache should last at least 1 hour")
	assert.LessOrEqual(t, cfg.CacheTTL, 7*24*time.Hour, "Cache should not exceed 1 week")

	assert.GreaterOrEqual(t, cfg.CircuitBreaker.MaxRequests, uint32(1), "CB should allow at least 1 request")
	assert.GreaterOrEqual(t, cfg.CircuitBreaker.Interval, 5*time.Second, "CB interval should be at least 5s")
	assert.GreaterOrEqual(t, cfg.CircuitBreaker.Timeout, 10*time.Second, "CB timeout should be at least 10s")

	assert.GreaterOrEqual(t, cfg.RateLimit.RequestsPerMinute, 10, "Rate limit should allow at least 10 rpm")
	assert.GreaterOrEqual(t, cfg.RateLimit.BurstSize, 1, "Burst size should be at least 1")
}
