package application

import (
	"github.com/Leon180/tabelogo-v2/internal/spider/application/usecases"
	"go.uber.org/fx"
)

// Module provides application layer dependencies
var Module = fx.Module("spider.application",
	fx.Provide(
		usecases.NewScrapeRestaurantUseCase,
		usecases.NewGetJobStatusUseCase,
	),
)
