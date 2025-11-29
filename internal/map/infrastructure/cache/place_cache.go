package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// CachedPlace represents a cached place details
type CachedPlace struct {
	PlaceID      string                 `json:"place_id"`
	LanguageCode string                 `json:"language_code"`
	Data         map[string]interface{} `json:"data"`
	CachedAt     time.Time              `json:"cached_at"`
	ExpiresAt    time.Time              `json:"expires_at"`
}

// PlaceCache manages caching of place details
type PlaceCache struct {
	redis  *redis.Client
	logger *zap.Logger
}

// NewPlaceCache creates a new PlaceCache
func NewPlaceCache(redis *redis.Client, logger *zap.Logger) *PlaceCache {
	return &PlaceCache{
		redis:  redis,
		logger: logger.With(zap.String("component", "place_cache")),
	}
}

// GetPlace retrieves a place from cache
func (c *PlaceCache) GetPlace(
	ctx context.Context,
	placeID string,
	languageCode string,
) (*CachedPlace, error) {
	key := c.buildKey(placeID, languageCode)

	val, err := c.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		c.logger.Debug("Cache miss",
			zap.String("place_id", placeID),
			zap.String("language", languageCode),
		)
		return nil, nil // Cache miss
	} else if err != nil {
		c.logger.Error("Redis get failed", zap.Error(err))
		return nil, fmt.Errorf("redis get failed: %w", err)
	}

	var cached CachedPlace
	if err := json.Unmarshal([]byte(val), &cached); err != nil {
		c.logger.Error("Failed to unmarshal cached data", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	c.logger.Info("Cache hit",
		zap.String("place_id", placeID),
		zap.String("language", languageCode),
		zap.Time("cached_at", cached.CachedAt),
	)

	return &cached, nil
}

// SetPlace stores a place in cache
func (c *PlaceCache) SetPlace(
	ctx context.Context,
	placeID string,
	languageCode string,
	data map[string]interface{},
	ttl time.Duration,
) error {
	key := c.buildKey(placeID, languageCode)

	cached := CachedPlace{
		PlaceID:      placeID,
		LanguageCode: languageCode,
		Data:         data,
		CachedAt:     time.Now(),
		ExpiresAt:    time.Now().Add(ttl),
	}

	jsonData, err := json.Marshal(cached)
	if err != nil {
		c.logger.Error("Failed to marshal cache data", zap.Error(err))
		return fmt.Errorf("failed to marshal: %w", err)
	}

	if err := c.redis.Set(ctx, key, jsonData, ttl).Err(); err != nil {
		c.logger.Error("Redis set failed", zap.Error(err))
		return fmt.Errorf("redis set failed: %w", err)
	}

	c.logger.Info("Cached place details",
		zap.String("place_id", placeID),
		zap.String("language", languageCode),
		zap.Duration("ttl", ttl),
	)

	return nil
}

// buildKey creates a cache key
func (c *PlaceCache) buildKey(placeID, languageCode string) string {
	return fmt.Sprintf("map:place:%s:%s", placeID, languageCode)
}
