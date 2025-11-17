package config

import (
	"go.uber.org/fx"
)

// Module provides config for FX dependency injection
// Usage:
//   fx.New(
//     config.Module,
//     fx.Invoke(func(cfg *config.Config) { ... }),
//   )
var Module = fx.Module("config",
	fx.Provide(New),
)

// New creates a new config instance for dependency injection
// This is the FX-compatible constructor
func New() (*Config, error) {
	return Load()
}

// NewWithPrefix creates a new config with environment variable prefix
func NewWithPrefix(prefix string) func() (*Config, error) {
	return func() (*Config, error) {
		return LoadWithPrefix(prefix)
	}
}
