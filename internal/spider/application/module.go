package application

import (
	"context"

	"github.com/Leon180/tabelogo-v2/internal/spider/application/services"
	"github.com/Leon180/tabelogo-v2/internal/spider/application/usecases"
	"github.com/Leon180/tabelogo-v2/internal/spider/domain/repositories"
	"github.com/Leon180/tabelogo-v2/internal/spider/infrastructure/metrics"
	"github.com/Leon180/tabelogo-v2/internal/spider/infrastructure/scraper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides application layer dependencies
var Module = fx.Module("spider.application",
	fx.Provide(
		usecases.NewScrapeRestaurantUseCase,
		usecases.NewGetJobStatusUseCase,
		NewJobProcessor,
	),
	fx.Invoke(StartJobProcessor),
)

// NewJobProcessor creates a job processor with 20 workers
func NewJobProcessor(
	jobRepo repositories.JobRepository,
	resultCache repositories.ResultCacheRepository,
	scraper *scraper.Scraper,
	metrics *metrics.SpiderMetrics,
	logger *zap.Logger,
) *services.JobProcessor {
	return services.NewJobProcessor(jobRepo, resultCache, scraper, metrics, logger, 20)
}

// StartJobProcessor starts the job processor with lifecycle management
func StartJobProcessor(lc fx.Lifecycle, processor *services.JobProcessor, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting job processor")
			// Use background context for long-running workers
			// The hook ctx is only for startup timeout, not for worker lifetime
			processor.Start(context.Background())
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping job processor")
			return processor.Stop(ctx)
		},
	})
}
