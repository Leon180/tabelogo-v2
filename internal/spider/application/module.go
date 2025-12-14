package application

import (
	"context"

	"github.com/Leon180/tabelogo-v2/internal/spider/application/services"
	"github.com/Leon180/tabelogo-v2/internal/spider/application/usecases"
	"go.uber.org/fx"
)

// Module provides application layer dependencies
var Module = fx.Module("application",
	// Services
	fx.Provide(
		services.NewJobProcessor,
		services.NewRateLimiter,
	),

	// Use cases
	fx.Provide(
		usecases.NewScrapeRestaurantUseCase,
		usecases.NewGetJobStatusUseCase,

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
