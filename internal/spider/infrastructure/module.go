package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/Leon180/tabelogo-v2/internal/spider/config"
	"github.com/Leon180/tabelogo-v2/internal/spider/domain/repositories"
	"github.com/Leon180/tabelogo-v2/internal/spider/infrastructure/metrics"
	"github.com/Leon180/tabelogo-v2/internal/spider/infrastructure/persistence"
	"github.com/Leon180/tabelogo-v2/internal/spider/infrastructure/scraper"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides infrastructure layer dependencies
var Module = fx.Module("spider.infrastructure",
	// Configuration
	fx.Provide(config.DefaultConfig),

	// Metrics
	fx.Provide(metrics.NewSpiderMetrics),

	// Persistence
	fx.Provide(
		fx.Annotate(
			persistence.NewRedisJobStore,
			fx.As(new(repositories.JobRepository)),
		),
		fx.Annotate(
			newRedisResultCache,
			fx.As(new(repositories.ResultCacheRepository)),
		),
	),

	// Scraper with circuit breaker
	fx.Provide(
		newCircuitBreaker,
		newScraper,
	),
)

// NewRedis creates a new Redis client and manages its lifecycle with Fx.
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

// newRedisResultCache creates a Redis result cache with configured TTL
func newRedisResultCache(client *redis.Client, logger *zap.Logger, cfg *config.SpiderConfig) repositories.ResultCacheRepository {
	return persistence.NewRedisResultCache(client, logger, cfg.CacheTTL)
}

// newCircuitBreaker creates a circuit breaker with configured settings
func newCircuitBreaker(logger *zap.Logger, metrics *metrics.SpiderMetrics, cfg *config.SpiderConfig) *scraper.CircuitBreaker {
	cbConfig := scraper.CircuitBreakerConfig{
		MaxRequests: cfg.CircuitBreaker.MaxRequests,
		Interval:    cfg.CircuitBreaker.Interval,
		Timeout:     cfg.CircuitBreaker.Timeout,
	}
	return scraper.NewCircuitBreaker(logger, metrics, cbConfig)
}

// newScraper creates a scraper with dependencies
func newScraper(logger *zap.Logger, metrics *metrics.SpiderMetrics, cb *scraper.CircuitBreaker) *scraper.Scraper {
	scraperConfig := &scraper.ScraperConfig{
		UserAgent:      "Mozilla/5.0 (compatible; TabelogoBot/1.0)",
		RequestTimeout: 30 * time.Second,
	}
	return scraper.NewScraper(logger, metrics, scraperConfig, cb)
}
