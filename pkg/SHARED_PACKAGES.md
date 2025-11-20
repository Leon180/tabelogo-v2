# 共用套件文檔 (Shared Packages Documentation)

本文檔說明專案中所有共用套件的使用方式。

---

## 📦 套件總覽

| 套件 | 說明 | 狀態 |
|------|------|------|
| `pkg/logger` | 統一日誌套件 (基於 zap) | ✅ 完成 |
| `pkg/config` | 配置載入與管理 | ✅ 完成 |
| `pkg/errors` | 統一錯誤處理 | ✅ 完成 |
| `pkg/migrations` | 資料庫 Migration 管理 | ✅ 完成 |
| `pkg/middleware` | HTTP/gRPC 中間件 | 🚧 開發中 |

---

## 1️⃣ Logger 套件 (`pkg/logger`)

### 功能特性

✅ 基於 `go.uber.org/zap` 的高效能日誌
✅ 支援多個日誌等級 (debug, info, warn, error, fatal)
✅ 結構化日誌 (JSON 格式)
✅ Context 支援（日誌追蹤）
✅ Uber FX 依賴注入整合

### 基本使用

```go
package main

import (
    "github.com/lileon/tabelogov2/pkg/logger"
    "go.uber.org/zap"
)

func main() {
    // 初始化日誌
    if err := logger.Init("info"); err != nil {
        panic(err)
    }
    defer logger.Sync()

    // 基本日誌
    logger.Info("application started")
    logger.Debug("debug message")
    logger.Warn("warning message")
    logger.Error("error message")

    // 帶欄位的日誌
    logger.Info("user logged in",
        zap.String("user_id", "123"),
        zap.String("ip", "192.168.1.1"),
    )

    // 建立子 logger
    childLogger := logger.With(
        zap.String("service", "auth"),
        zap.String("version", "1.0.0"),
    )
    childLogger.Info("service started")
}
```

### Context 支援

```go
import (
    "context"
    "github.com/lileon/tabelogov2/pkg/logger"
    "go.uber.org/zap"
)

func handleRequest(ctx context.Context) {
    // 從 context 獲取 logger
    log := logger.FromContext(ctx)
    log.Info("processing request")

    // 添加 fields 到 context
    ctx = logger.WithFields(ctx,
        zap.String("request_id", "req-123"),
        zap.String("user_id", "user-456"),
    )

    // 在其他函數中使用
    processData(ctx)
}

func processData(ctx context.Context) {
    log := logger.FromContext(ctx)
    // 這個 log 會自動包含 request_id 和 user_id
    log.Info("processing data")
}
```

### Uber FX 整合

```go
package main

import (
    "github.com/lileon/tabelogov2/pkg/logger"
    "go.uber.org/fx"
    "go.uber.org/zap"
)

func main() {
    app := fx.New(
        logger.Module, // 提供 *zap.Logger

        fx.Invoke(func(log *zap.Logger) {
            log.Info("application started with FX")
        }),
    )

    app.Run()
}
```

### 開發模式 vs 生產模式

```go
// 開發模式（彩色輸出，易讀）
logger.InitDevelopment()

// 生產模式（JSON 格式）
logger.Init("info")
```

---

## 2️⃣ Config 套件 (`pkg/config`)

### 功能特性

✅ 從環境變數載入配置
✅ 支援環境變數前綴（多服務配置）
✅ 完整的配置驗證
✅ 型別安全的配置存取
✅ Uber FX 依賴注入整合

### 配置結構

```go
type Config struct {
    Environment string    // development, staging, production, test
    LogLevel    string    // debug, info, warn, error, fatal
    ServerPort  int       // HTTP server port
    GRPCPort    int       // gRPC server port
    Database    DatabaseConfig
    Redis       RedisConfig
    Kafka       KafkaConfig
    JWT         JWTConfig
}
```

### 基本使用

```go
package main

import (
    "github.com/lileon/tabelogov2/pkg/config"
    "log"
)

func main() {
    // 載入配置
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err)
    }

    // 使用配置
    dsn := cfg.GetDatabaseDSN()
    redisAddr := cfg.GetRedisAddr()
    kafkaBrokers := cfg.GetKafkaBrokers()

    // 環境檢查
    if cfg.IsDevelopment() {
        log.Println("Running in development mode")
    }
}
```

### 環境變數範例

創建 `.env` 檔案：

```bash
# General
ENVIRONMENT=development
LOG_LEVEL=debug

# Server
SERVER_PORT=8080
GRPC_PORT=9090

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=tabelogo_db
DB_USER=postgres
DB_PASSWORD=postgres
DB_SSLMODE=disable
DB_MAX_OPEN_CONNS=100
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=1h

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Kafka
KAFKA_BROKERS=localhost:9092,localhost:9093
KAFKA_GROUP_ID=tabelogo-group

# JWT
JWT_SECRET=your-secret-key
JWT_ACCESS_TOKEN_EXPIRE=15m
JWT_REFRESH_TOKEN_EXPIRE=168h
```

### 使用環境變數前綴（多服務）

