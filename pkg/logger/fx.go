package logger

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides logger for FX dependency injection
// Usage:
//   fx.New(
//     logger.Module,
//     fx.Invoke(func(log *zap.Logger) { ... }),
//   )
var Module = fx.Module("logger",
	fx.Provide(New),
)

// Params holds dependencies for logger
type Params struct {
	fx.In

	// Optional: can provide LogLevel from config
	LogLevel string `optional:"true"`
}

// New creates a new logger instance for dependency injection
// This is the FX-compatible constructor
func New(params Params) (*zap.Logger, error) {
	level := params.LogLevel
	if level == "" {
		level = "info"
	}

	// Check if development mode (can be passed via params)
	if level == "debug" {
		return zap.NewDevelopment()
	}

	// Production mode
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(parseLevel(level))

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return logger, nil
}

// NewDevelopment creates a development logger for FX
func NewDevelopment() (*zap.Logger, error) {
	return zap.NewDevelopment()
}

// NewProduction creates a production logger for FX
func NewProduction() (*zap.Logger, error) {
	return zap.NewProduction()
}
