package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// FxParams FX 依賴注入參數
type FxParams struct {
	fx.In

	DB             *sql.DB
	Logger         *zap.Logger
	ServiceName    string `name:"service_name"`
	MigrationsPath string `name:"migrations_path"`
}

// ProvideFx 提供 migration manager 給 FX
func ProvideFx() fx.Option {
	return fx.Provide(NewManagerFromFx)
}

// NewManagerFromFx 從 FX 建立 migration manager
func NewManagerFromFx(params FxParams) (*Manager, error) {
	return NewManager(Config{
		DB:             params.DB,
		Logger:         params.Logger,
		MigrationsPath: params.MigrationsPath,
		ServiceName:    params.ServiceName,
	})
}

// FxInvokeParams FX invoke 參數
type FxInvokeParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Manager   *Manager
	Logger    *zap.Logger
}

// InvokeAutoMigrate 提供自動執行 migration 的 FX invoke
func InvokeAutoMigrate() fx.Option {
	return fx.Invoke(func(params FxInvokeParams) error {
		params.Lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				params.Logger.Info("running auto migrations on startup")

				if err := params.Manager.Validate(ctx); err != nil {
					params.Logger.Warn("migration validation failed, will attempt to migrate", zap.Error(err))
				}

				if err := params.Manager.Up(ctx); err != nil {
					return fmt.Errorf("auto migration failed: %w", err)
				}

				info, err := params.Manager.GetInfo(ctx)
				if err != nil {
					return fmt.Errorf("failed to get migration info: %w", err)
				}

				params.Logger.Info("auto migrations completed",
					zap.Uint("version", info.Version),
					zap.Bool("dirty", info.Dirty),
				)

				return nil
			},
			OnStop: func(ctx context.Context) error {
				params.Logger.Info("closing migration manager")
				return params.Manager.Close()
			},
		})
		return nil
	})
}
