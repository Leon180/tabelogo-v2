package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Leon180/tabelogo-v2/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimitConfig holds the configuration for rate limiting
type RateLimitConfig struct {
	// RedisClient is the Redis client for distributed rate limiting
	RedisClient *redis.Client
	// Limit is the maximum number of requests allowed
	Limit int
	// Window is the time window for rate limiting
	Window time.Duration
	// KeyPrefix is the prefix for Redis keys
	KeyPrefix string
	// SkipPaths are paths that skip rate limiting
	SkipPaths []string
}

// RateLimit returns a middleware that limits request rate using Redis
func RateLimit(config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if path should skip rate limiting
		for _, path := range config.SkipPaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		// Get client identifier (IP or user ID)
		clientID := getClientIdentifier(c)

		// Create Redis key
		key := fmt.Sprintf("%s:%s", config.KeyPrefix, clientID)

		// Check rate limit
		allowed, remaining, err := checkRateLimit(
			c.Request.Context(),
			config.RedisClient,
			key,
			config.Limit,
			config.Window,
		)

		if err != nil {
			// Log error but don't block request on Redis failure
			c.Next()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", config.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(config.Window).Unix()))

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    errors.ErrCodeRateLimitExceeded,
				"message": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkRateLimit checks if the client has exceeded the rate limit using sliding window algorithm
func checkRateLimit(
	ctx context.Context,
	client *redis.Client,
	key string,
	limit int,
	window time.Duration,
) (allowed bool, remaining int, err error) {
	now := time.Now()
	windowStart := now.Add(-window)

	pipe := client.Pipeline()

	// Remove old entries outside the window
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart.UnixNano()))

	// Count current requests in window
	countCmd := pipe.ZCard(ctx, key)

	// Add current request
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(now.UnixNano()),
		Member: fmt.Sprintf("%d", now.UnixNano()),
	})

	// Set expiration
	pipe.Expire(ctx, key, window)

	// Execute pipeline
	_, err = pipe.Exec(ctx)
	if err != nil {
		return false, 0, err
	}

	count := int(countCmd.Val())

	// Check if limit exceeded
	if count >= limit {
		return false, 0, nil
	}

	return true, limit - count - 1, nil
}

// getClientIdentifier returns a unique identifier for the client
func getClientIdentifier(c *gin.Context) string {
	// Try to get user ID from context (for authenticated users)
	if userID, exists := c.Get(UserIDKey); exists {
		if id, ok := userID.(string); ok {
			return fmt.Sprintf("user:%s", id)
		}
	}

	// Fall back to IP address
	return fmt.Sprintf("ip:%s", c.ClientIP())
}

// InMemoryRateLimit returns a simple in-memory rate limiter (for development)
// Note: This is not distributed and should only be used for testing
func InMemoryRateLimit(limit int, window time.Duration) gin.HandlerFunc {
	type clientInfo struct {
		requests []time.Time
	}

	clients := make(map[string]*clientInfo)

	return func(c *gin.Context) {
		clientID := getClientIdentifier(c)
		now := time.Now()

		// Get or create client info
		info, exists := clients[clientID]
		if !exists {
			info = &clientInfo{requests: []time.Time{}}
			clients[clientID] = info
		}

		// Remove old requests outside the window
		validRequests := []time.Time{}
		for _, reqTime := range info.requests {
			if now.Sub(reqTime) <= window {
				validRequests = append(validRequests, reqTime)
			}
		}
		info.requests = validRequests

		// Check rate limit
		if len(info.requests) >= limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    errors.ErrCodeRateLimitExceeded,
				"message": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		// Add current request
		info.requests = append(info.requests, now)

		// Set rate limit headers
		remaining := limit - len(info.requests)
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", now.Add(window).Unix()))

		c.Next()
	}
}
