package application

import (
	pkgconfig "github.com/Leon180/tabelogo-v2/pkg/config"
	"go.uber.org/fx"
)

// NewConfig creates application config from pkg config
func NewConfig(cfg *pkgconfig.Config) *Config {
	return &Config{
		DataFreshnessTTL: cfg.MapService.DataFreshnessTTL,
	}
}

// Module provides application layer dependencies
var Module = fx.Module("restaurant.application",
	fx.Provide(
		NewConfig,
		NewRestaurantService,
	),
)
