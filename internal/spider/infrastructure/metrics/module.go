package metrics

import (
	"go.uber.org/fx"
)

// Module provides metrics dependencies
var Module = fx.Module("spider.metrics",
	fx.Provide(
		NewSpiderMetrics,
	),
)
