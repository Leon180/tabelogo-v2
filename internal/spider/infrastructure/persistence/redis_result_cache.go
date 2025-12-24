package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/Leon180/tabelogo-v2/internal/spider/domain/models"
	"github.com/Leon180/tabelogo-v2/internal/spider/domain/repositories"
	redisclient "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisResultCache implements ResultCacheRepository using Redis
type RedisResultCache struct {
	client   *redisclient.Client
	helper   *RedisHelper // Added
	logger   *zap.Logger
	cacheTTL time.Duration
}

// NewRedisResultCache creates a new Redis result cache
func NewRedisResultCache(client *redisclient.Client, cacheTTL time.Duration, logger *zap.Logger) repositories.ResultCacheRepository { // Signature changed
	return &RedisResultCache{
		client:   client,
		helper:   NewRedisHelper(client, logger), // Added
		logger:   logger.With(zap.String("component", "redis_result_cache")),
		cacheTTL: cacheTTL,
	}
}

// cacheKey generates a Redis key for a given Google Place ID.
func (r *RedisResultCache) cacheKey(googleID string) string {
	return fmt.Sprintf("tabelog:results:%s", googleID)
}

// Get retrieves cached results by Google Place ID
func (r *RedisResultCache) Get(ctx context.Context, googleID string) (*models.CachedResult, error) {
	key := r.cacheKey(googleID)

	var cached models.CachedResult
	if err := r.helper.GetJSON(ctx, key, &cached); err != nil {
		// Not found is not an error, just return nil
		if err.Error() == fmt.Sprintf("key not found: %s", key) {
			return nil, nil
		}
		r.logger.Error("Failed to get cached result", zap.Error(err), zap.String("google_id", googleID))
		return nil, err
	}

	// Check if expired
	if cached.IsExpired() {
		r.logger.Info("Cached result expired", zap.String("google_id", googleID))
		r.helper.Delete(ctx, key)
		return nil, nil
	}

	r.logger.Info("Cache hit", zap.String("google_id", googleID))
	return &cached, nil
}

// Set stores cached results
func (r *RedisResultCache) Set(ctx context.Context, googleID string, restaurants []models.TabelogRestaurant, ttl time.Duration) error {
	key := r.cacheKey(googleID)

	// Convert domain models to DTOs for JSON serialization
	dtos := make([]models.TabelogRestaurantDTO, len(restaurants))
	for i, restaurant := range restaurants {
		dtos[i] = restaurant.ToDTO()
	}

	cached := &models.CachedResult{
		PlaceID:   googleID,
		Results:   dtos,
		CachedAt:  time.Now(),
		ExpiresAt: time.Now().Add(ttl),
	}

	if err := r.helper.SetJSON(ctx, key, cached, ttl); err != nil {
		r.logger.Error("Failed to set cached result", zap.Error(err), zap.String("google_id", googleID))
		return err
	}

	r.logger.Info("Cached result stored",
		zap.String("google_id", googleID),
		zap.Int("restaurant_count", len(restaurants)),
		zap.Duration("ttl", ttl),
	)

	return nil
}

// Delete removes cached results
func (r *RedisResultCache) Delete(ctx context.Context, googleID string) error {
	key := r.cacheKey(googleID)

	if err := r.helper.Delete(ctx, key); err != nil {
		r.logger.Error("Failed to delete cached result", zap.Error(err), zap.String("google_id", googleID))
		return err
	}

	r.logger.Info("Cached result deleted", zap.String("google_id", googleID))
	return nil
}
