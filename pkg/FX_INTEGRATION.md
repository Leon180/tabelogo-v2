# Uber FX 依賴注入整合指南

本專案的共用套件完全支援 Uber FX 依賴注入框架。

## 🎯 為什麼使用 Uber FX？

Uber FX 是 Uber 開源的依賴注入框架，被許多大型 Go 專案使用：

✅ **自動依賴解析** - 不需要手動管理依賴順序
✅ **生命週期管理** - 自動處理啟動/關閉邏輯
✅ **測試友善** - 輕鬆 mock 依賴
✅ **模組化** - 清晰的程式碼組織
✅ **型別安全** - 編譯期檢查依賴

---

## 📦 已提供的 FX Modules

### 1. config.Module
```go
import "github.com/Leon180/tabelogo-v2/pkg/config"

fx.New(
    config.Module,  // 自動提供 *config.Config
)
```

### 2. logger.Module
```go
import "github.com/Leon180/tabelogo-v2/pkg/logger"

fx.New(
    logger.Module,  // 自動提供 *zap.Logger
)
```

---

## 🚀 快速開始

### 基礎範例
```go
package main

import (
    "github.com/Leon180/tabelogo-v2/pkg/config"
    "github.com/Leon180/tabelogo-v2/pkg/logger"
    "go.uber.org/fx"
    "go.uber.org/zap"
)

func main() {
    fx.New(
        // 載入模組
        config.Module,
        logger.Module,

        // 使用依賴
        fx.Invoke(func(cfg *config.Config, log *zap.Logger) {
            log.Info("Application starting",
                zap.String("environment", cfg.Environment),
                zap.Int("port", cfg.ServerPort),
            )
        }),
    ).Run()
}
```

---

## 📚 使用模式

### Pattern 1: 簡單服務

```go
type UserService struct {
    config *config.Config
    logger *zap.Logger
}

// FX 會自動注入依賴
func NewUserService(cfg *config.Config, log *zap.Logger) *UserService {
    return &UserService{
        config: cfg,
        logger: log,
    }
}

func (s *UserService) GetUser(id string) {
    s.logger.Info("Getting user", zap.String("id", id))
    // ... business logic
}

func main() {
    fx.New(
        config.Module,
        logger.Module,

        // 提供你的服務
        fx.Provide(NewUserService),

        // 使用服務
        fx.Invoke(func(svc *UserService) {
            svc.GetUser("123")
        }),
    ).Run()
}
```

### Pattern 2: 生命週期管理

```go
type Server struct {
    config *config.Config
    logger *zap.Logger
}

func NewServer(cfg *config.Config, log *zap.Logger) *Server {
    return &Server{config: cfg, logger: log}
}

func (s *Server) Start(ctx context.Context) error {
    s.logger.Info("Server starting", zap.Int("port", s.config.ServerPort))
    // 啟動 HTTP 服務器
    return nil
}

func (s *Server) Stop(ctx context.Context) error {
    s.logger.Info("Server stopping")
    // 優雅關閉
    return nil
}

func main() {
    fx.New(
        config.Module,
        logger.Module,
        fx.Provide(NewServer),

        // 註冊生命週期 Hook
        fx.Invoke(func(lc fx.Lifecycle, server *Server) {
            lc.Append(fx.Hook{
                OnStart: server.Start,
                OnStop:  server.Stop,
            })
        }),
    ).Run()
}
```

### Pattern 3: 依賴鏈

FX 會自動解析依賴順序：

```go
// Database 依賴 Config 和 Logger
type Database struct {
    config *config.Config
    logger *zap.Logger
}

func NewDatabase(cfg *config.Config, log *zap.Logger) *Database {
    log.Info("Connecting to database", zap.String("host", cfg.Database.Host))
    return &Database{config: cfg, logger: log}
}

// Repository 依賴 Database
type UserRepository struct {
    db     *Database
    logger *zap.Logger
}

func NewUserRepository(db *Database, log *zap.Logger) *UserRepository {
    return &UserRepository{db: db, logger: log}
}

// Service 依賴 Repository
type UserService struct {
    repo   *UserRepository
    logger *zap.Logger
}

func NewUserService(repo *UserRepository, log *zap.Logger) *UserService {
    return &UserService{repo: repo, logger: log}
}

func main() {
    fx.New(
        config.Module,
        logger.Module,

        // FX 會自動按正確順序建立：Config → Logger → Database → Repository → Service
        fx.Provide(
            NewDatabase,
            NewUserRepository,
            NewUserService,
        ),

        fx.Invoke(func(svc *UserService) {
            // Service 已經準備好，所有依賴都已注入
        }),
    ).Run()
}
```

### Pattern 4: 模組化 (推薦)

將相關的 Providers 組織成模組：

```go
// internal/user/module.go
package user

import (
    "github.com/Leon180/tabelogo-v2/pkg/config"
    "github.com/Leon180/tabelogo-v2/pkg/logger"
    "go.uber.org/fx"
)

var Module = fx.Module("user",
    // 包含基礎模組
    config.Module,
    logger.Module,

    // 提供此模組的服務
    fx.Provide(
        NewDatabase,
        NewUserRepository,
        NewUserService,
    ),
)

// cmd/user-service/main.go
package main

import (
    "github.com/Leon180/tabelogo-v2/internal/user"
    "go.uber.org/fx"
)

func main() {
    fx.New(
        // 只需要載入模組
        user.Module,

        fx.Invoke(func(svc *user.UserService) {
            // Ready to use
        }),
    ).Run()
}
```

