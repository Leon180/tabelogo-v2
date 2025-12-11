package services

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// DynamicRateLimiter implements adaptive rate limiting with backoff
type DynamicRateLimiter struct {
	limiter       *rate.Limiter
	mu            sync.RWMutex
	baseRate      float64 // requests per second
	currentRate   float64
	minRate       float64 // minimum 5 req/min = 0.083 req/s
	backoffFactor float64
	recoveryRate  float64
	logger        *zap.Logger
}

// NewDynamicRateLimiter creates a new dynamic rate limiter
func NewDynamicRateLimiter(requestsPerMinute int, logger *zap.Logger) *DynamicRateLimiter {
	baseRate := float64(requestsPerMinute) / 60.0 // Convert to req/sec
	minRate := 5.0 / 60.0                         // 5 req/min minimum

	limiter := rate.NewLimiter(rate.Limit(baseRate), 1)

	return &DynamicRateLimiter{
		limiter:       limiter,
		baseRate:      baseRate,
		currentRate:   baseRate,
		minRate:       minRate,
		backoffFactor: 0.5, // Reduce by 50% on rate limit
		recoveryRate:  1.1, // Increase by 10% on success
		logger:        logger.With(zap.String("component", "rate_limiter")),
	}
}

// Wait waits for permission to proceed
func (r *DynamicRateLimiter) Wait(ctx context.Context) error {
	return r.limiter.Wait(ctx)
}

// OnRateLimitHit is called when a 429 response is received
func (r *DynamicRateLimiter) OnRateLimitHit() {
	r.mu.Lock()
	defer r.mu.Unlock()

	oldRate := r.currentRate
	r.currentRate = r.currentRate * r.backoffFactor

	// Don't go below minimum
	if r.currentRate < r.minRate {
		r.currentRate = r.minRate
	}

	// Update limiter
	r.limiter.SetLimit(rate.Limit(r.currentRate))

	r.logger.Warn("Rate limit hit, reducing rate",
		zap.Float64("old_rate_per_sec", oldRate),
		zap.Float64("new_rate_per_sec", r.currentRate),
		zap.Float64("old_rate_per_min", oldRate*60),
		zap.Float64("new_rate_per_min", r.currentRate*60),
	)
}

// OnSuccess is called after a successful request
func (r *DynamicRateLimiter) OnSuccess() {
	r.mu.Lock()
	defer r.mu.Unlock()

	oldRate := r.currentRate
	r.currentRate = r.currentRate * r.recoveryRate

	// Don't exceed base rate
	if r.currentRate > r.baseRate {
		r.currentRate = r.baseRate
	}

	// Only update if rate changed significantly (avoid too frequent updates)
	if r.currentRate != oldRate {
		r.limiter.SetLimit(rate.Limit(r.currentRate))

		r.logger.Debug("Rate limit recovered",
			zap.Float64("old_rate_per_sec", oldRate),
			zap.Float64("new_rate_per_sec", r.currentRate),
			zap.Float64("rate_per_min", r.currentRate*60),
		)
	}
}

// CurrentRate returns the current rate in requests per minute
func (r *DynamicRateLimiter) CurrentRate() float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.currentRate * 60.0
}

// SetRate manually sets the rate (for testing or manual adjustment)
func (r *DynamicRateLimiter) SetRate(requestsPerMinute float64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ratePerSec := requestsPerMinute / 60.0
	if ratePerSec < r.minRate {
		ratePerSec = r.minRate
	}
	if ratePerSec > r.baseRate {
		ratePerSec = r.baseRate
	}

	r.currentRate = ratePerSec
	r.limiter.SetLimit(rate.Limit(ratePerSec))

	r.logger.Info("Rate limit manually adjusted",
		zap.Float64("rate_per_min", requestsPerMinute),
	)
}
