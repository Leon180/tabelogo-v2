package infrastructure

import (
	"context"
	"fmt"

	authpostgres "github.com/Leon180/tabelogo-v2/internal/auth/infrastructure/postgres"
	authredis "github.com/Leon180/tabelogo-v2/internal/auth/infrastructure/redis"
	"github.com/Leon180/tabelogo-v2/pkg/config"
	redisclient "github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Module provides infrastructure dependencies
var Module = fx.Module("auth.infrastructure",
	fx.Provide(
		NewDatabase,
		NewRedis,
		authpostgres.NewUserRepository,
		authredis.NewTokenRepository,
	),
)

// NewDatabase creates a new database connection with lifecycle management
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
			logger.Info("Connecting to database", zap.String("host", cfg.Database.Host))
			if err := sqlDB.PingContext(ctx); err != nil {
				return fmt.Errorf("failed to ping database: %w", err)
			}
			logger.Info("Database connected successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing database connection")
			return sqlDB.Close()
		},
	})

	return db, nil
}

// NewRedis creates a new Redis client with lifecycle management
func NewRedis(cfg *config.Config, lc fx.Lifecycle, logger *zap.Logger) *redisclient.Client {
	rdb := redisclient.NewClient(&redisclient.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Connecting to Redis", zap.String("addr", cfg.GetRedisAddr()))
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