```go
// Auth Service
cfg, err := config.LoadWithPrefix("AUTH")
// 會讀取: AUTH_DB_HOST, AUTH_DB_PORT, AUTH_SERVER_PORT 等

// Restaurant Service
cfg, err := config.LoadWithPrefix("RESTAURANT")
// 會讀取: RESTAURANT_DB_HOST, RESTAURANT_DB_PORT 等
```

**環境變數範例**：

```bash
# Auth Service
AUTH_DB_NAME=auth_db
AUTH_DB_PORT=15432
AUTH_SERVER_PORT=8081
AUTH_GRPC_PORT=9091

# Restaurant Service
RESTAURANT_DB_NAME=restaurant_db
RESTAURANT_DB_PORT=5433
RESTAURANT_SERVER_PORT=8082
RESTAURANT_GRPC_PORT=9092
```

### Uber FX 整合

```go
package main

import (
    "github.com/lileon/tabelogov2/pkg/config"
    "go.uber.org/fx"
)

func main() {
    app := fx.New(
        config.Module, // 提供 *config.Config

        fx.Invoke(func(cfg *config.Config) {
            fmt.Printf("Running on port %d\n", cfg.ServerPort)
        }),
    )

    app.Run()
}
```

### 配置驗證

```go
cfg, err := config.Load()
if err != nil {
    // 配置驗證失敗時會返回詳細錯誤訊息
    // 例如: "DB_NAME is required"
    //      "JWT_SECRET must be changed in production environment"
    log.Fatal(err)
}
```

---

## 3️⃣ Errors 套件 (`pkg/errors`)

### 功能特性

✅ 統一的錯誤碼系統
✅ HTTP 狀態碼自動映射
✅ gRPC 錯誤支援
✅ 錯誤包裝與追蹤
✅ 詳細資訊附加

### 錯誤碼定義

```go
const (
    // General errors
    ErrCodeInternal       = "INTERNAL_ERROR"
    ErrCodeInvalidRequest = "INVALID_REQUEST"
    ErrCodeNotFound       = "NOT_FOUND"
    ErrCodeUnauthorized   = "UNAUTHORIZED"
    ErrCodeForbidden      = "FORBIDDEN"
    ErrCodeConflict       = "CONFLICT"

    // Auth errors
    ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"
    ErrCodeTokenExpired       = "TOKEN_EXPIRED"
    ErrCodeTokenInvalid       = "TOKEN_INVALID"

    // Business logic errors
    ErrCodeUserNotFound       = "USER_NOT_FOUND"
    ErrCodeRestaurantNotFound = "RESTAURANT_NOT_FOUND"
    ErrCodeBookingNotFound    = "BOOKING_NOT_FOUND"
    ErrCodeBookingConflict    = "BOOKING_CONFLICT"
)
```

### 基本使用

```go
package main

import (
    "github.com/lileon/tabelogov2/pkg/errors"
)

func getUserByID(id string) (*User, error) {
    // 創建錯誤
    if id == "" {
        return nil, errors.NewInvalidRequestError("user ID is required")
    }

    user, err := db.FindUser(id)
    if err != nil {
        // 包裝錯誤
        return nil, errors.Wrap(err, errors.ErrCodeInternal, "failed to query user")
    }

    if user == nil {
        return nil, errors.New(errors.ErrCodeUserNotFound, "user not found")
    }

    return user, nil
}
```

### 添加詳細資訊

```go
err := errors.NewInvalidRequestError("validation failed").
    WithDetails(map[string]interface{}{
        "field": "email",
        "value": "invalid-email",
        "requirement": "valid email format",
    })

// 錯誤會包含這些詳細資訊，方便 debug
```

### HTTP 處理

```go
func handleError(w http.ResponseWriter, err error) {
    appErr, ok := errors.AsAppError(err)
    if !ok {
        // 不是 AppError，使用預設處理
        appErr = errors.NewInternalError("internal server error")
    }

    w.WriteHeader(appErr.HTTPStatus)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "error": appErr.Code,
        "message": appErr.Message,
        "details": appErr.Details,
    })
}
```

### gRPC 處理

```go
import (
    "github.com/lileon/tabelogov2/pkg/errors"
)

// 服務端：轉換為 gRPC 錯誤
func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    user, err := s.userService.GetByID(req.Id)
    if err != nil {
        // 轉換為 gRPC 錯誤
        if appErr, ok := errors.AsAppError(err); ok {
            return nil, appErr.ToGRPCError()
        }
        return nil, err
    }
    return user, nil
}

// 客戶端：從 gRPC 錯誤轉換
func callGRPC() error {
    resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: "123"})
    if err != nil {
        // 轉換回 AppError
        appErr := errors.FromGRPCError(err)
        return appErr
    }
    return nil
}
```

### 錯誤檢查

```go
err := getUserByID("123")

// 檢查是否為 AppError
if errors.IsAppError(err) {
    appErr, _ := errors.AsAppError(err)
    fmt.Println(appErr.Code)
    fmt.Println(appErr.HTTPStatus)
}

// 使用 errors.As (Go 1.13+)
var appErr *errors.AppError
if errors.As(err, &appErr) {
    fmt.Println(appErr.Code)
}
```