---

## 🧪 測試

### 單元測試 - Mock 依賴

```go
func TestUserService(t *testing.T) {
    // 建立 mock config
    mockConfig := &config.Config{
        Environment: "test",
        ServerPort:  8888,
    }

    // 建立 test logger
    mockLogger, _ := zap.NewDevelopment()

    // 使用 FX 測試
    var svc *UserService

    app := fxtest.New(t,
        // 提供 mocks
        fx.Supply(mockConfig),
        fx.Supply(mockLogger),

        // 提供要測試的服務
        fx.Provide(NewUserService),

        // 取得實例
        fx.Populate(&svc),
    )

    app.RequireStart()
    defer app.RequireStop()

    // 測試你的服務
    svc.GetUser("123")
}
```

### 整合測試

```go
func TestIntegration(t *testing.T) {
    app := fxtest.New(t,
        config.Module,
        logger.Module,

        // 使用真實服務
        fx.Provide(
            NewDatabase,
            NewUserRepository,
            NewUserService,
        ),

        fx.Invoke(func(svc *UserService) {
            // 執行整合測試
            user := svc.GetUser("123")
            assert.NotNil(t, user)
        }),
    )

    app.RequireStart()
    defer app.RequireStop()
}
```

---

## 🎨 進階用法

### 條件性提供 (根據環境)

```go
func main() {
    fx.New(
        config.Module,

        // 根據環境提供不同的 logger
        fx.Provide(func(cfg *config.Config) (*zap.Logger, error) {
            if cfg.IsDevelopment() {
                return logger.NewDevelopment()
            }
            return logger.NewProduction()
        }),

        fx.Invoke(func(log *zap.Logger) {
            log.Info("Logger configured based on environment")
        }),
    ).Run()
}
```

### 選擇性依賴 (Optional)

```go
type MyService struct {
    fx.In

    Config   *config.Config
    Logger   *zap.Logger
    Cache    *redis.Client `optional:"true"` // 可選依賴
}

func NewMyService(params MyService) *MyService {
    if params.Cache != nil {
        params.Logger.Info("Cache is available")
    } else {
        params.Logger.Info("Running without cache")
    }
    return &MyService{}
}
```

### 多個相同類型的依賴 (Named)

```go
type Result struct {
    fx.Out

    AuthDB      *gorm.DB `name:"auth"`
    RestaurantDB *gorm.DB `name:"restaurant"`
}

func NewDatabases(cfg *config.Config) (Result, error) {
    authDB := // connect to auth_db
    restaurantDB := // connect to restaurant_db

    return Result{
        AuthDB: authDB,
        RestaurantDB: restaurantDB,
    }, nil
}

type UserService struct {
    fx.In

    DB *gorm.DB `name:"auth"` // 注入 auth DB
}
```

---

## 📖 與現有程式碼的對比

### ❌ 傳統方式（手動管理依賴）

```go
func main() {
    // 手動載入 config
    cfg, err := config.Load()
    if err != nil {
        panic(err)
    }

    // 手動初始化 logger
    log, err := logger.New(cfg.LogLevel)
    if err != nil {
        panic(err)
    }
    defer log.Sync()

    // 手動建立 database
    db := NewDatabase(cfg, log)

    // 手動建立 repository
    repo := NewUserRepository(db, log)

    // 手動建立 service
    svc := NewUserService(repo, log)

    // 手動啟動
    if err := svc.Start(); err != nil {
        panic(err)
    }

    // 手動關閉（容易忘記）
    defer svc.Stop()
}
```

### ✅ FX 方式（自動管理）

```go
func main() {
    fx.New(
        config.Module,
        logger.Module,

        fx.Provide(
            NewDatabase,
            NewUserRepository,
            NewUserService,
        ),

        fx.Invoke(func(lc fx.Lifecycle, svc *UserService) {
            lc.Append(fx.Hook{
                OnStart: func(ctx context.Context) error {
                    return svc.Start()
                },
                OnStop: func(ctx context.Context) error {
                    return svc.Stop()
                },
            })
        }),
    ).Run() // 自動處理啟動、執行、優雅關閉
}
```

---

## ✅ 優點總結

| 特性 | 傳統方式 | FX 方式 |
|------|---------|---------|
| 依賴順序 | 手動管理 | ✅ 自動解析 |
| 生命週期 | 手動處理 | ✅ 自動管理 |
| 錯誤處理 | 需要大量 if err | ✅ 集中處理 |
| 測試 | Mock 複雜 | ✅ 輕鬆 Mock |
| 程式碼量 | 冗長 | ✅ 簡潔 |
| 可維護性 | 低 | ✅ 高 |

---

## 📚 延伸閱讀

1. [Uber FX 官方文檔](https://uber-go.github.io/fx/)
2. [Uber Go Style Guide](https://github.com/uber-go/guide)
3. [FX 最佳實踐](https://github.com/uber-go/fx/blob/master/docs/best-practices.md)

---

## 💡 建議

**使用 FX 當：**
- ✅ 專案有多個服務/模組
- ✅ 需要管理複雜的依賴關係
- ✅ 需要生命週期管理（啟動/關閉）
- ✅ 需要良好的測試性

**不使用 FX 當：**
- ❌ 簡單的 CLI 工具
- ❌ 單一檔案的腳本
- ❌ 依賴關係非常簡單

**對於本專案（微服務架構）：強烈建議使用 FX！**
