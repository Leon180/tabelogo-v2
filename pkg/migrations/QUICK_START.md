# Migration 快速開始指南

## 🚀 5 分鐘快速測試

### 1. 建置 CLI 工具

```bash
cd /Users/lileon/goproject/tabelogov2/pkg
export GOWORK=off
go build -o migrations/bin/migrate ./migrations/cmd/migrate/main.go
```

### 2. 準備測試資料庫（選擇一種）

**選項 A: Docker (推薦)**

```bash
docker run -d \
  --name test-postgres \
  -e POSTGRES_USER=testuser \
  -e POSTGRES_PASSWORD=testpass \
  -e POSTGRES_DB=testdb \
  -p 5433:5432 \
  postgres:15-alpine
```

**選項 B: 本地 PostgreSQL**

```sql
CREATE DATABASE testdb;
```

### 3. 執行 Migration

```bash
cd /Users/lileon/goproject/tabelogov2/pkg/migrations

# 設定資料庫連線
export DB_DSN="postgres://testuser:testpass@localhost:5433/testdb?sslmode=disable"

# 執行 up
./bin/migrate \
  -dsn "$DB_DSN" \
  -path "../../migrations/auth" \
  -service "auth" \
  -command up
```

### 4. 驗證結果

```bash
# 檢查版本
./bin/migrate -dsn "$DB_DSN" -path "../../migrations/auth" -service "auth" -command version

# 預期輸出: Current version: 2
```

### 5. 查看資料庫

```bash
# 使用 Docker
docker exec -it test-postgres psql -U testuser -d testdb -c "\dt"
docker exec -it test-postgres psql -U testuser -d testdb -c "SELECT * FROM schema_migrations_auth"

# 或使用 psql
psql -h localhost -p 5433 -U testuser -d testdb -c "\dt"
```

## 📋 常用命令

### 執行 Up

```bash
./bin/migrate -dsn "$DB_DSN" -path "../../migrations/auth" -service "auth" -command up
```

### 回滾 (Down)

```bash
./bin/migrate -dsn "$DB_DSN" -path "../../migrations/auth" -service "auth" -command down
```

### 執行指定步數

```bash
# 向上 2 步
./bin/migrate -dsn "$DB_DSN" -path "../../migrations/auth" -service "auth" -command steps -steps 2

# 向下 1 步
./bin/migrate -dsn "$DB_DSN" -path "../../migrations/auth" -service "auth" -command steps -steps -1
```

### 檢查版本

```bash
./bin/migrate -dsn "$DB_DSN" -path "../../migrations/auth" -service "auth" -command version
```

### 驗證狀態

```bash
./bin/migrate -dsn "$DB_DSN" -path "../../migrations/auth" -service "auth" -command validate
```

## 🧪 使用測試腳本

### 基本測試（不需要資料庫）

```bash
cd /Users/lileon/goproject/tabelogov2/pkg/migrations
./test_simple.sh
```

### 完整測試（需要 Docker）

```bash
cd /Users/lileon/goproject/tabelogov2/pkg/migrations
./test_migrations.sh
```

## 🔧 在程式碼中使用

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
    db, err := sql.Open("postgres", "postgres://user:pass@localhost/db?sslmode=disable")
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

    // 取得版本
    version, dirty, _ := mgr.Version()
    logger.Info("Migration completed",
        zap.Uint("version", version),
        zap.Bool("dirty", dirty),
    )
}
```

## 📚 更多資訊

- [完整文檔](README.md)
- [測試結果](TEST_RESULTS.md)
- [Import 指南](IMPORT_GUIDE.md)
- [修正記錄](FIXED_IMPORTS.md)

## ⚠️ 注意事項

1. **生產環境**: 執行前務必備份資料庫
2. **版本號**: 建議使用時間戳格式避免衝突
3. **Dirty State**: 如果遇到 dirty 狀態，參考[完整文檔](README.md#疑難排解)
4. **Docker**: 記得在測試完成後清理容器

```bash
docker stop test-postgres && docker rm test-postgres
```
