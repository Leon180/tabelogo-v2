package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
	"go.uber.org/zap"
)

// RateLimiterConfig holds rate limiter configuration
type RateLimiterConfig struct {
	Rate   limiter.Rate
	Logger *zap.Logger
}

// NewRateLimiter creates a new rate limiter middleware
func NewRateLimiter(redisClient *redis.Client, config RateLimiterConfig) gin.HandlerFunc {
	// Create Redis store
	store, err := sredis.NewStoreWithOptions(redisClient, limiter.StoreOptions{
		Prefix:   "rate_limit",
		MaxRetry: 3,
	})
	if err != nil {
		config.Logger.Fatal("Failed to create rate limiter store", zap.Error(err))
	}

	// Create rate limiter
	rate := config.Rate
	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		// Get client IP
		key := c.ClientIP()

		config.Logger.Debug("Rate limiter checking",
			zap.String("ip", key),
			zap.String("path", c.Request.URL.Path),
		)

		// Get limit context
		context, err := instance.Get(c.Request.Context(), key)
		if err != nil {
			config.Logger.Error("Rate limiter error", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "rate_limiter_error",
				"message": "Internal server error",
			})
			c.Abort()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", context.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", context.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", context.Reset))

		config.Logger.Debug("Rate limit headers set",
			zap.Int64("limit", context.Limit),
			zap.Int64("remaining", context.Remaining),
			zap.Bool("reached", context.Reached),
		)

		// Check if limit exceeded
		if context.Reached {
			config.Logger.Warn("Rate limit exceeded",
				zap.String("ip", key),
				zap.Int64("limit", context.Limit),
			)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate_limit_exceeded",
				"message":     "Too many requests, please try again later",
				"retry_after": context.Reset,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
