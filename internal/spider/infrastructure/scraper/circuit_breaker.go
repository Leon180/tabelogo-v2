package scraper

import (
	"time"

	"github.com/Leon180/tabelogo-v2/internal/spider/infrastructure/metrics"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

// CircuitBreakerConfig holds circuit breaker configuration
type CircuitBreakerConfig struct {
	MaxRequests uint32        // Max requests allowed in half-open state
	Interval    time.Duration // Interval to clear internal counts
	Timeout     time.Duration // Timeout to switch from open to half-open
}

// DefaultCircuitBreakerConfig returns default circuit breaker settings
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		MaxRequests: 3,                // Allow 3 requests in half-open state
		Interval:    60 * time.Second, // Clear counts every 60s
		Timeout:     30 * time.Second, // Try half-open after 30s
	}
}

// NewCircuitBreaker creates a new circuit breaker for the scraper
func NewCircuitBreaker(logger *zap.Logger, metrics *metrics.SpiderMetrics, config CircuitBreakerConfig) *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        "tabelog-scraper",
		MaxRequests: config.MaxRequests,
		Interval:    config.Interval,
		Timeout:     config.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Open circuit after 5 consecutive failures
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			logger.Warn("Circuit breaker state changed",
				zap.String("circuit", name),
				zap.String("from", from.String()),
				zap.String("to", to.String()),
			)

			// Record state change in metrics
			var stateValue float64
			switch to {
			case gobreaker.StateClosed:
				stateValue = 0
			case gobreaker.StateOpen:
				stateValue = 1
				metrics.RecordCircuitBreakerFailure(name)
			case gobreaker.StateHalfOpen:
				stateValue = 2
			}
			metrics.SetCircuitBreakerState(name, stateValue)
		},
	}

	cb := gobreaker.NewCircuitBreaker(settings)

	// Initialize state metric
	metrics.SetCircuitBreakerState("tabelog-scraper", 0) // Start as closed

	return cb
}

// IsCircuitBreakerError checks if error is from circuit breaker
func IsCircuitBreakerError(err error) bool {
	return err == gobreaker.ErrOpenState || err == gobreaker.ErrTooManyRequests
}
