# Import 問題已修正 ✅

## 問題描述

原本 `pkg/migrations/` 有獨立的 `go.mod`，導致 import 錯誤和依賴管理複雜。

## 解決方案

### 1. 移除獨立的 go.mod

已刪除:
- ~~`pkg/migrations/go.mod`~~
- ~~`pkg/migrations/go.sum`~~

### 2. 統一在 pkg 層級管理

現在所有依賴在 `pkg/go.mod` 中統一管理:

```
pkg/
├── go.mod          ← 唯一的 module 定義
├── go.sum          ← 依賴鎖定
├── migrations/     ← package migrations (無獨立 go.mod)
├── logger/
├── config/
└── ...
```

### 3. 更新的依賴

`pkg/go.mod` 現在包含:

```go
module github.com/Leon180/tabelogo-v2/pkg

require (
    github.com/gin-gonic/gin v1.10.0
    github.com/golang-migrate/migrate/v4 v4.17.0  ← migrations 需要
    github.com/lib/pq v1.10.9                      ← migrations 需要
    go.uber.org/fx v1.20.1
    go.uber.org/zap v1.26.0
)
```

## 如何使用

### Import 方式

```go
// ✅ 正確
import "github.com/Leon180/tabelogo-v2/pkg/migrations"

// ❌ 錯誤 (舊的方式)
import "github.com/Leon180/tabelogo-v2/pkg/migrations/migrations"
```

### 在專案中使用

```go
package main

import (
    "context"
    "database/sql"

    "github.com/Leon180/tabelogo-v2/pkg/migrations"
    "github.com/Leon180/tabelogo-v2/pkg/logger"
    _ "github.com/lib/pq"
)

func main() {
    log := logger.New()

    db, err := sql.Open("postgres", "postgres://...")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    mgr, err := migrations.NewManager(migrations.Config{
        DB:             db,
        Logger:         log,
        MigrationsPath: "file://migrations/auth",
        ServiceName:    "auth",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer mgr.Close()

    if err := mgr.Up(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

### 使用 Uber FX

```go
import (
    "github.com/Leon180/tabelogo-v2/pkg/migrations"
    "go.uber.org/fx"
)

func main() {
    app := fx.New(
        // 提供依賴
        fx.Provide(
            NewDatabase,
            NewLogger,
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

        // 啟動時自動執行
        migrations.InvokeAutoMigrate(),
    )

    app.Run()
}
```

## 驗證

### 方法 1: 編譯檢查

```bash
cd /Users/lileon/goproject/tabelogov2/pkg
go build ./migrations
```

應該沒有任何錯誤。

### 方法 2: 執行測試

```bash
cd /Users/lileon/goproject/tabelogov2/pkg/migrations
go test -v
```

### 方法 3: 使用驗證腳本

```bash
cd /Users/lileon/goproject/tabelogov2/pkg
./verify_imports.sh
```

## 常見問題

### Q: IDE 還是顯示錯誤?

A: 重新載入 Go module cache:

1. 在 pkg 目錄執行:
   ```bash
   cd pkg
   go mod download
   go mod tidy
   ```

2. 重啟 IDE 或重新載入工作區

### Q: 如何添加新的依賴?

A: 在 pkg 目錄執行:

```bash
cd pkg
go get github.com/new/package@version
```

### Q: 可以在 migrations 下執行 go test 嗎?

A: 可以! Go 會自動找到上層的 go.mod:

```bash
cd pkg/migrations
go test -v  # ✅ 正常工作
```

## 相關文件

- [pkg/README.md](../README.md) - pkg 整體說明
- [pkg/migrations/README.md](README.md) - Migrations 完整文檔
- [pkg/migrations/IMPORT_GUIDE.md](IMPORT_GUIDE.md) - Import 使用指南

## 修正記錄

- **2025-01-18**: 移除獨立的 `go.mod`，統一在 pkg 層級管理
- **2025-01-18**: 更新所有文檔和範例
- **2025-01-18**: 建立驗證腳本
