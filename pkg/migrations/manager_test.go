package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// 測試用的資料庫配置
const (
	testDBHost     = "localhost"
	testDBPort     = "5432"
	testDBUser     = "postgres"
	testDBPassword = "postgres"
	testDBName     = "test_migrations"
)

// setupTestDB 建立測試資料庫
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	// 連接到 postgres 系統資料庫
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		testDBHost, testDBPort, testDBUser, testDBPassword)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Skipf("Cannot connect to postgres: %v. Skipping test.", err)
		return nil
	}

	// 檢查連線
	if err := db.Ping(); err != nil {
		t.Skipf("Cannot ping postgres: %v. Skipping test.", err)
		return nil
	}

	// 刪除測試資料庫 (如果存在)
	_, _ = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))

	// 建立測試資料庫
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", testDBName))
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	db.Close()

	// 連接到測試資料庫
	testDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		testDBHost, testDBPort, testDBUser, testDBPassword, testDBName)

	testDB, err := sql.Open("postgres", testDSN)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	return testDB
}

// cleanupTestDB 清理測試資料庫
func cleanupTestDB(t *testing.T, db *sql.DB) {
	t.Helper()

	if db != nil {
		db.Close()
	}

	// 連接到 postgres 系統資料庫
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		testDBHost, testDBPort, testDBUser, testDBPassword)

	sysDB, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Logf("Failed to connect to postgres for cleanup: %v", err)
		return
	}
	defer sysDB.Close()

	// 刪除測試資料庫
	_, _ = sysDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))
}

// createTestMigrations 建立測試用的 migration 檔案
func createTestMigrations(t *testing.T) string {
	t.Helper()

	// 建立臨時目錄
	tmpDir, err := os.MkdirTemp("", "test_migrations_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// 建立 migration 檔案
	migrations := map[string]string{
		"000001_create_test_table.up.sql": `
CREATE TABLE test_users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);`,
		"000001_create_test_table.down.sql": `DROP TABLE IF EXISTS test_users;`,
		"000002_add_username.up.sql": `
ALTER TABLE test_users ADD COLUMN username VARCHAR(50);`,
		"000002_add_username.down.sql": `
ALTER TABLE test_users DROP COLUMN IF EXISTS username;`,
	}

	for filename, content := range migrations {
		path := filepath.Join(tmpDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write migration file %s: %v", filename, err)
		}
	}

	return tmpDir
}

func TestNewManager(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	migrationsPath := createTestMigrations(t)
	defer os.RemoveAll(migrationsPath)

	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				DB:             db,
				Logger:         logger,
				MigrationsPath: "file://" + migrationsPath,
				ServiceName:    "test",
			},
			wantErr: false,
		},
		{
			name: "missing db",
			config: Config{
				Logger:         logger,
				MigrationsPath: "file://" + migrationsPath,
				ServiceName:    "test",
			},
			wantErr: true,
		},
		{
			name: "missing migrations path",
			config: Config{
				DB:          db,
				Logger:      logger,
				ServiceName: "test",
			},
			wantErr: true,
		},
		{
			name: "missing service name",
			config: Config{
				DB:             db,
				Logger:         logger,
				MigrationsPath: "file://" + migrationsPath,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr, err := NewManager(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if mgr != nil {
				defer mgr.Close()
			}
		})
	}
}

func TestManager_Up(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	migrationsPath := createTestMigrations(t)
	defer os.RemoveAll(migrationsPath)

	logger, _ := zap.NewDevelopment()

	mgr, err := NewManager(Config{
		DB:             db,
		Logger:         logger,
		MigrationsPath: "file://" + migrationsPath,
		ServiceName:    "test",
	})
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer mgr.Close()

	ctx := context.Background()

	// 執行 up
	err = mgr.Up(ctx)
	if err != nil {
		t.Fatalf("Up() failed: %v", err)
	}

	// 檢查版本
	version, dirty, err := mgr.Version()
	if err != nil {
		t.Fatalf("Version() failed: %v", err)
	}

	if version != 2 {
		t.Errorf("Expected version 2, got %d", version)
	}

	if dirty {
		t.Errorf("Migration is dirty")
	}

	// 檢查表格是否存在
	var exists bool
	err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'test_users')").Scan(&exists)
	if err != nil {
		t.Fatalf("Failed to check table existence: %v", err)
	}

	if !exists {
		t.Errorf("Table test_users does not exist")
	}

	// 檢查 username 欄位是否存在
	err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'test_users' AND column_name = 'username')").Scan(&exists)
	if err != nil {
		t.Fatalf("Failed to check column existence: %v", err)
	}

	if !exists {
		t.Errorf("Column username does not exist")
	}
}

