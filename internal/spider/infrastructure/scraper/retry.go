package scraper

import (
	"context"
	"fmt"
	"math"
	"time"

	"go.uber.org/zap"
)

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxRetries     int           // Maximum number of retry attempts
	InitialDelay   time.Duration // Initial delay before first retry
	MaxDelay       time.Duration // Maximum delay between retries
	BackoffFactor  float64       // Exponential backoff multiplier
	RetryableTypes []ErrorType   // Error types that should be retried
}

// DefaultRetryConfig returns default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:    3,
		InitialDelay:  1 * time.Second,
		MaxDelay:      10 * time.Second,
		BackoffFactor: 2.0,
		RetryableTypes: []ErrorType{
			ErrorTypeTransient,
			ErrorTypeRateLimit,
		},
	}
}

// RetryableFunc is a function that can be retried
type RetryableFunc func() error

// WithRetry executes a function with retry logic
func WithRetry(ctx context.Context, logger *zap.Logger, config RetryConfig, fn RetryableFunc) error {
	var lastErr error

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// Execute the function
		err := fn()
		if err == nil {
			if attempt > 0 {
				logger.Info("Operation succeeded after retry",
					zap.Int("attempt", attempt),
				)
			}
			return nil
		}

		lastErr = err

		// Check if we should retry
		errType := ClassifyError(err)
		shouldRetry := false
		for _, retryableType := range config.RetryableTypes {
			if errType == retryableType {
				shouldRetry = true
				break
			}
		}

		// Don't retry permanent errors
		if !shouldRetry {
			logger.Warn("Non-retryable error encountered",
				zap.Error(err),
				zap.String("error_type", errorTypeToString(errType)),
			)
			return err
		}

		// Don't retry if we've exhausted attempts
		if attempt >= config.MaxRetries {
			logger.Error("Max retries exhausted",
				zap.Error(err),
				zap.Int("attempts", attempt+1),
			)
			return fmt.Errorf("max retries (%d) exhausted: %w", config.MaxRetries, err)
		}

		// Calculate delay with exponential backoff
		delay := calculateBackoff(attempt, config)

		logger.Warn("Operation failed, retrying",
			zap.Error(err),
			zap.Int("attempt", attempt+1),
			zap.Int("max_retries", config.MaxRetries),
			zap.Duration("delay", delay),
			zap.String("error_type", errorTypeToString(errType)),
		)

		// Wait before retry, respecting context cancellation
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled: %w", ctx.Err())
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	return lastErr
}

// calculateBackoff calculates the delay for the next retry attempt
func calculateBackoff(attempt int, config RetryConfig) time.Duration {
	// Exponential backoff: initialDelay * (backoffFactor ^ attempt)
	delay := float64(config.InitialDelay) * math.Pow(config.BackoffFactor, float64(attempt))

	// Cap at max delay
	if delay > float64(config.MaxDelay) {
		delay = float64(config.MaxDelay)
	}

	return time.Duration(delay)
}

// errorTypeToString converts ErrorType to string for logging
func errorTypeToString(errType ErrorType) string {
	switch errType {
	case ErrorTypeTransient:
		return "transient"
	case ErrorTypePermanent:
		return "permanent"
	case ErrorTypeRateLimit:
		return "rate_limit"
	default:
		return "unknown"
	}
}
