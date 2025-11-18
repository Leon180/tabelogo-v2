package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

// Manager 管理資料庫 migration
type Manager struct {
	db             *sql.DB
	logger         *zap.Logger
	migrationsPath string
	serviceName    string
	client         *migrate.Migrate
}

// Config migration 管理器配置
type Config struct {
	DB             *sql.DB
	Logger         *zap.Logger
	MigrationsPath string // migrations 檔案路徑，例如 "file://migrations/auth"
	ServiceName    string // 服務名稱，例如 "auth"
}

// NewManager 建立 migration 管理器
func NewManager(cfg Config) (*Manager, error) {
	if cfg.DB == nil {
		return nil, fmt.Errorf("database connection is required")
	}
	if cfg.MigrationsPath == "" {
		return nil, fmt.Errorf("migrations path is required")
	}
	if cfg.ServiceName == "" {
		return nil, fmt.Errorf("service name is required")
	}

	logger := cfg.Logger
	if logger == nil {
		logger, _ = zap.NewProduction()
	}

	// 建立 postgres driver instance
	driver, err := postgres.WithInstance(cfg.DB, &postgres.Config{
		MigrationsTable: fmt.Sprintf("schema_migrations_%s", cfg.ServiceName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// 建立 migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		cfg.MigrationsPath,
		cfg.ServiceName,
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return &Manager{
		db:             cfg.DB,
		logger:         logger,
		migrationsPath: cfg.MigrationsPath,
		serviceName:    cfg.ServiceName,
		client:         m,
	}, nil
}

// Up 執行所有未執行的 migrations
func (m *Manager) Up(ctx context.Context) error {
	m.logger.Info("running migrations up",
		zap.String("service", m.serviceName),
		zap.String("path", m.migrationsPath),
	)

	if err := m.client.Up(); err != nil && err != migrate.ErrNoChange {
		m.logger.Error("failed to run migrations up", zap.Error(err))
		return fmt.Errorf("failed to run migrations up: %w", err)
	}

	version, dirty, err := m.client.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	m.logger.Info("migrations completed",
		zap.Uint("version", version),
		zap.Bool("dirty", dirty),
	)

	return nil
}

// Down 回滾一個 migration
func (m *Manager) Down(ctx context.Context) error {
	m.logger.Info("running migration down",
		zap.String("service", m.serviceName),
	)

	if err := m.client.Down(); err != nil && err != migrate.ErrNoChange {
		m.logger.Error("failed to run migration down", zap.Error(err))
		return fmt.Errorf("failed to run migration down: %w", err)
	}

	version, dirty, err := m.client.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	m.logger.Info("migration down completed",
		zap.Uint("version", version),
		zap.Bool("dirty", dirty),
	)

	return nil
}

// Steps 執行指定步數的 migration (正數為 up，負數為 down)
func (m *Manager) Steps(ctx context.Context, n int) error {
	m.logger.Info("running migration steps",
		zap.String("service", m.serviceName),
		zap.Int("steps", n),
	)

	if err := m.client.Steps(n); err != nil && err != migrate.ErrNoChange {
		m.logger.Error("failed to run migration steps", zap.Error(err))
		return fmt.Errorf("failed to run migration steps: %w", err)
	}

	version, dirty, err := m.client.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	m.logger.Info("migration steps completed",
		zap.Uint("version", version),
		zap.Bool("dirty", dirty),
	)

	return nil
}

// Migrate 遷移到指定版本
func (m *Manager) Migrate(ctx context.Context, version uint) error {
	m.logger.Info("migrating to version",
		zap.String("service", m.serviceName),
		zap.Uint("target_version", version),
	)

	if err := m.client.Migrate(version); err != nil && err != migrate.ErrNoChange {
		m.logger.Error("failed to migrate to version", zap.Error(err))
		return fmt.Errorf("failed to migrate to version %d: %w", version, err)
	}

	currentVersion, dirty, err := m.client.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	m.logger.Info("migration completed",
		zap.Uint("current_version", currentVersion),
		zap.Bool("dirty", dirty),
	)

	return nil
}

// Version 取得當前 migration 版本
func (m *Manager) Version() (version uint, dirty bool, err error) {
	version, dirty, err = m.client.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}
	return version, dirty, nil
}

// Force 強制設定 migration 版本 (用於修復 dirty 狀態)
func (m *Manager) Force(version int) error {
	m.logger.Warn("forcing migration version",
		zap.String("service", m.serviceName),
		zap.Int("version", version),
	)

	if err := m.client.Force(version); err != nil {
		m.logger.Error("failed to force migration version", zap.Error(err))
		return fmt.Errorf("failed to force migration version: %w", err)
	}

	m.logger.Info("migration version forced", zap.Int("version", version))
	return nil
}

// Drop 刪除所有表格 (危險操作，僅用於開發環境)
func (m *Manager) Drop(ctx context.Context) error {
	m.logger.Warn("dropping all tables",
		zap.String("service", m.serviceName),
	)

	if err := m.client.Drop(); err != nil {
		m.logger.Error("failed to drop tables", zap.Error(err))
		return fmt.Errorf("failed to drop tables: %w", err)
	}

	m.logger.Info("all tables dropped")
	return nil
}

// Close 關閉 migration manager
func (m *Manager) Close() error {
	srcErr, dbErr := m.client.Close()
	if srcErr != nil {
		return srcErr
	}
	if dbErr != nil {
		return dbErr
	}
	return nil
}

// MigrationInfo 包含 migration 資訊
type MigrationInfo struct {
	Version   uint
	Dirty     bool
	Timestamp time.Time
}

// GetInfo 取得當前 migration 資訊
func (m *Manager) GetInfo(ctx context.Context) (*MigrationInfo, error) {
	version, dirty, err := m.Version()
	if err != nil {
		return nil, err
	}

	return &MigrationInfo{
		Version:   version,
		Dirty:     dirty,
		Timestamp: time.Now(),
	}, nil
}

// Validate 驗證 migration 狀態
func (m *Manager) Validate(ctx context.Context) error {
	version, dirty, err := m.Version()
	if err != nil {
		return fmt.Errorf("failed to get version: %w", err)
	}

	if dirty {
		return fmt.Errorf("migration is in dirty state at version %d, please fix manually", version)
	}

	m.logger.Info("migration state is valid",
		zap.Uint("version", version),
	)

	return nil
}
