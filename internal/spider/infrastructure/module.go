package infrastructure

import (
	"context"
	"fmt"

	"github.com/Leon180/tabelogo-v2/internal/spider/domain/repositories"
	"github.com/Leon180/tabelogo-v2/internal/spider/infrastructure/persistence"
	"github.com/Leon180/tabelogo-v2/internal/spider/infrastructure/scraper"
	"github.com/Leon180/tabelogo-v2/pkg/config"
	redisclient "github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides infrastructure layer dependencies
var Module = fx.Module("spider.infrastructure",
	fx.Provide(
		NewRedis,
		// Repositories
		fx.Annotate(
			persistence.NewInMemoryJobRepository,
			fx.As(new(repositories.JobRepository)),
		),
		fx.Annotate(
			persistence.NewRedisResultCache,
			fx.As(new(repositories.ResultCacheRepository)),
		),
		// Scraper
		scraper.NewScraper,
	),
)

// NewRedis creates a new Redis client with lifecycle management
// Spider Service uses Redis DB 2
func NewRedis(cfg *config.Config, lc fx.Lifecycle, logger *zap.Logger) *redisclient.Client {
	rdb := redisclient.NewClient(&redisclient.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.Redis.Password,
		DB:       2, // Use DB 2 for Spider Service
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Connecting to Redis", zap.String("addr", cfg.GetRedisAddr()), zap.Int("db", 2))
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
