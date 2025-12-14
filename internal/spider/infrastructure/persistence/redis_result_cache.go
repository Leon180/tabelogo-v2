package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Leon180/tabelogo-v2/internal/spider/domain/models"
	redisclient "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisResultCache implements ResultCacheRepository using Redis
type RedisResultCache struct {
	client   *redisclient.Client
	logger   *zap.Logger
	cacheTTL time.Duration
}

// NewRedisResultCache creates a new Redis-based result cache
func NewRedisResultCache(client *redisclient.Client, logger *zap.Logger, cacheTTL time.Duration) *RedisResultCache {
	return &RedisResultCache{
		client:   client,
		logger:   logger.With(zap.String("component", "redis_result_cache")),
		cacheTTL: cacheTTL,
	}
}

// Get retrieves cached results for a place
func (r *RedisResultCache) Get(ctx context.Context, placeID string) (*models.CachedResult, error) {
	key := fmt.Sprintf("tabelog:results:%s", placeID)

	data, err := r.client.Get(ctx, key).Bytes()
	if err == redisclient.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		r.logger.Error("Failed to get cached results", zap.Error(err), zap.String("place_id", placeID))
		return nil, fmt.Errorf("failed to get cached results: %w", err)
	}

	var cached models.CachedResult
	if err := json.Unmarshal(data, &cached); err != nil {
		r.logger.Error("Failed to unmarshal cached results", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal cached results: %w", err)
	}

	// Check if expired
	if cached.IsExpired() {
		r.logger.Info("Cached results expired", zap.String("place_id", placeID))
		_ = r.Delete(ctx, placeID)
		return nil, nil
	}

	r.logger.Info("Cache hit", zap.String("place_id", placeID), zap.Int("results_count", len(cached.Results)))
	return &cached, nil
}

// Set stores results in cache with TTL
func (r *RedisResultCache) Set(ctx context.Context, placeID string, results []models.TabelogRestaurant, ttl time.Duration) error {
	key := fmt.Sprintf("tabelog:results:%s", placeID)

	// Convert domain models to DTOs for JSON serialization
	dtos := make([]models.TabelogRestaurantDTO, len(results))
	for i, restaurant := range results {
		dtos[i] = restaurant.ToDTO()
	}

	cached := models.CachedResult{
		PlaceID:   placeID,
		Results:   dtos,
		CachedAt:  time.Now(),
		ExpiresAt: time.Now().Add(r.cacheTTL),
	}

	data, err := json.Marshal(cached)
	if err != nil {
		r.logger.Error("Failed to marshal cached results", zap.Error(err))
		return fmt.Errorf("failed to marshal cached results: %w", err)
	}

	// Store in Redis with TTL
	err = r.client.Set(ctx, key, data, r.cacheTTL).Err()
	if err != nil {
		r.logger.Error("Failed to cache results",
			zap.String("place_id", placeID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to cache results: %w", err)
	}

	r.logger.Info("Cached results",
		zap.String("place_id", placeID),
		zap.Int("results_count", len(results)),
		zap.Duration("ttl", r.cacheTTL),
	)

	return nil
}

// Delete removes cached results
func (r *RedisResultCache) Delete(ctx context.Context, placeID string) error {
	key := fmt.Sprintf("tabelog:results:%s", placeID)

	if err := r.client.Del(ctx, key).Err(); err != nil {
		r.logger.Error("Failed to delete cached results", zap.Error(err), zap.String("place_id", placeID))
		return fmt.Errorf("failed to delete cached results: %w", err)
	}

	r.logger.Info("Deleted cached results", zap.String("place_id", placeID))
	return nil
}
