package scraper

import (
	"errors"
	"testing"
	"time"

	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

func TestCircuitBreaker(t *testing.T) {
	logger := zap.NewNop()
	config := DefaultCircuitBreakerConfig()
	cb := NewCircuitBreaker(logger, config)

	t.Run("allows requests when closed", func(t *testing.T) {
		_, err := cb.Execute(func() (interface{}, error) {
			return "success", nil
		})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("opens after consecutive failures", func(t *testing.T) {
		// Reset circuit breaker
		cb = NewCircuitBreaker(logger, config)

		// Cause failures to open circuit
		testErr := errors.New("test error")
		for i := 0; i < 5; i++ {
			cb.Execute(func() (interface{}, error) {
				return nil, testErr
			})
		}

		// Circuit should be open now
		_, err := cb.Execute(func() (interface{}, error) {
			return "should not execute", nil
		})

		if err != gobreaker.ErrOpenState {
			t.Errorf("expected ErrOpenState, got %v", err)
		}
	})

	t.Run("transitions to half-open after timeout", func(t *testing.T) {
		// Use shorter timeout for testing
		shortConfig := CircuitBreakerConfig{
			MaxRequests: 1,
			Interval:    1 * time.Second,
			Timeout:     100 * time.Millisecond,
		}
		cb = NewCircuitBreaker(logger, shortConfig)

		// Open the circuit
		testErr := errors.New("test error")
		for i := 0; i < 5; i++ {
			cb.Execute(func() (interface{}, error) {
				return nil, testErr
			})
		}

		// Wait for timeout
		time.Sleep(150 * time.Millisecond)

		// Should allow one request in half-open state
		_, err := cb.Execute(func() (interface{}, error) {
			return "success", nil
		})
		if err != nil {
			t.Errorf("expected request to succeed in half-open state, got %v", err)
		}
	})
}

func TestIsCircuitBreakerError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "open state error",
			err:      gobreaker.ErrOpenState,
			expected: true,
		},
		{
			name:     "too many requests error",
			err:      gobreaker.ErrTooManyRequests,
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("some error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsCircuitBreakerError(tt.err)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDefaultCircuitBreakerConfig(t *testing.T) {
	config := DefaultCircuitBreakerConfig()

	if config.MaxRequests != 3 {
		t.Errorf("expected MaxRequests=3, got %d", config.MaxRequests)
	}
	if config.Interval != 60*time.Second {
		t.Errorf("expected Interval=60s, got %v", config.Interval)
	}
	if config.Timeout != 30*time.Second {
		t.Errorf("expected Timeout=30s, got %v", config.Timeout)
	}
}
