package services

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestDynamicRateLimiter_BaseRate(t *testing.T) {
	logger := zap.NewNop()
	limiter := NewDynamicRateLimiter(30, logger) // 30 req/min

	// Should allow requests at base rate
	ctx := context.Background()
	start := time.Now()

	// First request should be immediate
	err := limiter.Wait(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	elapsed := time.Since(start)
	if elapsed > 100*time.Millisecond {
		t.Errorf("First request took too long: %v", elapsed)
	}

	// Check current rate
	rate := limiter.CurrentRate()
	if rate != 30.0 {
		t.Errorf("Expected rate 30, got %f", rate)
	}
}

func TestDynamicRateLimiter_OnRateLimitHit(t *testing.T) {
	logger := zap.NewNop()
	limiter := NewDynamicRateLimiter(30, logger)

	initialRate := limiter.CurrentRate()

	// Simulate rate limit hit
	limiter.OnRateLimitHit()

	newRate := limiter.CurrentRate()

	// Rate should be reduced by 50%
	expectedRate := initialRate * 0.5
	if newRate != expectedRate {
		t.Errorf("Expected rate %f after backoff, got %f", expectedRate, newRate)
	}

	// Hit again
	limiter.OnRateLimitHit()
	finalRate := limiter.CurrentRate()

	// Should not go below minimum (5 req/min)
	if finalRate < 5.0 {
		t.Errorf("Rate went below minimum: %f", finalRate)
	}
}

func TestDynamicRateLimiter_OnSuccess(t *testing.T) {
	logger := zap.NewNop()
	limiter := NewDynamicRateLimiter(30, logger)

	// Reduce rate first
	limiter.OnRateLimitHit()
	reducedRate := limiter.CurrentRate()

	// Simulate success
	limiter.OnSuccess()
	recoveredRate := limiter.CurrentRate()

	// Rate should increase by 10%
	expectedRate := reducedRate * 1.1
	if recoveredRate != expectedRate {
		t.Errorf("Expected rate %f after recovery, got %f", expectedRate, recoveredRate)
	}

	// Multiple successes should not exceed base rate
	for i := 0; i < 20; i++ {
		limiter.OnSuccess()
	}

	finalRate := limiter.CurrentRate()
	if finalRate > 30.0 {
		t.Errorf("Rate exceeded base rate: %f", finalRate)
	}
}

func TestDynamicRateLimiter_SetRate(t *testing.T) {
	logger := zap.NewNop()
	limiter := NewDynamicRateLimiter(30, logger)

	// Set custom rate
	limiter.SetRate(20.0)

	rate := limiter.CurrentRate()
	if rate != 20.0 {
		t.Errorf("Expected rate 20, got %f", rate)
	}

	// Should not allow below minimum
	limiter.SetRate(3.0)
	rate = limiter.CurrentRate()
	if rate < 5.0 {
		t.Errorf("Rate went below minimum: %f", rate)
	}

	// Should not allow above base rate
	limiter.SetRate(50.0)
	rate = limiter.CurrentRate()
	if rate > 30.0 {
		t.Errorf("Rate exceeded base rate: %f", rate)
	}
}

func TestDynamicRateLimiter_ContextCancellation(t *testing.T) {
	logger := zap.NewNop()
	limiter := NewDynamicRateLimiter(1, logger) // Very slow rate

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// First request succeeds
	err := limiter.Wait(ctx)
	if err != nil {
		t.Fatalf("First request failed: %v", err)
	}

	// Second request should timeout
	err = limiter.Wait(ctx)
	if err == nil {
		t.Error("Expected context timeout error")
	}
}

func BenchmarkDynamicRateLimiter_Wait(b *testing.B) {
	logger := zap.NewNop()
	limiter := NewDynamicRateLimiter(1000, logger) // High rate for benchmark
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.Wait(ctx)
	}
}
