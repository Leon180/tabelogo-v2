package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/Leon180/tabelogo-v2/pkg/migrations"
	"go.uber.org/zap"
)

func main() {
	var (
		dsn            string
		migrationsPath string
		serviceName    string
		command        string
		steps          int
		version        int
	)

	flag.StringVar(&dsn, "dsn", "", "Database connection string (required)")
	flag.StringVar(&migrationsPath, "path", "", "Path to migrations directory (required)")
	flag.StringVar(&serviceName, "service", "", "Service name (required)")
	flag.StringVar(&command, "command", "up", "Command: up, down, steps, migrate, version, validate, force, drop")
	flag.IntVar(&steps, "steps", 0, "Number of steps for 'steps' command")
	flag.IntVar(&version, "version", 0, "Target version for 'migrate' or 'force' command")
	flag.Parse()

	// 驗證必要參數
	if dsn == "" || migrationsPath == "" || serviceName == "" {
		fmt.Println("Error: -dsn, -path, and -service are required")
		flag.Usage()
		os.Exit(1)
	}

	// 建立 logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// 連接資料庫
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Fatal("Failed to open database", zap.Error(err))
	}
	defer db.Close()

	// 測試連線
	if err := db.Ping(); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}

	// 建立 migration manager
	mgr, err := migrations.NewManager(migrations.Config{
		DB:             db,
		Logger:         logger,
		MigrationsPath: "file://" + migrationsPath,
		ServiceName:    serviceName,
	})
	if err != nil {
		logger.Fatal("Failed to create migration manager", zap.Error(err))
	}
	defer mgr.Close()

	ctx := context.Background()

	// 執行命令
	switch command {
	case "up":
		if err := mgr.Up(ctx); err != nil {
			logger.Fatal("Failed to run migrations up", zap.Error(err))
		}
		logger.Info("Migrations up completed")
		printVersion(mgr)

	case "down":
		if err := mgr.Down(ctx); err != nil {
			logger.Fatal("Failed to run migration down", zap.Error(err))
		}
		logger.Info("Migration down completed")
		printVersion(mgr)

	case "steps":
		if steps == 0 {
			logger.Fatal("steps parameter is required for 'steps' command")
		}
		if err := mgr.Steps(ctx, steps); err != nil {
			logger.Fatal("Failed to run migration steps", zap.Error(err))
		}
		logger.Info("Migration steps completed", zap.Int("steps", steps))
		printVersion(mgr)

	case "migrate":
		if version == 0 {
			logger.Fatal("version parameter is required for 'migrate' command")
		}
		if err := mgr.Migrate(ctx, uint(version)); err != nil {
			logger.Fatal("Failed to migrate to version", zap.Error(err))
		}
		logger.Info("Migration completed", zap.Int("target_version", version))
		printVersion(mgr)

	case "version":
		printVersion(mgr)

	case "validate":
		if err := mgr.Validate(ctx); err != nil {
			logger.Error("Migration validation failed", zap.Error(err))
			os.Exit(1)
		}
		logger.Info("Migration state is valid")
		printVersion(mgr)

	case "force":
		if version == 0 {
			logger.Fatal("version parameter is required for 'force' command")
		}
		if err := mgr.Force(version); err != nil {
			logger.Fatal("Failed to force version", zap.Error(err))
		}
		logger.Info("Version forced", zap.Int("version", version))
		printVersion(mgr)

	case "drop":
		fmt.Print("Are you sure you want to drop all tables? (yes/no): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			logger.Info("Drop cancelled")
			return
		}
		if err := mgr.Drop(ctx); err != nil {
			logger.Fatal("Failed to drop tables", zap.Error(err))
		}
		logger.Info("All tables dropped")

	default:
		logger.Fatal("Unknown command", zap.String("command", command))
	}
}

func printVersion(mgr *migrations.Manager) {
	version, dirty, err := mgr.Version()
	if err != nil {
		fmt.Printf("Failed to get version: %v\n", err)
		return
	}

	dirtyStr := ""
	if dirty {
		dirtyStr = " (DIRTY)"
	}

	fmt.Printf("Current version: %d%s\n", version, dirtyStr)
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Migration management tool for database schema versioning.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nCommands:\n")
		fmt.Fprintf(os.Stderr, "  up       - Run all pending migrations\n")
		fmt.Fprintf(os.Stderr, "  down     - Rollback the last migration\n")
		fmt.Fprintf(os.Stderr, "  steps    - Run N migrations (positive=up, negative=down)\n")
		fmt.Fprintf(os.Stderr, "  migrate  - Migrate to a specific version\n")
		fmt.Fprintf(os.Stderr, "  version  - Show current migration version\n")
		fmt.Fprintf(os.Stderr, "  validate - Validate migration state\n")
		fmt.Fprintf(os.Stderr, "  force    - Force set migration version (use with caution)\n")
		fmt.Fprintf(os.Stderr, "  drop     - Drop all tables (dangerous!)\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  # Run all pending migrations\n")
		fmt.Fprintf(os.Stderr, "  %s -dsn \"postgres://user:pass@localhost/db\" -path migrations/auth -service auth -command up\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  # Rollback last migration\n")
		fmt.Fprintf(os.Stderr, "  %s -dsn \"postgres://user:pass@localhost/db\" -path migrations/auth -service auth -command down\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  # Run 2 migrations forward\n")
		fmt.Fprintf(os.Stderr, "  %s -dsn \"postgres://user:pass@localhost/db\" -path migrations/auth -service auth -command steps -steps 2\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  # Migrate to version 5\n")
		fmt.Fprintf(os.Stderr, "  %s -dsn \"postgres://user:pass@localhost/db\" -path migrations/auth -service auth -command migrate -version 5\n\n", os.Args[0])
	}
}

// getEnv 從環境變數取得值,如果不存在則返回預設值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt 從環境變數取得整數值
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
