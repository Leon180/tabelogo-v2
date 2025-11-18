# Database Migrations

本專案的資料庫 migration 管理說明。

## 目錄結構

```
migrations/
├── auth/                    # Auth Service 的 migrations
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_refresh_tokens_table.up.sql
│   └── 000002_create_refresh_tokens_table.down.sql
├── restaurant/              # Restaurant Service 的 migrations
│   ├── 000001_create_restaurants_table.up.sql
│   └── 000001_create_restaurants_table.down.sql
├── booking/                 # Booking Service 的 migrations
├── spider/                  # Spider Service 的 migrations
└── mail/                    # Mail Service 的 migrations
```

## 快速開始

### 1. 安裝 migrate CLI

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### 2. 建立新的 Migration

```bash
# 時間戳格式 (推薦)
migrate create -ext sql -dir migrations/auth create_new_table

# 序列號格式
migrate create -ext sql -dir migrations/auth -seq create_new_table
```

### 3. 執行 Migrations

**使用 CLI 工具**:

```bash
# 設定環境變數
export DB_DSN="postgres://user:pass@localhost/auth_db?sslmode=disable"
export MIGRATIONS_PATH="migrations/auth"
export SERVICE_NAME="auth"

# 執行所有未執行的 migrations
cd pkg/migrations
make example-up

# 或直接使用 CLI
go run cmd/migrate/main.go \
    -dsn "$DB_DSN" \
    -path "$MIGRATIONS_PATH" \
    -service "$SERVICE_NAME" \
    -command up
```

**使用程式碼**:

```go
import (
    "github.com/Leon180/tabelogo-v2/pkg/migrations"
)

mgr, err := migrations.NewManager(migrations.Config{
    DB:             db,
    Logger:         logger,
    MigrationsPath: "file://migrations/auth",
    ServiceName:    "auth",
})
defer mgr.Close()

// 執行 migrations
err = mgr.Up(context.Background())
```

### 4. 檢查 Migration 狀態

```bash
# 使用提供的檢查腳本
./pkg/migrations/scripts/check_migrations.sh auth migrations/auth

# 查看版本
psql -d auth_db -c "SELECT * FROM schema_migrations_auth;"

# 列出所有版本
cd pkg/migrations/scripts
go run list_versions.go -path ../../../migrations/auth
```

## 版本控制機制

### 版本號來源

Migration 的版本號來自**檔案名稱的前綴數字**:

```
000001_create_users.up.sql    → 版本: 1
000002_add_email.up.sql        → 版本: 2
20250118120000_add_roles.up.sql → 版本: 20250118120000
```

### 查看當前版本

**方法 1: SQL 查詢**

```sql
-- 查看 auth service 的版本
SELECT * FROM schema_migrations_auth;

-- 輸出:
-- version | dirty
-- --------|------
-- 2       | false
```

**方法 2: 使用程式**

```go
version, dirty, err := mgr.Version()
fmt.Printf("Version: %d, Dirty: %v\n", version, dirty)
```

**方法 3: 使用 CLI**

```bash
./migrate -dsn "..." -path migrations/auth -service auth -command version
```

### 查看可用的 Migrations

**檔案系統**:

```bash
ls -1 migrations/auth/*.up.sql | sed 's/.*\/\([0-9]*\)_.*/\1/'
# 輸出:
# 000001
# 000002
```

**使用工具**:

```bash
cd pkg/migrations/scripts
go run list_versions.go -path ../../../migrations/auth
```

## Migration 檔案格式

### 檔案命名

```
{version}_{description}.{up|down}.sql

範例:
000001_create_users_table.up.sql     ✅ 正確
000001_create_users_table.down.sql   ✅ 正確
20250118120000_add_roles.up.sql      ✅ 正確 (時間戳格式)
001-create-users.sql                 ❌ 錯誤 (格式不對)
create_users.up.sql                  ❌ 錯誤 (缺少版本號)
```

### 檔案內容範例

