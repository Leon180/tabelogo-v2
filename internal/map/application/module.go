package application

import (
	"github.com/Leon180/tabelogo-v2/internal/map/application/usecases"
	"go.uber.org/fx"
)

// Module provides application layer dependencies
var Module = fx.Module("map.application",
	fx.Provide(
		usecases.NewQuickSearchUseCase,
		usecases.NewAdvanceSearchUseCase,
	),
)
