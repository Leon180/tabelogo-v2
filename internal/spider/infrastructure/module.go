package infrastructure

import (
	"time"

	"github.com/Leon180/tabelogo-v2/internal/spider/config"
	"github.com/Leon180/tabelogo-v2/internal/spider/domain/models"
	"github.com/Leon180/tabelogo-v2/internal/spider/domain/repositories"
	"github.com/Leon180/tabelogo-v2/internal/spider/infrastructure/metrics"
	"github.com/Leon180/tabelogo-v2/internal/spider/infrastructure/persistence"
	"github.com/Leon180/tabelogo-v2/internal/spider/infrastructure/scraper"
	pkgconfig "github.com/Leon180/tabelogo-v2/pkg/config"
	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides infrastructure layer dependencies
var Module = fx.Module("spider.infrastructure",
	// Configuration
	fx.Provide(config.DefaultConfig),

	// Metrics
	fx.Provide(metrics.NewSpiderMetrics),

	// Redis client
	fx.Provide(newRedisClient),

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

// newRedisClient creates a Redis client from main config
func newRedisClient(cfg *pkgconfig.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
}

// newRedisResultCache creates a Redis result cache with configured TTL
func newRedisResultCache(client *redis.Client, logger *zap.Logger, cfg *config.SpiderConfig) repositories.ResultCacheRepository {
	return persistence.NewRedisResultCache(client, cfg.CacheTTL, logger)
}

// newCircuitBreaker creates a circuit breaker with configured settings
func newCircuitBreaker(logger *zap.Logger, m *metrics.SpiderMetrics, cfg *config.SpiderConfig) *gobreaker.CircuitBreaker {
	cbConfig := scraper.CircuitBreakerConfig{
		MaxRequests: cfg.CircuitBreaker.MaxRequests,
		Interval:    cfg.CircuitBreaker.Interval,
		Timeout:     cfg.CircuitBreaker.Timeout,
	}
	return scraper.NewCircuitBreaker(logger, m, cbConfig)
}

// newScraper creates a scraper with dependencies
func newScraper(logger *zap.Logger, m *metrics.SpiderMetrics, cb *gobreaker.CircuitBreaker) *scraper.Scraper {
	scraperConfig := models.NewScraperConfig().
		WithTimeout(30 * time.Second)
	return scraper.NewScraper(logger, m, scraperConfig, cb)
}