**Up Migration** (`000001_create_users.up.sql`):

```sql
-- Migration: 000001
-- Description: Create users table
-- Author: Your Name
-- Date: 2025-01-18

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);

COMMENT ON TABLE users IS 'User accounts';
```

**Down Migration** (`000001_create_users.down.sql`):

```sql
-- Rollback: 000001
-- Description: Drop users table

DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
```

## 版本狀態

### Clean State (正常)

```
schema_migrations_auth
┌─────────┬───────┐
│ version │ dirty │
├─────────┼───────┤
│    2    │ false │  ✅
└─────────┴───────┘
```

### Dirty State (異常)

```
schema_migrations_auth
┌─────────┬───────┐
│ version │ dirty │
├─────────┼───────┤
│    2    │ true  │  ❌ 需要修復!
└─────────┴───────┘
```

**修復方法**: 參考 [pkg/migrations/README.md](../pkg/migrations/README.md#疑難排解)

## 各服務的資料庫

| 服務 | 資料庫名稱 | 版本控制表 |
|------|-----------|-----------|
| Auth | `auth_db` | `schema_migrations_auth` |
| Restaurant | `restaurant_db` | `schema_migrations_restaurant` |
| Booking | `booking_db` | `schema_migrations_booking` |
| Spider | `spider_db` | `schema_migrations_spider` |
| Mail | `mail_db` | `schema_migrations_mail` |

## 執行流程

當執行 `Up()` 時:

```
1. 掃描 migrations/auth/ 目錄
   → 找到: [1, 2, 3, 4, 5]

2. 查詢資料庫當前版本
   → SELECT version FROM schema_migrations_auth
   → 結果: 2

3. 計算需要執行的 migrations
   → 需要執行: [3, 4, 5]

4. 按順序執行
   → 執行 3 ✅
   → 執行 4 ✅
   → 執行 5 ✅

5. 更新版本記錄
   → UPDATE schema_migrations_auth SET version = 5
```

## 最佳實踐

1. ✅ **使用時間戳格式版本號** - 避免多人開發衝突
2. ✅ **總是提供 down migration** - 確保可回滾
3. ✅ **使用 IF NOT EXISTS / IF EXISTS** - 避免重複執行錯誤
4. ✅ **一個 migration 只做一件事** - 保持簡單
5. ✅ **在測試環境先驗證** - 不要直接在生產環境執行
6. ✅ **添加詳細註解** - 說明目的和相關資訊
7. ✅ **執行前備份資料庫** - 以防萬一

## 工具與腳本

### Migration 健康檢查

```bash
./pkg/migrations/scripts/check_migrations.sh auth migrations/auth
```

檢查項目:
- ✅ Up/Down 檔案配對
- ✅ 版本號重複
- ✅ 檔名格式
- ✅ SQL 語法 (基本)

### 列出版本

```bash
cd pkg/migrations/scripts
go run list_versions.go -path ../../../migrations/auth
```

## 常見問題

### Q: Migration 執行失敗怎麼辦?

參考 [pkg/migrations/README.md - 疑難排解](../pkg/migrations/README.md#疑難排解) 章節。

### Q: 如何回滾 Migration?

```bash
# 回滾最後一個
./migrate -command down

# 回滾多步
./migrate -command steps -steps -2

# 回滾到特定版本
./migrate -command migrate -version 3
```

### Q: 如何在多個環境管理 Migrations?

使用環境變數區分資料庫連線，migration 檔案保持一致:

```bash
# Development
export DB_DSN="postgres://localhost/mydb_dev"

# Production
export DB_DSN="postgres://prod-db/mydb"
```

## 更多資訊

詳細的使用說明、API 參考、進階主題和疑難排解，請參考:

📖 **[pkg/migrations/README.md](../pkg/migrations/README.md)**

包含:
- 完整的 API 參考
- 版本控制機制詳解
- Dirty State 處理
- 多服務管理
- 環境隔離
- 最佳實踐
- 常見問題與解答
- 參考資料
