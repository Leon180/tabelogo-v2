package application

import (
	"context"

	"github.com/Leon180/tabelogo-v2/internal/spider/application/services"
	"github.com/Leon180/tabelogo-v2/internal/spider/application/usecases"
	"github.com/Leon180/tabelogo-v2/internal/spider/config"
	"github.com/Leon180/tabelogo-v2/internal/spider/domain/repositories"
	"github.com/Leon180/tabelogo-v2/internal/spider/infrastructure/metrics"
	"github.com/Leon180/tabelogo-v2/internal/spider/infrastructure/scraper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// newJobProcessor creates a JobProcessor with workerCount from config
func newJobProcessor(
	jobRepo repositories.JobRepository,
	resultCache repositories.ResultCacheRepository,
	scraperInstance *scraper.Scraper,
	metrics *metrics.SpiderMetrics,
	logger *zap.Logger,
	cfg *config.SpiderConfig,
) *services.JobProcessor {
	return services.NewJobProcessor(
		jobRepo,
		resultCache,
		scraperInstance,
		metrics,
		logger,
		cfg.WorkerCount,
	)
}

// Module provides application layer dependencies
var Module = fx.Module("application",
	// Services
	fx.Provide(
		newJobProcessor, // Use our custom provider that injects workerCount
	),

	// Use cases
	fx.Provide(
		usecases.NewScrapeRestaurantUseCase,
		usecases.NewGetJobStatusUseCase,
	),

	// Lifecycle hooks for job processor
	fx.Invoke(registerJobProcessorLifecycle),
)

// registerJobProcessorLifecycle registers lifecycle hooks for the job processor
func registerJobProcessorLifecycle(lc fx.Lifecycle, processor *services.JobProcessor) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Start the job processor with a background context
			// We use context.Background() instead of ctx because the lifecycle context
			// gets cancelled after startup, which would cause all workers to exit
			go processor.Start(context.Background())
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Stop the processor gracefully
			return processor.Stop(ctx)
		},
	})
}
