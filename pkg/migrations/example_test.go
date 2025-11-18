package migrations_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/Leon180/tabelogo-v2/pkg/migrations"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Example_basic 基本使用範例
func Example_basic() {
	// 建立資料庫連線
	db, err := sql.Open("postgres", "postgres://user:pass@localhost/mydb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 建立 logger
	logger, _ := zap.NewProduction()

	// 建立 migration manager
	mgr, err := migrations.NewManager(migrations.Config{
		DB:             db,
		Logger:         logger,
		MigrationsPath: "file://migrations/auth",
		ServiceName:    "auth",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer mgr.Close()

	// 執行所有 migrations
	ctx := context.Background()
	if err := mgr.Up(ctx); err != nil {
		log.Fatal(err)
	}

	// 取得當前版本
	version, dirty, err := mgr.Version()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current version: %d, dirty: %v\n", version, dirty)
}

// Example_steps 步進式 migration 範例
func Example_steps() {
	db, _ := sql.Open("postgres", "postgres://user:pass@localhost/mydb?sslmode=disable")
	defer db.Close()

	logger, _ := zap.NewProduction()

	mgr, _ := migrations.NewManager(migrations.Config{
		DB:             db,
		Logger:         logger,
		MigrationsPath: "file://migrations/auth",
		ServiceName:    "auth",
	})
	defer mgr.Close()

	ctx := context.Background()

	// 執行 1 步 up
	if err := mgr.Steps(ctx, 1); err != nil {
		log.Fatal(err)
	}

	// 執行 2 步 down
	if err := mgr.Steps(ctx, -2); err != nil {
		log.Fatal(err)
	}
}

// Example_migrate_to_version 遷移到指定版本
func Example_migrateToVersion() {
	db, _ := sql.Open("postgres", "postgres://user:pass@localhost/mydb?sslmode=disable")
	defer db.Close()

	logger, _ := zap.NewProduction()

	mgr, _ := migrations.NewManager(migrations.Config{
		DB:             db,
		Logger:         logger,
		MigrationsPath: "file://migrations/auth",
		ServiceName:    "auth",
	})
	defer mgr.Close()

	ctx := context.Background()

	// 遷移到版本 3
	if err := mgr.Migrate(ctx, 3); err != nil {
		log.Fatal(err)
	}
}

// Example_validate 驗證 migration 狀態
func Example_validate() {
	db, _ := sql.Open("postgres", "postgres://user:pass@localhost/mydb?sslmode=disable")
	defer db.Close()

	logger, _ := zap.NewProduction()

	mgr, _ := migrations.NewManager(migrations.Config{
		DB:             db,
		Logger:         logger,
		MigrationsPath: "file://migrations/auth",
		ServiceName:    "auth",
	})
	defer mgr.Close()

	ctx := context.Background()

	// 驗證狀態
	if err := mgr.Validate(ctx); err != nil {
		log.Printf("Migration state is invalid: %v", err)

		// 如果是 dirty 狀態,可以嘗試修復
		version, dirty, _ := mgr.Version()
		if dirty {
			log.Printf("Migration is dirty at version %d", version)
			// 手動修復後,強制設定版本
			if err := mgr.Force(int(version)); err != nil {
				log.Fatal(err)
			}
		}
	}
}

// Example_fx 使用 Uber FX
func Example_fx() {
	app := fx.New(
		// 提供依賴
		fx.Provide(
			func() (*sql.DB, error) {
				return sql.Open("postgres", "postgres://user:pass@localhost/mydb?sslmode=disable")
			},
			func() (*zap.Logger, error) {
				return zap.NewProduction()
			},
			// 提供服務名稱
			fx.Annotate(
				func() string { return "auth" },
				fx.ResultTags(`name:"service_name"`),
			),
			// 提供 migrations 路徑
			fx.Annotate(
				func() string { return "file://migrations/auth" },
				fx.ResultTags(`name:"migrations_path"`),
			),
		),

		// 註冊 migration manager
		migrations.ProvideFx(),

		// 啟動時自動執行 migrations
		migrations.InvokeAutoMigrate(),

		// 也可以手動處理 migrations
		fx.Invoke(func(lc fx.Lifecycle, mgr *migrations.Manager, logger *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					info, err := mgr.GetInfo(ctx)
					if err != nil {
						return err
					}
					logger.Info("Migration info",
						zap.Uint("version", info.Version),
						zap.Bool("dirty", info.Dirty),
					)
					return nil
				},
			})
		}),
	)

	app.Run()
}

// Example_multiService 多服務 migration 管理
func Example_multiService() {
	db, _ := sql.Open("postgres", "postgres://user:pass@localhost/mydb?sslmode=disable")
	defer db.Close()

	logger, _ := zap.NewProduction()
	ctx := context.Background()

	// Auth Service
	authMgr, _ := migrations.NewManager(migrations.Config{
		DB:             db,
		Logger:         logger,
		MigrationsPath: "file://migrations/auth",
		ServiceName:    "auth",
	})
	defer authMgr.Close()

	if err := authMgr.Up(ctx); err != nil {
		log.Fatal(err)
	}

	// Restaurant Service
	restaurantMgr, _ := migrations.NewManager(migrations.Config{
		DB:             db,
		Logger:         logger,
		MigrationsPath: "file://migrations/restaurant",
		ServiceName:    "restaurant",
	})
	defer restaurantMgr.Close()

	if err := restaurantMgr.Up(ctx); err != nil {
		log.Fatal(err)
	}

	// 顯示所有服務的版本
	authVersion, _, _ := authMgr.Version()
	restaurantVersion, _, _ := restaurantMgr.Version()

	fmt.Printf("Auth version: %d\n", authVersion)
	fmt.Printf("Restaurant version: %d\n", restaurantVersion)
}
