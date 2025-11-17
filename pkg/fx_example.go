package pkg

// This file demonstrates how to use pkg with Uber FX dependency injection
// DO NOT include this file in production builds

/*

import (
	"context"

	"github.com/Leon180/tabelogo-v2/pkg/config"
	"github.com/Leon180/tabelogo-v2/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Example 1: Basic FX Application
func ExampleBasicApp() {
	app := fx.New(
		// Provide config
		config.Module,

		// Provide logger
		logger.Module,

		// Use them in your service
		fx.Invoke(func(cfg *config.Config, log *zap.Logger) {
			log.Info("Application starting",
				zap.String("environment", cfg.Environment),
				zap.Int("port", cfg.ServerPort),
			)
		}),
	)

	app.Run()
}

// Example 2: Custom Service with Dependencies
type MyService struct {
	config *config.Config
	logger *zap.Logger
}

func NewMyService(cfg *config.Config, log *zap.Logger) *MyService {
	return &MyService{
		config: cfg,
		logger: log,
	}
}

func (s *MyService) Start(ctx context.Context) error {
	s.logger.Info("Service starting")
	return nil
}

func (s *MyService) Stop(ctx context.Context) error {
	s.logger.Info("Service stopping")
	return nil
}

func ExampleServiceWithLifecycle() {
	app := fx.New(
		// Modules
		config.Module,
		logger.Module,

		// Provide your service
		fx.Provide(NewMyService),

		// Hook into lifecycle
		fx.Invoke(func(lc fx.Lifecycle, svc *MyService) {
			lc.Append(fx.Hook{
				OnStart: svc.Start,
				OnStop:  svc.Stop,
			})
		}),
	)

	app.Run()
}

// Example 3: Multiple Services with Dependency Graph
type DatabaseService struct {
	config *config.Config
	logger *zap.Logger
}

func NewDatabaseService(cfg *config.Config, log *zap.Logger) *DatabaseService {
	log.Info("Creating database service", zap.String("dsn", cfg.GetDatabaseDSN()))
	return &DatabaseService{config: cfg, logger: log}
}

type APIService struct {
	db     *DatabaseService
	logger *zap.Logger
}

func NewAPIService(db *DatabaseService, log *zap.Logger) *APIService {
	return &APIService{db: db, logger: log}
}

func ExampleDependencyGraph() {
	app := fx.New(
		// Base modules
		config.Module,
		logger.Module,

		// Services (FX will resolve dependency order automatically)
		fx.Provide(
			NewDatabaseService,
			NewAPIService,
		),

		// Invoke to start
		fx.Invoke(func(api *APIService) {
			api.logger.Info("API service ready")
		}),
	)

	app.Run()
}

// Example 4: Custom Logger Configuration
func ExampleCustomLoggerConfig() {
	app := fx.New(
		// Load config first
		config.Module,

		// Provide logger with config
		fx.Provide(func(cfg *config.Config) (*zap.Logger, error) {
			if cfg.IsDevelopment() {
				return logger.NewDevelopment()
			}
			return logger.NewProduction()
		}),

		fx.Invoke(func(log *zap.Logger, cfg *config.Config) {
			log.Info("Using custom logger",
				zap.String("env", cfg.Environment),
				zap.String("level", cfg.LogLevel),
			)
		}),
	)

	app.Run()
}

// Example 5: Testing with FX
func ExampleTesting() {
	// In tests, you can easily mock dependencies
	mockConfig := &config.Config{
		Environment: "test",
		ServerPort:  8888,
	}

	mockLogger, _ := zap.NewDevelopment()

	app := fx.New(
		// Provide mocks
		fx.Supply(mockConfig),
		fx.Supply(mockLogger),

		// Test your service
		fx.Provide(NewMyService),

		fx.Invoke(func(svc *MyService) {
			// Assertions here
		}),
	)

	app.Run()
}

// Example 6: Module Pattern (Recommended for large apps)
var MyModule = fx.Module("myapp",
	// Include base modules
	config.Module,
	logger.Module,

	// Provide services
	fx.Provide(
		NewDatabaseService,
		NewAPIService,
	),
)

func ExampleModulePattern() {
	app := fx.New(
		// Just use your module
		MyModule,

		fx.Invoke(func(api *APIService) {
			api.logger.Info("App ready")
		}),
	)

	app.Run()
}

*/
