package infrastructure

import (
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

// newRedisResultCache creates a Redis result cache with configured TTL
func newRedisResultCache(client *redis.Client, logger *zap.Logger, cfg *config.SpiderConfig) repositories.ResultCacheRepository {
	return persistence.NewRedisResultCache(client, logger, cfg.CacheTTL)
}

// newCircuitBreaker creates a circuit breaker with configured settings
func newCircuitBreaker(logger *zap.Logger, m *metrics.SpiderMetrics, cfg *config.SpiderConfig) *scraper.CircuitBreaker {
	cbConfig := scraper.CircuitBreakerConfig{
		MaxRequests: cfg.CircuitBreaker.MaxRequests,
		Interval:    cfg.CircuitBreaker.Interval,
		Timeout:     cfg.CircuitBreaker.Timeout,
	}
	return scraper.NewCircuitBreaker(logger, m, cbConfig)
}

// newScraper creates a scraper with dependencies
func newScraper(logger *zap.Logger, m *metrics.SpiderMetrics, cb *scraper.CircuitBreaker) *scraper.Scraper {
	scraperConfig := scraper.ScraperConfig{
		UserAgent:      "Mozilla/5.0 (compatible; TabelogoBot/1.0)",
		RequestTimeout: 30 * time.Second,
	}
	return scraper.NewScraper(logger, m, &scraperConfig, cb)
}
