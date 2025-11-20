# Migration 功能測試結果

## ✅ 測試完成日期
2025-01-19

## 測試環境

- **Go Version**: 1.21
- **操作系統**: macOS
- **測試範圍**: Migration 核心功能

## 測試結果總覽

| 項目 | 狀態 | 說明 |
|------|------|------|
| 程式碼編譯 | ✅ 通過 | 所有套件可正常編譯 |
| CLI 工具建置 | ✅ 通過 | migrate CLI 成功建置 (7.8MB) |
| 幫助資訊 | ✅ 通過 | CLI 可正常顯示幫助 |
| Import 錯誤修正 | ✅ 完成 | 統一在 pkg/go.mod 管理 |

## 已測試的功能

### 1. 建置系統 ✅

```bash
cd /Users/lileon/goproject/tabelogov2/pkg
go build -o bin/migrate ./migrations/cmd/migrate/main.go
```

**結果**: 成功建置，產生 7.8MB 的可執行檔

### 2. CLI 工具 ✅

執行 `./migrate` 顯示正確的幫助資訊：

```
Commands:
  up       - Run all pending migrations
  down     - Rollback the last migration
  steps    - Run N migrations (positive=up, negative=down)
  migrate  - Migrate to a specific version
  version  - Show current migration version
  validate - Validate migration state
  force    - Force set migration version (use with caution)
  drop     - Drop all tables (dangerous!)
```

### 3. Migration 檔案 ✅

已建立測試用的 migration 檔案：

- `migrations/auth/000001_create_users_table.up.sql`
- `migrations/auth/000001_create_users_table.down.sql`
- `migrations/auth/000002_create_refresh_tokens_table.up.sql`
- `migrations/auth/000002_create_refresh_tokens_table.down.sql`

檔案格式正確，包含：
- CREATE TABLE 語句
- 索引建立
- 觸發器 (updated_at)
- 註解 (COMMENT ON)

### 4. 程式碼結構 ✅

```
pkg/
├── go.mod               ✅ 統一管理
├── go.sum               ✅ 依賴鎖定
└── migrations/
    ├── manager.go       ✅ 核心管理器
    ├── fx.go            ✅ FX 整合
    ├── errors.go        ✅ 錯誤定義
    ├── cmd/migrate/     ✅ CLI 工具
    └── test_simple.sh   ✅ 測試腳本
```

## 需要資料庫的測試

以下測試需要 PostgreSQL 資料庫：

### 準備測試資料庫

**選項 1: 使用 Docker**

```bash
docker run -d \
  --name test-postgres \
  -e POSTGRES_USER=testuser \
  -e POSTGRES_PASSWORD=testpass \
  -e POSTGRES_DB=testdb \
  -p 5433:5432 \
  postgres:15-alpine
```

**選項 2: 使用本地 PostgreSQL**

```sql
CREATE DATABASE testdb;
CREATE USER testuser WITH PASSWORD 'testpass';
GRANT ALL PRIVILEGES ON DATABASE testdb TO testuser;
```

### 執行測試

```bash
# 設定環境變數
export DB_DSN="postgres://testuser:testpass@localhost:5433/testdb?sslmode=disable"

# 建置 CLI
cd /Users/lileon/goproject/tabelogov2/pkg
go build -o migrations/bin/migrate ./migrations/cmd/migrate/main.go

# 執行 migrations
cd migrations
./bin/migrate \
  -dsn "$DB_DSN" \
  -path "../../migrations/auth" \
  -service "auth" \
  -command up

# 檢查版本
./bin/migrate \
  -dsn "$DB_DSN" \
  -path "../../migrations/auth" \
  -service "auth" \
  -command version

# 驗證狀態
./bin/migrate \
  -dsn "$DB_DSN" \
  -path "../../migrations/auth" \
  -service "auth" \
  -command validate
```

### 預期結果

**Up Migration**:
- 建立 `users` 表
- 建立 `refresh_tokens` 表
- 建立所有索引
- 建立 `update_updated_at_column` 函數
- 建立觸發器
- 版本記錄在 `schema_migrations_auth` 表

**版本檢查**:
```
Current version: 2
```

**驗證狀態**:
```
Migration state is valid
```

## 手動驗證 SQL

連接到資料庫後：

```sql
-- 查看所有表
\dt

-- 查看版本控制表
SELECT * FROM schema_migrations_auth;

-- 查看 users 表結構
\d users

-- 查看索引
\di

-- 查看觸發器
SELECT tgname FROM pg_trigger WHERE tgrelid = 'users'::regclass;
```

## 已知問題與解決方案

### 問題 1: Docker Daemon 未運行

**錯誤**:
```
Cannot connect to the Docker daemon at unix:///Users/lileon/.docker/run/docker.sock
```

**解決方案**:
1. 啟動 Docker Desktop
2. 或使用本地 PostgreSQL 進行測試

### 問題 2: Import 錯誤

**已修正**: 移除 `pkg/migrations/go.mod`，統一在 `pkg/go.mod` 管理

### 問題 3: Go 版本格式

**已修正**: 將 `go 1.23.0` 改為 `go 1.21`，移除 `toolchain` 指令

## 測試腳本

已提供以下測試腳本：

1. **test_simple.sh**: 基本建置和 CLI 測試
   ```bash
   cd pkg/migrations
   ./test_simple.sh
   ```

2. **test_migrations.sh**: 完整測試（需要 Docker）
   ```bash
   cd pkg/migrations
   ./test_migrations.sh
   ```

3. **check_migrations.sh**: 檢查 migration 檔案健康狀態
   ```bash
   cd pkg/migrations/scripts
   ./check_migrations.sh auth ../../migrations/auth
   ```

## 下一步測試建議

1. **整合測試**: 在實際的 PostgreSQL 環境中測試
2. **Up/Down 循環測試**: 測試多次 up/down 的穩定性
3. **並行測試**: 測試多個服務同時執行 migration
4. **錯誤恢復測試**: 測試 dirty state 的修復
5. **FX 整合測試**: 在實際微服務中測試 FX 自動執行

## 結論

✅ **Migration 核心功能已實作完成並通過基本測試**

- 程式碼可正常編譯
- CLI 工具運作正常
- Migration 檔案格式正確
- 文檔完整

需要連接實際資料庫才能測試完整的 migration 執行流程。
