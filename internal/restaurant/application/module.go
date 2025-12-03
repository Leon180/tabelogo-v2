package application

import (
	"go.uber.org/fx"
)

// Module provides application layer dependencies
var Module = fx.Module("restaurant.application",
	fx.Provide(
		NewRestaurantService,
	),
)