func TestManager_Down(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	migrationsPath := createTestMigrations(t)
	defer os.RemoveAll(migrationsPath)

	logger, _ := zap.NewDevelopment()

	mgr, err := NewManager(Config{
		DB:             db,
		Logger:         logger,
		MigrationsPath: "file://" + migrationsPath,
		ServiceName:    "test",
	})
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer mgr.Close()

	ctx := context.Background()

	// 先執行 up
	if err := mgr.Up(ctx); err != nil {
		t.Fatalf("Up() failed: %v", err)
	}

	// 執行 down
	if err := mgr.Down(ctx); err != nil {
		t.Fatalf("Down() failed: %v", err)
	}

	// 檢查版本
	version, _, err := mgr.Version()
	if err != nil {
		t.Fatalf("Version() failed: %v", err)
	}

	if version != 1 {
		t.Errorf("Expected version 1, got %d", version)
	}

	// 檢查 username 欄位是否不存在
	var exists bool
	err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'test_users' AND column_name = 'username')").Scan(&exists)
	if err != nil {
		t.Fatalf("Failed to check column existence: %v", err)
	}

	if exists {
		t.Errorf("Column username should not exist after down migration")
	}
}

func TestManager_Steps(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	migrationsPath := createTestMigrations(t)
	defer os.RemoveAll(migrationsPath)

	logger, _ := zap.NewDevelopment()

	mgr, err := NewManager(Config{
		DB:             db,
		Logger:         logger,
		MigrationsPath: "file://" + migrationsPath,
		ServiceName:    "test",
	})
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer mgr.Close()

	ctx := context.Background()

	// 執行 1 步 up
	if err := mgr.Steps(ctx, 1); err != nil {
		t.Fatalf("Steps(1) failed: %v", err)
	}

	version, _, err := mgr.Version()
	if err != nil {
		t.Fatalf("Version() failed: %v", err)
	}

	if version != 1 {
		t.Errorf("Expected version 1, got %d", version)
	}

	// 再執行 1 步 up
	if err := mgr.Steps(ctx, 1); err != nil {
		t.Fatalf("Steps(1) failed: %v", err)
	}

	version, _, err = mgr.Version()
	if err != nil {
		t.Fatalf("Version() failed: %v", err)
	}

	if version != 2 {
		t.Errorf("Expected version 2, got %d", version)
	}

	// 執行 1 步 down
	if err := mgr.Steps(ctx, -1); err != nil {
		t.Fatalf("Steps(-1) failed: %v", err)
	}

	version, _, err = mgr.Version()
	if err != nil {
		t.Fatalf("Version() failed: %v", err)
	}

	if version != 1 {
		t.Errorf("Expected version 1 after down, got %d", version)
	}
}

func TestManager_GetInfo(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	migrationsPath := createTestMigrations(t)
	defer os.RemoveAll(migrationsPath)

	logger, _ := zap.NewDevelopment()

	mgr, err := NewManager(Config{
		DB:             db,
		Logger:         logger,
		MigrationsPath: "file://" + migrationsPath,
		ServiceName:    "test",
	})
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer mgr.Close()

	ctx := context.Background()

	// 執行 migrations
	if err := mgr.Up(ctx); err != nil {
		t.Fatalf("Up() failed: %v", err)
	}

	// 取得資訊
	info, err := mgr.GetInfo(ctx)
	if err != nil {
		t.Fatalf("GetInfo() failed: %v", err)
	}

	if info.Version != 2 {
		t.Errorf("Expected version 2, got %d", info.Version)
	}

	if info.Dirty {
		t.Errorf("Expected not dirty")
	}

	// 檢查 timestamp 是否合理
	if time.Since(info.Timestamp) > time.Second {
		t.Errorf("Timestamp seems incorrect: %v", info.Timestamp)
	}
}

func TestManager_Validate(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	migrationsPath := createTestMigrations(t)
	defer os.RemoveAll(migrationsPath)

	logger, _ := zap.NewDevelopment()

	mgr, err := NewManager(Config{
		DB:             db,
		Logger:         logger,
		MigrationsPath: "file://" + migrationsPath,
		ServiceName:    "test",
	})
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer mgr.Close()

	ctx := context.Background()

	// 執行 migrations
	if err := mgr.Up(ctx); err != nil {
		t.Fatalf("Up() failed: %v", err)
	}

	// 驗證
	if err := mgr.Validate(ctx); err != nil {
		t.Errorf("Validate() failed: %v", err)
	}
}