---

## 🎯 完整範例：整合所有套件

```go
package main

import (
    "context"
    "net/http"

    "github.com/lileon/tabelogov2/pkg/config"
    "github.com/lileon/tabelogov2/pkg/errors"
    "github.com/lileon/tabelogov2/pkg/logger"
    "go.uber.org/fx"
    "go.uber.org/zap"
)

func main() {
    app := fx.New(
        // 提供所有共用套件
        config.Module,
        logger.Module,

        // 提供服務
        fx.Provide(NewUserService),
        fx.Provide(NewHTTPServer),

        // 啟動服務
        fx.Invoke(func(*http.Server) {}),
    )

    app.Run()
}

// UserService 範例
type UserService struct {
    cfg    *config.Config
    logger *zap.Logger
}

func NewUserService(cfg *config.Config, logger *zap.Logger) *UserService {
    return &UserService{
        cfg:    cfg,
        logger: logger.With(zap.String("service", "user")),
    }
}

func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
    log := logger.FromContext(ctx)

    if id == "" {
        return nil, errors.NewInvalidRequestError("user ID is required")
    }

    log.Info("fetching user", zap.String("user_id", id))

    // ... 業務邏輯

    return user, nil
}

// HTTP Server 範例
type HTTPServer struct {
    cfg         *config.Config
    logger      *zap.Logger
    userService *UserService
}

func NewHTTPServer(
    lc fx.Lifecycle,
    cfg *config.Config,
    logger *zap.Logger,
    userService *UserService,
) *http.Server {
    mux := http.NewServeMux()

    server := &HTTPServer{
        cfg:         cfg,
        logger:      logger,
        userService: userService,
    }

    mux.HandleFunc("/users/", server.handleGetUser)

    httpServer := &http.Server{
        Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
        Handler: mux,
    }

    lc.Append(fx.Hook{
        OnStart: func(context.Context) error {
            go httpServer.ListenAndServe()
            logger.Info("HTTP server started", zap.Int("port", cfg.ServerPort))
            return nil
        },
        OnStop: func(ctx context.Context) error {
            return httpServer.Shutdown(ctx)
        },
    })

    return httpServer
}

func (s *HTTPServer) handleGetUser(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // 添加 request_id 到 context
    requestID := r.Header.Get("X-Request-ID")
    ctx = logger.WithFields(ctx, zap.String("request_id", requestID))

    id := r.URL.Query().Get("id")
    user, err := s.userService.GetUser(ctx, id)
    if err != nil {
        s.handleError(w, err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func (s *HTTPServer) handleError(w http.ResponseWriter, err error) {
    appErr, ok := errors.AsAppError(err)
    if !ok {
        appErr = errors.NewInternalError("internal server error")
    }

    s.logger.Error("request error",
        zap.String("code", string(appErr.Code)),
        zap.Error(err),
    )

    w.WriteHeader(appErr.HTTPStatus)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "error":   appErr.Code,
        "message": appErr.Message,
        "details": appErr.Details,
    })
}
```

---

## 📝 測試

所有套件都包含完整的單元測試：

```bash
# 執行所有測試
go test ./pkg/...

# 執行特定套件測試
go test ./pkg/logger
go test ./pkg/config
go test ./pkg/errors

# 查看測試覆蓋率
go test ./pkg/... -cover

# 生成覆蓋率報告
go test ./pkg/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## 🎓 最佳實踐

### 1. Logger

- ✅ 使用結構化日誌（zap.Field）
- ✅ 在 context 中傳遞 logger
- ✅ 為每個服務創建專屬的 child logger
- ❌ 避免使用 fmt.Printf 或 log.Println

### 2. Config

- ✅ 在應用啟動時驗證配置
- ✅ 使用環境變數前綴區分不同服務
- ✅ 敏感資訊（密碼、密鑰）不要硬編碼
- ❌ 不要直接使用 os.Getenv，使用 config 套件

### 3. Errors

- ✅ 使用語意化的錯誤碼
- ✅ 包裝底層錯誤以保留 stack trace
- ✅ 為錯誤添加有用的 context
- ❌ 不要忽略錯誤或返回 nil

---

## 📊 效能考量

| 套件 | 效能特性 |
|------|----------|
| Logger | zap 是最快的 Go logger 之一（比 logrus 快 4-10x） |
| Config | 配置只在啟動時載入一次，無執行期開銷 |
| Errors | 輕量級結構，minimal allocation |

---

## 🔗 相關連結

- [Uber Zap Documentation](https://pkg.go.dev/go.uber.org/zap)
- [Uber FX Documentation](https://uber-go.github.io/fx/)
- [12-Factor App Configuration](https://12factor.net/config)

---

## ✅ 更新記錄

- **2025-11-20**: 初始版本完成
  - ✅ logger 套件實作完成
  - ✅ config 套件實作完成
  - ✅ errors 套件實作完成
  - ✅ 添加完整單元測試
  - ✅ 添加 Context 支援（logger）
  - ✅ 添加 gRPC 支援（errors）
