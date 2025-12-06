package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/Leon180/tabelogo-v2/internal/restaurant/application"
	"github.com/Leon180/tabelogo-v2/internal/restaurant/infrastructure/grpc"
	restaurantpostgres "github.com/Leon180/tabelogo-v2/internal/restaurant/infrastructure/postgres"
	"github.com/Leon180/tabelogo-v2/pkg/config"
	redisclient "github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
	grpclib "google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewMapServiceConnection creates a gRPC connection to Map Service
func NewMapServiceConnection(cfg *config.Config, logger *zap.Logger) (*grpclib.ClientConn, error) {
	grpcConfig := &grpc.ConnectionConfig{
		Address:          cfg.MapService.GRPCAddr,
		Timeout:          cfg.MapService.Timeout,
		MaxRetries:       3,
		KeepAliveTime:    30 * time.Second,
		KeepAliveTimeout: 10 * time.Second,
	}
	return grpc.NewMapServiceConnection(grpcConfig, logger)
}

// NewMapServiceClient creates a Map Service gRPC client
func NewMapServiceClient(conn *grpclib.ClientConn, cfg *config.Config, logger *zap.Logger) application.MapServiceClient {
	return grpc.NewMapServiceClient(conn, logger, cfg.MapService.Timeout)
}

// Module provides infrastructure dependencies
var Module = fx.Module("restaurant.infrastructure",
	fx.Provide(
		NewDatabase,
		NewRedis,
		restaurantpostgres.NewRestaurantRepository,
		restaurantpostgres.NewFavoriteRepository,
		// Map Service integration
		NewMapServiceConnection,
		NewMapServiceClient,
	),
)

// NewDatabase creates a new database connection with lifecycle management
// Restaurant Service connects to restaurant_db on port 5433
func NewDatabase(cfg *config.Config, lc fx.Lifecycle, logger *zap.Logger) (*gorm.DB, error) {
	dsn := cfg.GetDatabaseDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Connecting to restaurant database", zap.String("host", cfg.Database.Host))
			if err := sqlDB.PingContext(ctx); err != nil {
				return fmt.Errorf("failed to ping database: %w", err)
			}
			logger.Info("Restaurant database connected successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing restaurant database connection")
			return sqlDB.Close()
		},
	})

	return db, nil
}

// NewRedis creates a new Redis client with lifecycle management
// Restaurant Service uses Redis DB 1
func NewRedis(cfg *config.Config, lc fx.Lifecycle, logger *zap.Logger) *redisclient.Client {
	rdb := redisclient.NewClient(&redisclient.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Connecting to Redis", zap.String("addr", cfg.GetRedisAddr()), zap.Int("db", cfg.Redis.DB))
			if err := rdb.Ping(ctx).Err(); err != nil {
				return fmt.Errorf("failed to connect to Redis: %w", err)
			}
			logger.Info("Redis connected successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing Redis connection")
			return rdb.Close()
		},
	})

	return rdb
}
