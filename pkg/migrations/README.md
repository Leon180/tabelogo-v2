# Migrations Package

這個套件提供了資料庫 migration 的版本控制與管理功能，基於 `golang-migrate/migrate`。

## 目錄

- [功能特性](#功能特性)
- [核心概念](#核心概念)
- [版本控制機制詳解](#版本控制機制詳解)
- [使用方式](#使用方式)
- [Migration 檔案管理](#migration-檔案管理)
- [API 參考](#api-參考)
- [進階主題](#進階主題)
- [最佳實踐](#最佳實踐)
- [常見問題](#常見問題)
- [疑難排解](#疑難排解)

## 功能特性

- ✅ 支援多服務獨立 migration 管理
- ✅ 版本控制與狀態追蹤
- ✅ Up/Down migration 支援
- ✅ 步進式 migration (Steps)
- ✅ 版本強制設定 (Force) - 用於修復 dirty 狀態
- ✅ Migration 狀態驗證
- ✅ 整合 Uber FX 依賴注入
- ✅ 結構化日誌 (zap)
- ✅ CLI 工具支援

## 核心概念

### 什麼是 Migration?

Migration 是一種版本化的資料庫變更腳本，讓你可以：
1. **追蹤資料庫結構變更歷史**
2. **在不同環境間同步資料庫結構** (開發、測試、生產)
3. **回滾錯誤的變更**
4. **協同開發時避免衝突**

### Migration 的生命週期

```
1. 建立 Migration 檔案
   ├── 000001_create_users.up.sql   (向上遷移)
   └── 000001_create_users.down.sql (向下回滾)

2. 執行 Migration (up)
   ├── 讀取未執行的 migration 檔案
   ├── 按照版本順序執行
   └── 記錄到版本控制表

3. 版本狀態追蹤
   └── 儲存在 schema_migrations_{service} 表

4. 需要時回滾 (down)
   ├── 執行 down.sql
   └── 更新版本記錄
```

## 使用方式

### 1. 基本使用

```go
package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/Leon180/tabelogo-v2/pkg/migrations"
	"go.uber.org/zap"
)

func main() {
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

	// 執行 migrations
	ctx := context.Background()
	if err := mgr.Up(ctx); err != nil {
		log.Fatal(err)
	}

	// 取得當前版本
	version, dirty, err := mgr.Version()
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("Current version", zap.Uint("version", version), zap.Bool("dirty", dirty))
}
```

### 2. 使用 Uber FX

```go
package main

import (
	"database/sql"

	"github.com/Leon180/tabelogo-v2/pkg/migrations"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		// 提供依賴
		fx.Provide(
			func() (*sql.DB, error) {
				return sql.Open("postgres", "postgres://user:pass@localhost/mydb?sslmode=disable")
			},
			func() (*zap.Logger, error) {
				return zap.NewProduction()
			},
			fx.Annotate(
				func() string { return "auth" },
				fx.ResultTags(`name:"service_name"`),
			),
			fx.Annotate(
				func() string { return "file://migrations/auth" },
				fx.ResultTags(`name:"migrations_path"`),
			),
		),

		// 註冊 migration manager
		migrations.ProvideFx(),

		// 啟動時自動執行 migrations
		migrations.InvokeAutoMigrate(),
	)

	app.Run()
}
```

## 版本控制機制詳解

### 版本號的來源與意義

Migration 的版本號來自於**檔案名稱的前綴數字**，系統會自動解析並排序。

#### 版本號格式

有兩種常見的版本號格式：

**1. 序列號格式** (Sequence)
```
000001_create_users.up.sql
000002_add_email_verified.up.sql
000003_create_roles.up.sql
```
- 優點：簡潔易讀
- 缺點：多人協作時容易衝突

**2. 時間戳格式** (Timestamp) ⭐ **推薦**
```
20250118120000_create_users.up.sql
20250118130000_add_email_verified.up.sql
20250119100000_create_roles.up.sql
```
- 格式：`YYYYMMDDHHmmss`
- 優點：避免多人開發時的版本號衝突
- 缺點：版本號較長

#### 如何建立帶版本號的 Migration?

**方法 1: 使用 golang-migrate CLI (推薦)**

```bash
# 安裝 migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 建立序列號格式的 migration
migrate create -ext sql -dir migrations/auth -seq create_users_table

# 建立時間戳格式的 migration
migrate create -ext sql -dir migrations/auth create_users_table
```

**方法 2: 手動建立**

```bash
# 序列號格式
touch migrations/auth/000001_create_users_table.up.sql
touch migrations/auth/000001_create_users_table.down.sql

# 時間戳格式 (使用當前時間)
touch migrations/auth/$(date +%Y%m%d%H%M%S)_create_users_table.up.sql
touch migrations/auth/$(date +%Y%m%d%H%M%S)_create_users_table.down.sql
```

### 版本控制表結構

每個服務在資料庫中會建立獨立的版本控制表，表名格式為：`schema_migrations_{service}`

#### 表結構

```sql
-- 例如 auth service 的版本控制表
CREATE TABLE schema_migrations_auth (
    version BIGINT NOT NULL PRIMARY KEY,
    dirty BOOLEAN NOT NULL
);
```

**欄位說明**：
- `version`: 當前執行到的 migration 版本號
- `dirty`: 標記 migration 是否執行失敗 (true=失敗/未完成)

#### 查看版本控制表

**方法 1: 使用 SQL 查詢**

```sql
-- 查看 auth service 的 migration 狀態
SELECT * FROM schema_migrations_auth;

-- 範例輸出：
-- version  | dirty
-- ---------|------
-- 2        | false

-- 查看所有服務的 migration 表
SELECT table_name
FROM information_schema.tables
WHERE table_name LIKE 'schema_migrations_%';
```

**方法 2: 使用 Migration Manager**

```go
version, dirty, err := mgr.Version()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Current version: %d, dirty: %v\n", version, dirty)
```

**方法 3: 使用 CLI 工具**

```bash
./migrate -dsn "postgres://..." -path migrations/auth -service auth -command version
```

### 如何檢視有哪些可用的 Migration?

#### 方法 1: 檢視檔案系統

```bash
# 列出所有 migration 檔案
ls -1 migrations/auth/

# 輸出：
# 000001_create_users_table.up.sql
# 000001_create_users_table.down.sql
# 000002_add_email_verified.up.sql
# 000002_add_email_verified.down.sql
# 000003_create_roles.up.sql
# 000003_create_roles.down.sql

# 只列出版本號
ls migrations/auth/*.up.sql | sed 's/.*\/\([0-9]*\)_.*/\1/'

# 輸出：
# 000001
# 000002
# 000003
```

#### 方法 2: 使用程式碼列出

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "regexp"
    "sort"
    "strconv"
)

func listMigrations(dir string) []uint {
    var versions []uint
    pattern := regexp.MustCompile(`^(\d+)_.+\.up\.sql$`)

    files, _ := os.ReadDir(dir)
    for _, file := range files {
        if matches := pattern.FindStringSubmatch(file.Name()); len(matches) > 1 {
            version, _ := strconv.ParseUint(matches[1], 10, 64)
            versions = append(versions, uint(version))
        }
    }

    sort.Slice(versions, func(i, j int) bool {
        return versions[i] < versions[j]
    })

    return versions
}

func main() {
    versions := listMigrations("migrations/auth")
    fmt.Printf("Available migrations: %v\n", versions)
    // 輸出: Available migrations: [1 2 3]
}
```

### 版本狀態說明

#### Clean State (正常狀態)

```
┌─────────────────────────────┐
│ schema_migrations_auth      │
├─────────────────────────────┤
│ version | dirty             │
│ 2       | false             │  ✅ 正常狀態
└─────────────────────────────┘

可用的 migrations:
  ✅ 000001 (已執行)
  ✅ 000002 (已執行)
  ⏸️  000003 (未執行)
  ⏸️  000004 (未執行)
```

#### Dirty State (異常狀態)

```
┌─────────────────────────────┐
│ schema_migrations_auth      │
├─────────────────────────────┤
│ version | dirty             │
│ 2       | true              │  ❌ 異常狀態！
└─────────────────────────────┘

說明：
- version 2 的 migration 執行失敗或未完成
- 需要手動修復後才能繼續執行其他 migrations
```

### 版本比對與執行

當執行 `Up()` 時，系統會：

1. **讀取檔案系統中的所有 migration 檔案**
   ```
   Available files: [1, 2, 3, 4, 5]
   ```

2. **讀取資料庫中的當前版本**
   ```sql
   SELECT version FROM schema_migrations_auth;
   -- 結果: 2
   ```

3. **計算需要執行的 migrations**
   ```
   Current version: 2
   Available: [1, 2, 3, 4, 5]
   To execute: [3, 4, 5]  (所有大於當前版本的)
   ```

4. **按順序執行**
   ```
   Executing 3... ✅
   Executing 4... ✅
   Executing 5... ✅

   Final version: 5
   ```

### Migration 檔案範例

Migration 檔案應該放在 `migrations/{service}/` 目錄下:

```
migrations/
├── auth/
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_add_email_verified_column.up.sql
│   └── 000002_add_email_verified_column.down.sql
├── restaurant/
│   ├── 000001_create_restaurants_table.up.sql
│   └── 000001_create_restaurants_table.down.sql
```

**檔案命名規則**: `{version}_{description}.{up|down}.sql`

範例 `000001_create_users_table.up.sql`:
```sql
-- Migration: 000001
-- Description: Create users table for authentication
-- Author: Team Backend
-- Date: 2025-01-18

BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    username VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- Comments
COMMENT ON TABLE users IS 'User accounts for authentication service';
COMMENT ON COLUMN users.id IS 'Primary key - UUID v4';
COMMENT ON COLUMN users.email IS 'User email address - unique identifier';

COMMIT;
```

範例 `000001_create_users_table.down.sql`:
```sql
-- Rollback: 000001
-- Description: Drop users table

BEGIN;

DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;

COMMIT;
```

## API 參考

### Manager 方法

#### Up(ctx context.Context) error
執行所有未執行的 migrations

#### Down(ctx context.Context) error
回滾最後一個 migration

#### Steps(ctx context.Context, n int) error
執行指定步數的 migration
- `n > 0`: 向上執行 n 步
- `n < 0`: 向下回滾 n 步

#### Migrate(ctx context.Context, version uint) error
遷移到指定版本

#### Version() (version uint, dirty bool, err error)
取得當前 migration 版本

#### Force(version int) error
強制設定版本號 (用於修復 dirty 狀態)

#### Drop(ctx context.Context) error
刪除所有表格 ⚠️ 危險操作，僅用於開發環境

#### Validate(ctx context.Context) error
驗證 migration 狀態是否正常

#### GetInfo(ctx context.Context) (*MigrationInfo, error)
取得當前 migration 詳細資訊

## 進階主題

### 多服務 Migration 管理

在微服務架構中，每個服務有獨立的資料庫和 migration：

```go
package main

import (
    "context"
    "database/sql"
    "log"
    "sync"

    "github.com/Leon180/tabelogo-v2/pkg/migrations"
    "go.uber.org/zap"
)

func main() {
    logger, _ := zap.NewProduction()
    ctx := context.Background()

    // 定義所有服務
    services := []struct {
        name string
        dsn  string
        path string
    }{
        {"auth", "postgres://user:pass@localhost/auth_db", "migrations/auth"},
        {"restaurant", "postgres://user:pass@localhost/restaurant_db", "migrations/restaurant"},
        {"booking", "postgres://user:pass@localhost/booking_db", "migrations/booking"},
    }

    var wg sync.WaitGroup
    errors := make(chan error, len(services))

    // 並行執行所有服務的 migrations
    for _, svc := range services {
        wg.Add(1)
        go func(service struct{ name, dsn, path string }) {
            defer wg.Done()

            db, err := sql.Open("postgres", service.dsn)
            if err != nil {
                errors <- err
                return
            }
            defer db.Close()

            mgr, err := migrations.NewManager(migrations.Config{
                DB:             db,
                Logger:         logger.Named(service.name),
                MigrationsPath: "file://" + service.path,
                ServiceName:    service.name,
            })
            if err != nil {
                errors <- err
                return
            }
            defer mgr.Close()

            if err := mgr.Up(ctx); err != nil {
                errors <- err
                return
            }

            version, _, _ := mgr.Version()
            logger.Info("Migration completed",
                zap.String("service", service.name),
                zap.Uint("version", version),
            )
        }(svc)
    }

    wg.Wait()
    close(errors)

    // 檢查錯誤
    for err := range errors {
        log.Printf("Migration error: %v", err)
    }
}
```

### 版本控制表說明

每個服務會在資料庫中建立獨立的版本控制表:

| 服務 | 版本控制表名稱 | 資料庫 |
|------|--------------|--------|
| Auth Service | `schema_migrations_auth` | `auth_db` |
| Restaurant Service | `schema_migrations_restaurant` | `restaurant_db` |
| Booking Service | `schema_migrations_booking` | `booking_db` |
| Spider Service | `schema_migrations_spider` | `spider_db` |
| Mail Service | `schema_migrations_mail` | `mail_db` |

**查詢所有服務的版本狀態**：

```sql
-- 在 auth_db 中
SELECT 'auth' as service, version, dirty FROM schema_migrations_auth
UNION ALL
-- 在 restaurant_db 中
SELECT 'restaurant' as service, version, dirty FROM schema_migrations_restaurant
-- ... 其他服務
```

### 環境隔離

在不同環境中管理 migrations：

```go
package main

import (
    "fmt"
    "os"
)

func getDSN() string {
    env := os.Getenv("APP_ENV") // dev, test, staging, production

    switch env {
    case "production":
        return os.Getenv("PROD_DB_DSN")
    case "staging":
        return os.Getenv("STAGING_DB_DSN")
    case "test":
        return os.Getenv("TEST_DB_DSN")
    default: // dev
        return os.Getenv("DEV_DB_DSN")
    }
}

func main() {
    dsn := getDSN()
    fmt.Printf("Using database: %s\n", dsn)

    // ... 建立 migration manager
}
```

**環境變數設定**：

```bash
# .env.dev
APP_ENV=dev
DEV_DB_DSN=postgres://user:pass@localhost/mydb_dev

# .env.staging
APP_ENV=staging
STAGING_DB_DSN=postgres://user:pass@staging-db/mydb

# .env.production
APP_ENV=production
PROD_DB_DSN=postgres://user:pass@prod-db/mydb
```

## 疑難排解

### Dirty State 處理

當 migration 執行失敗時，會進入 "dirty" 狀態。

#### 什麼情況會導致 Dirty State?

1. **SQL 語法錯誤**
   ```sql
   -- 錯誤的 SQL
   CREATE TABLEE users (  -- 拼字錯誤
       id INT
   );
   ```

2. **違反資料庫約束**
   ```sql
   -- 嘗試建立已存在的表
   CREATE TABLE users (...);  -- 如果表已存在會失敗
   ```

3. **連線中斷**
   - 網路問題
   - 資料庫重啟
   - 超時

4. **權限不足**
   ```sql
   CREATE EXTENSION postgis;  -- 需要超級用戶權限
   ```

#### 修復 Dirty State 的步驟

**步驟 1: 檢查狀態**

```bash
# 使用 CLI
./migrate -dsn "postgres://..." -path migrations/auth -service auth -command version

# 輸出:
# Current version: 3 (DIRTY)
```

或使用程式碼：

```go
version, dirty, err := mgr.Version()
if dirty {
    fmt.Printf("⚠️  Migration is dirty at version %d\n", version)
}
```

**步驟 2: 查看失敗的 Migration**

```bash
# 找出版本 3 的 migration 檔案
cat migrations/auth/000003_*.up.sql
```

**步驟 3: 檢查資料庫實際狀態**

```sql
-- 檢查 migration 是否部分執行
-- 例如，如果 migration 要建立 3 個表，檢查哪些已建立

SELECT table_name
FROM information_schema.tables
WHERE table_schema = 'public'
  AND table_name IN ('table1', 'table2', 'table3');
```

**步驟 4: 手動修復**

有兩種策略：

**策略 A: 完成 Migration (推薦)**

如果 migration 部分執行，手動完成剩餘部分：

```sql
-- 手動執行未完成的部分
BEGIN;
CREATE TABLE table3 (...);  -- 假設這個沒執行
COMMIT;
```

然後強制設定版本：

```go
// 強制設定為版本 3 (已完成)
err := mgr.Force(3)
if err != nil {
    log.Fatal(err)
}
```

**策略 B: 回滾 Migration**

如果 migration 執行錯誤，回滾到上一個版本：

```sql
-- 手動執行 down migration
BEGIN;
-- 執行 000003_xxx.down.sql 的內容
COMMIT;
```

然後強制設定版本：

```go
// 強制設定為版本 2 (回滾到上一版)
err := mgr.Force(2)
if err != nil {
    log.Fatal(err)
}
```

**步驟 5: 驗證修復**

```go
// 驗證狀態
if err := mgr.Validate(ctx); err != nil {
    log.Fatal("Still dirty:", err)
}

version, dirty, _ := mgr.Version()
fmt.Printf("✅ Version: %d, Dirty: %v\n", version, dirty)
```

#### 完整修復範例

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"

    "github.com/Leon180/tabelogo-v2/pkg/migrations"
    "go.uber.org/zap"
)

func fixDirtyMigration(mgr *migrations.Manager, db *sql.DB) error {
    ctx := context.Background()

    // 1. 檢查狀態
    version, dirty, err := mgr.Version()
    if err != nil {
        return err
    }

    if !dirty {
        fmt.Println("✅ Migration is not dirty")
        return nil
    }

    fmt.Printf("⚠️  Found dirty migration at version %d\n", version)

    // 2. 詢問用戶如何處理
    fmt.Println("How to fix?")
    fmt.Println("1. Force complete (migration is done, just dirty)")
    fmt.Println("2. Force rollback (undo this migration)")
    fmt.Print("Choice (1/2): ")

    var choice int
    fmt.Scanln(&choice)

    switch choice {
    case 1:
        // 強制設定為當前版本 (標記為完成)
        fmt.Printf("Forcing version to %d (complete)...\n", version)
        if err := mgr.Force(int(version)); err != nil {
            return err
        }

    case 2:
        // 先手動清理，然後回滾到上一版
        fmt.Println("Please manually clean up the database first.")
        fmt.Print("Press Enter when done...")
        fmt.Scanln()

        fmt.Printf("Forcing version to %d (rollback)...\n", version-1)
        if err := mgr.Force(int(version - 1)); err != nil {
            return err
        }

    default:
        return fmt.Errorf("invalid choice")
    }

    // 3. 驗證
    if err := mgr.Validate(ctx); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    fmt.Println("✅ Migration fixed successfully!")
    return nil
}

func main() {
    db, _ := sql.Open("postgres", "postgres://...")
    defer db.Close()

    logger, _ := zap.NewProduction()

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

    if err := fixDirtyMigration(mgr, db); err != nil {
        log.Fatal(err)
    }
}
```

### 常見錯誤排除

#### 錯誤 1: "no change" 錯誤

```
Error: no change
```

**原因**: 所有 migrations 都已執行完畢，沒有新的可執行。

**解決方法**: 這不是錯誤，只是資訊提示。

```go
err := mgr.Up(ctx)
if err != nil && err != migrate.ErrNoChange {
    log.Fatal(err)  // 只處理真正的錯誤
}
```

#### 錯誤 2: "file does not exist"

```
Error: file://migrations/auth does not exist
```

**原因**: Migration 路徑不正確。

**解決方法**:

```go
// 確認路徑格式正確
// ✅ 正確
MigrationsPath: "file://migrations/auth"
MigrationsPath: "file:///absolute/path/to/migrations/auth"

// ❌ 錯誤
MigrationsPath: "migrations/auth"  // 缺少 file:// 前綴
```

#### 錯誤 3: "Dirty database version"

```
Error: Dirty database version 3. Fix and force version.
```

**解決方法**: 參考上面的 "Dirty State 處理" 章節。

#### 錯誤 4: 版本衝突

```
Error: migration version 5 already exists
```

**原因**: 多人開發時使用了相同的序列號。

**解決方法**: 改用時間戳格式的版本號：

```bash
# 使用時間戳格式
migrate create -ext sql -dir migrations/auth create_new_table
# 生成: 20250118153045_create_new_table.up.sql
```

## 最佳實踐

### 1. 版本號管理

**✅ 推薦: 使用時間戳格式**

```bash
migrate create -ext sql -dir migrations/auth create_users_table
# 生成: 20250118120000_create_users_table.up.sql
```

**優點**:
- 避免多人開發時的版本號衝突
- 可以從版本號看出建立時間
- 自然排序

### 2. Migration 設計原則

**✅ DO (應該做的)**

```sql
-- ✅ 使用交易
BEGIN;

CREATE TABLE users (...);
CREATE INDEX idx_users_email ON users(email);

COMMIT;

-- ✅ 使用 IF NOT EXISTS
CREATE TABLE IF NOT EXISTS users (...);

-- ✅ 一個 migration 一個目的
-- 000001_create_users.sql - 只建立 users 表
-- 000002_add_user_roles.sql - 只處理角色相關

-- ✅ 提供詳細註解
-- Migration: 000001
-- Description: Create users table for authentication
-- Author: John Doe
-- Date: 2025-01-18
-- Related: JIRA-123
```

**❌ DON'T (不應該做的)**

```sql
-- ❌ 不使用交易 (可能導致部分執行)
CREATE TABLE users (...);
CREATE TABLE roles (...);  -- 如果這裡失敗，users 已建立

-- ❌ 混合 schema 變更和資料遷移
CREATE TABLE users (...);
INSERT INTO users VALUES (...);  -- 應該分成兩個 migration

-- ❌ 不提供 down migration
-- 總是要有對應的 .down.sql 檔案

-- ❌ 使用 DROP TABLE (應該加 IF EXISTS)
DROP TABLE users;  -- 如果不存在會失敗
```

### 3. 測試流程

在執行 migration 之前，應該在測試環境驗證：

```bash
# 1. 在測試環境執行
export DB_DSN="postgres://test-db/..."
./migrate -command up

# 2. 驗證結果
./migrate -command version

# 3. 測試回滾
./migrate -command down

# 4. 再次執行 (確保可重複執行)
./migrate -command up

# 5. 確認無誤後，才在生產環境執行
export DB_DSN="postgres://prod-db/..."
./migrate -command up
```

### 4. 團隊協作

**Git 工作流程**:

```bash
# 1. 建立 feature branch
git checkout -b feature/add-user-roles

# 2. 建立 migration
migrate create -ext sql -dir migrations/auth add_user_roles

# 3. 撰寫 up 和 down SQL

# 4. 本地測試
./migrate -command up
# ... 測試 ...
./migrate -command down

# 5. Commit
git add migrations/auth/20250118120000_add_user_roles.*.sql
git commit -m "feat(auth): add user roles migration"

# 6. Push 並建立 PR
git push origin feature/add-user-roles
```

**Code Review 檢查清單**:

- [ ] Migration 檔名格式正確
- [ ] 同時有 up 和 down 檔案
- [ ] 使用交易包裹
- [ ] 使用 IF EXISTS / IF NOT EXISTS
- [ ] 有適當的註解說明
- [ ] Down migration 能正確回滾
- [ ] 已在測試環境驗證

### 5. 生產環境部署

**部署前檢查**:

```bash
# 1. 備份資料庫
pg_dump -h prod-db -U user dbname > backup_$(date +%Y%m%d_%H%M%S).sql

# 2. 檢查當前版本
./migrate -command version

# 3. 檢查待執行的 migrations
ls migrations/auth/*.up.sql

# 4. Dry-run (可選，需要工具支援)
# 在測試環境先執行一次

# 5. 執行 migration
./migrate -command up

# 6. 驗證
./migrate -command validate

# 7. 檢查應用程式是否正常
curl http://api/health
```

**回滾計畫**:

```bash
# 如果出問題，立即回滾
./migrate -command down

# 或回滾到特定版本
./migrate -command migrate -version 5

# 還原資料庫備份 (最後手段)
psql -h prod-db -U user dbname < backup_20250118_120000.sql
```

### 6. 監控與告警

**建議監控的項目**:

```go
// 在應用程式啟動時記錄 migration 狀態
version, dirty, _ := mgr.Version()
logger.Info("Database migration status",
    zap.Uint("version", version),
    zap.Bool("dirty", dirty),
)

// 如果是 dirty，發送告警
if dirty {
    alerting.Send("Database migration is dirty!", map[string]interface{}{
        "version": version,
        "service": "auth",
    })
}
```

### 7. 文檔記錄

為每個重要的 migration 建立文檔：

```markdown
# Migration 000005: 新增使用者角色系統

## 目的
實作 RBAC (角色基礎存取控制) 系統

## 變更內容
- 新增 `roles` 表
- 新增 `user_roles` 關聯表
- 新增必要的索引

## 影響範圍
- Auth Service
- 需要更新應用程式碼以支援新的角色系統

## 回滾影響
- 會刪除所有角色資料
- 需要重新配置使用者權限

## 相關連結
- JIRA: AUTH-123
- Design Doc: docs/rbac-design.md
```

## 實用工具

### 檢查 Migration 健康狀態

我們提供了一個腳本來檢查 migration 檔案的健康狀態：

```bash
# 檢查 auth service 的 migrations
./pkg/migrations/scripts/check_migrations.sh auth migrations/auth

# 輸出範例:
# ======================================
# Migration Health Check: auth
# ======================================
#
# ✅ Directory exists: migrations/auth
#
# 📋 Available migration files:
# ------------------------------
# 000001_create_users_table.up.sql
# 000001_create_users_table.down.sql
# 000002_add_email_verified.up.sql
# 000002_add_email_verified.down.sql
#
# 🔍 Checking up/down pairs:
# ------------------------------
# ✅ 000001 - has both up and down
# ✅ 000002 - has both up and down
#
# ...
#
# ✅ All checks passed!
```

### 列出所有版本

使用 Go 工具列出所有可用的 migration 版本：

```bash
cd pkg/migrations/scripts
go run list_versions.go -path ../../../migrations/auth

# 輸出範例:
# Version    Description                                        Up     Down
# ---------------------------------------------------------------------------------------------
# 1          create_users_table                                 ✅     ✅
# 2          add_email_verified                                 ✅     ✅
# 3          create_roles                                       ✅     ✅
#
# Summary:
# Total migrations: 3
# Complete (up+down): 3
# Incomplete: 0
#
# Version sequence check:
# ✅ Versions are sequential
```

### 檢視資料庫狀態

**SQL 查詢工具**:

```sql
-- 查看當前版本
SELECT * FROM schema_migrations_auth;

-- 查看所有服務的版本
SELECT
    'auth' as service,
    version,
    dirty,
    CASE WHEN dirty THEN '❌ DIRTY' ELSE '✅ Clean' END as status
FROM schema_migrations_auth
UNION ALL
SELECT
    'restaurant',
    version,
    dirty,
    CASE WHEN dirty THEN '❌ DIRTY' ELSE '✅ Clean' END
FROM schema_migrations_restaurant;

-- 列出所有 migration 相關的表
SELECT
    schemaname,
    tablename,
    tableowner
FROM pg_tables
WHERE tablename LIKE 'schema_migrations_%'
ORDER BY tablename;
```

### 批次操作腳本

**批次執行所有服務的 migrations**:

```bash
#!/bin/bash
# run_all_migrations.sh

set -e

SERVICES=("auth" "restaurant" "booking" "spider" "mail")
BASE_PATH="migrations"

for service in "${SERVICES[@]}"; do
    echo "Running migrations for $service..."

    DB_DSN="${service}_db_dsn"  # 從環境變數讀取

    ./migrate \
        -dsn "${!DB_DSN}" \
        -path "$BASE_PATH/$service" \
        -service "$service" \
        -command up

    echo "✅ $service migrations completed"
    echo ""
done

echo "🎉 All migrations completed!"
```

**批次檢查所有服務狀態**:

```bash
#!/bin/bash
# check_all_status.sh

SERVICES=("auth" "restaurant" "booking" "spider" "mail")
BASE_PATH="migrations"

echo "Migration Status Report"
echo "======================="
echo ""

for service in "${SERVICES[@]}"; do
    printf "%-15s: " "$service"

    DB_DSN="${service}_db_dsn"

    version_output=$(./migrate \
        -dsn "${!DB_DSN}" \
        -path "$BASE_PATH/$service" \
        -service "$service" \
        -command version 2>&1)

    echo "$version_output"
done
```

## 常見問題 (FAQ)

### Q1: 如何建立新的 migration?

**方法 1: 使用 migrate CLI (推薦)**

```bash
# 安裝
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 建立時間戳格式 (推薦)
migrate create -ext sql -dir migrations/auth create_roles_table
# 生成: 20250118120000_create_roles_table.up.sql
# 生成: 20250118120000_create_roles_table.down.sql

# 建立序列號格式
migrate create -ext sql -dir migrations/auth -seq create_roles_table
# 生成: 000001_create_roles_table.up.sql
# 生成: 000001_create_roles_table.down.sql
```

**方法 2: 手動建立**

```bash
# 使用當前時間戳
TIMESTAMP=$(date +%Y%m%d%H%M%S)
touch "migrations/auth/${TIMESTAMP}_create_roles_table.up.sql"
touch "migrations/auth/${TIMESTAMP}_create_roles_table.down.sql"
```

### Q2: 如何在多個環境中管理 migration?

使用環境變數配置資料庫連線字串，migration 檔案在所有環境保持一致：

```bash
# .env.development
DATABASE_URL=postgres://user:pass@localhost/mydb_dev

# .env.staging
DATABASE_URL=postgres://user:pass@staging-db/mydb

# .env.production
DATABASE_URL=postgres://user:pass@prod-db/mydb
```

### Q3: Migration 執行順序是如何決定的?

Migration 按照**檔案名稱的版本號**排序執行：

```
000001_xxx.sql  → 先執行
000002_xxx.sql  → 接著執行
000003_xxx.sql  → 最後執行
```

系統會：
1. 掃描 migrations 目錄
2. 解析檔案名稱中的版本號
3. 排序版本號
4. 執行大於當前資料庫版本的 migrations

### Q4: 可以跳過某個 migration 嗎?

**不建議跳過，但如果必須**:

```go
// 1. 先遷移到跳過前的版本
mgr.Migrate(ctx, 5)  // 假設要跳過版本 6

// 2. 手動在資料庫執行版本 6 的變更 (如果需要)

// 3. 強制設定為版本 6 (標記為已執行)
mgr.Force(6)

// 4. 繼續正常執行後續版本
mgr.Up(ctx)
```

### Q5: 如何測試 migration?

**單元測試範例**:

```go
func TestMigration(t *testing.T) {
    // 使用 testcontainers 建立測試資料庫
    db := setupTestDB(t)
    defer db.Close()

    logger, _ := zap.NewDevelopment()

    mgr, err := migrations.NewManager(migrations.Config{
        DB:             db,
        Logger:         logger,
        MigrationsPath: "file://migrations/auth",
        ServiceName:    "auth",
    })
    require.NoError(t, err)
    defer mgr.Close()

    ctx := context.Background()

    // 測試 up
    err = mgr.Up(ctx)
    assert.NoError(t, err)

    // 驗證版本
    version, dirty, err := mgr.Version()
    assert.NoError(t, err)
    assert.False(t, dirty)
    assert.Greater(t, version, uint(0))

    // 測試 down
    err = mgr.Down(ctx)
    assert.NoError(t, err)

    // 驗證可以再次 up
    err = mgr.Up(ctx)
    assert.NoError(t, err)
}
```

### Q6: Migration 失敗了該怎麼辦?

參考 [疑難排解 - Dirty State 處理](#dirty-state-處理) 章節。

### Q7: 如何在 CI/CD 中執行 migrations?

**GitHub Actions 範例**:

```yaml
name: Run Migrations

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  migrate:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v3

      - name: Install migrate
        run: |
          curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
          sudo mv migrate /usr/local/bin/

      - name: Run migrations
        env:
          DATABASE_URL: postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable
        run: |
          migrate -path migrations/auth -database "$DATABASE_URL" up

      - name: Verify migration
        run: |
          migrate -path migrations/auth -database "$DATABASE_URL" version
```

### Q8: 多個開發者同時建立 migration 怎麼辦?

使用**時間戳格式**可以避免版本號衝突：

```bash
# 開發者 A
migrate create -ext sql -dir migrations/auth add_user_roles
# 生成: 20250118120000_add_user_roles.sql

# 開發者 B (同一時間)
migrate create -ext sql -dir migrations/auth add_permissions
# 生成: 20250118120001_add_permissions.sql (時間戳不同)
```

合併時按照時間戳順序執行，自然不會衝突。

## 參考資料

### 官方文檔

- [golang-migrate/migrate](https://github.com/golang-migrate/migrate) - 官方 GitHub 倉庫
- [Migration Best Practices](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md) - 最佳實踐指南
- [PostgreSQL Documentation](https://www.postgresql.org/docs/) - PostgreSQL 官方文檔

### 相關文章

- [Database Migrations Done Right](https://www.brunobrito.pt/database-migrations/)
- [Evolutionary Database Design](https://martinfowler.com/articles/evodb.html) - Martin Fowler
- [Microservices and Database per Service](https://microservices.io/patterns/data/database-per-service.html)

### 工具

- [migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) - 命令列工具
- [Atlas](https://atlasgo.io/) - 替代方案，提供更多功能
- [Flyway](https://flywaydb.org/) - Java 生態系的 migration 工具 (概念類似)
