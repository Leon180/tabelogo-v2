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
	),

	// Lifecycle hooks for job processor
	fx.Invoke(registerJobProcessorLifecycle),
)

// registerJobProcessorLifecycle registers lifecycle hooks for the job processor
func registerJobProcessorLifecycle(lc fx.Lifecycle, processor *services.JobProcessor) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Start the job processor with a background context
			// The processor runs continuously and manages its own lifecycle
			go processor.Start(ctx)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Stop the processor gracefully
			return processor.Stop(ctx)
		},
	})
}
