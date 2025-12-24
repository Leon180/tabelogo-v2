package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redisclient "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisHelper provides common Redis operations with JSON marshaling
type RedisHelper struct {
	client *redisclient.Client
	logger *zap.Logger
}

// NewRedisHelper creates a new Redis helper
func NewRedisHelper(client *redisclient.Client, logger *zap.Logger) *RedisHelper {
	return &RedisHelper{
		client: client,
		logger: logger.With(zap.String("component", "redis_helper")),
	}
}

// SetJSON marshals value to JSON and stores it in Redis with TTL
func (h *RedisHelper) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		h.logger.Error("Failed to marshal value",
			zap.Error(err),
			zap.String("key", key))
		return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
	}

	if err := h.client.Set(ctx, key, data, ttl).Err(); err != nil {
		h.logger.Error("Failed to set value in Redis",
			zap.Error(err),
			zap.String("key", key))
		return fmt.Errorf("failed to set value for key %s: %w", key, err)
	}

	return nil
}

// GetJSON retrieves value from Redis and unmarshals it from JSON
func (h *RedisHelper) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := h.client.Get(ctx, key).Bytes()
	if err == redisclient.Nil {
		return fmt.Errorf("key not found: %s", key)
	}
	if err != nil {
		h.logger.Error("Failed to get value from Redis",
			zap.Error(err),
			zap.String("key", key))
		return fmt.Errorf("failed to get value for key %s: %w", key, err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		h.logger.Error("Failed to unmarshal value",
			zap.Error(err),
			zap.String("key", key))
		return fmt.Errorf("failed to unmarshal value for key %s: %w", key, err)
	}

	return nil
}

// Delete removes a key from Redis
func (h *RedisHelper) Delete(ctx context.Context, key string) error {
	if err := h.client.Del(ctx, key).Err(); err != nil {
		h.logger.Error("Failed to delete key from Redis",
			zap.Error(err),
			zap.String("key", key))
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

// Exists checks if a key exists in Redis
func (h *RedisHelper) Exists(ctx context.Context, key string) (bool, error) {
	count, err := h.client.Exists(ctx, key).Result()
	if err != nil {
		h.logger.Error("Failed to check key existence",
			zap.Error(err),
			zap.String("key", key))
		return false, fmt.Errorf("failed to check existence of key %s: %w", key, err)
	}
	return count > 0, nil
}

// SetAdd adds members to a Redis set
func (h *RedisHelper) SetAdd(ctx context.Context, key string, members ...interface{}) error {
	if err := h.client.SAdd(ctx, key, members...).Err(); err != nil {
		h.logger.Error("Failed to add to set",
			zap.Error(err),
			zap.String("key", key))
		return fmt.Errorf("failed to add to set %s: %w", key, err)
	}
	return nil
}

// SetRemove removes members from a Redis set
func (h *RedisHelper) SetRemove(ctx context.Context, key string, members ...interface{}) error {
	if err := h.client.SRem(ctx, key, members...).Err(); err != nil {
		h.logger.Error("Failed to remove from set",
			zap.Error(err),
			zap.String("key", key))
		return fmt.Errorf("failed to remove from set %s: %w", key, err)
	}
	return nil
}

// SetMembers retrieves all members of a Redis set
func (h *RedisHelper) SetMembers(ctx context.Context, key string) ([]string, error) {
	members, err := h.client.SMembers(ctx, key).Result()
	if err != nil {
		h.logger.Error("Failed to get set members",
			zap.Error(err),
			zap.String("key", key))
		return nil, fmt.Errorf("failed to get members of set %s: %w", key, err)
	}
	return members, nil
}

// Expire sets a TTL on a key
func (h *RedisHelper) Expire(ctx context.Context, key string, ttl time.Duration) error {
	if err := h.client.Expire(ctx, key, ttl).Err(); err != nil {
		h.logger.Error("Failed to set expiration",
			zap.Error(err),
			zap.String("key", key))
		return fmt.Errorf("failed to set expiration for key %s: %w", key, err)
	}
	return nil
}
