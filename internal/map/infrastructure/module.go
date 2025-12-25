package infrastructure

import (
	"context"
	"fmt"

	"github.com/Leon180/tabelogo-v2/internal/map/infrastructure/cache"
	"github.com/Leon180/tabelogo-v2/internal/map/infrastructure/external"
	"github.com/Leon180/tabelogo-v2/pkg/config"
	redisclient "github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides infrastructure dependencies
var Module = fx.Module("map.infrastructure",
	fx.Provide(
		NewRedis,
		cache.NewPlaceCache,
		external.NewPlacesClient, // Factory that returns either Google or Mock client
	),
)

// NewRedis creates a new Redis client with lifecycle management
func NewRedis(cfg *config.Config, lc fx.Lifecycle, logger *zap.Logger) *redisclient.Client {
	rdb := redisclient.NewClient(&redisclient.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Connecting to Redis", zap.String("addr", cfg.GetRedisAddr()), zap.Int("db", cfg.Redis.DB))
			if err := rdb.Ping(ctx).Err(); err != nil {
				return fmt.Errorf("failed to connect to Redis: %w", err)
			}
			logger.Info("Redis connected successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing Redis connection")
			return rdb.Close()
		},
	})

	return rdb
}
