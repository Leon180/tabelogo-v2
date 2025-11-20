# 共用套件 (pkg) 設計說明

本專案的共用套件設計參考了多個知名 Go 開源專案的最佳實踐。

## 📚 參考專案

1. **Kubernetes** - 大型微服務架構典範
2. **Docker/Moby** - 容器化平台
3. **Grafana** - 可觀測性平台
4. **Gin** - 高效能 Web 框架
5. **Go-kit** - 微服務工具包

---

## 📦 套件說明

### 1. pkg/logger - 日誌套件

**設計理念：**
- 參考 **Kubernetes** 和 **Grafana** 的日誌設計
- 使用結構化日誌 (Structured Logging)，方便後續分析與查詢
- 統一的日誌介面，所有微服務使用相同的日誌格式

**為什麼選擇 Zap？**
```
Uber 開源的高效能日誌庫，被廣泛用於生產環境：
✅ 極高的效能（零記憶體分配）
✅ 結構化日誌（JSON 格式）
✅ 豐富的欄位類型
✅ 支援日誌層級控制
```

**設計決策：**

1. **全域 Logger vs 依賴注入**
   - 我們使用全域 Logger + Getter 模式
   - 原因：方便使用，減少傳遞參數的複雜度
   - Kubernetes 也採用類似設計 (klog)

2. **開發模式 vs 生產模式**
   ```go
   // 開發模式：人類可讀的 console 輸出
   logger.InitDevelopment()

   // 生產模式：結構化 JSON 輸出，方便 ELK/Grafana Loki 收集
   logger.Init("info")
   ```

3. **Caller Skip**
   ```go
   log, err = config.Build(zap.AddCallerSkip(1))
   ```
   - 為什麼 Skip 1？因為我們包裝了 logger，需要跳過包裝層顯示真實呼叫位置
   - 參考：Grafana Loki 的 logger 包裝

**使用範例：**
```go
import "github.com/Leon180/tabelogo-v2/pkg/logger"

// 初始化
logger.Init("info")
defer logger.Sync()

// 使用
logger.Info("User logged in",
    zap.String("user_id", userID),
    zap.Duration("login_time", duration),
)
```

---

### 2. pkg/config - 配置管理

**設計理念：**
- 參考 **12-Factor App** 方法論
- 環境變數優先 (Environment Variables)
- 參考 **Docker** 的配置管理方式

**為什麼不使用 Viper？**
```
雖然 Viper 功能強大，但我們選擇簡單的環境變數方式：
✅ 符合 12-Factor App 原則
✅ 容器化友善（Docker/Kubernetes 原生支援）
✅ 減少依賴，程式碼更簡潔
✅ 避免配置檔案的複雜性

參考：Kubernetes 本身也是直接讀取環境變數
```

**設計決策：**

1. **型別安全**
   ```go
   type Config struct {
       Database DatabaseConfig  // 結構化配置
       Redis    RedisConfig
       JWT      JWTConfig
   }
   ```
   - 使用強型別結構，編譯期檢查錯誤
   - 參考：Go-kit 的配置設計

2. **預設值 + 驗證**
   ```go
   func (c *Config) Validate() error {
       if c.Database.Name == "" {
           return fmt.Errorf("DB_NAME is required")
       }
       return nil
   }
   ```
   - 提供合理的預設值
   - 啟動時驗證必要欄位
   - Fail-fast 原則：盡早發現配置錯誤

3. **DSN 生成器**
   ```go
   func (c *Config) GetDatabaseDSN() string
   ```
   - 封裝連線字串生成邏輯
   - 避免在各服務重複撰寫
   - 參考：GORM 社群最佳實踐

**使用範例：**
```go
import "github.com/Leon180/tabelogo-v2/pkg/config"

cfg, err := config.Load()
if err != nil {
    log.Fatal("Failed to load config", zap.Error(err))
}

// 取得 DB 連線字串
dsn := cfg.GetDatabaseDSN()
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
```

---

### 3. pkg/errors - 錯誤處理

**設計理念：**
- 參考 **Google Cloud Go SDK** 的錯誤設計
- 參考 **Twirp RPC** 的錯誤碼設計
- 統一的錯誤格式，方便前端處理

**為什麼需要自訂錯誤？**
```
Go 標準錯誤太簡單，無法滿足微服務需求：
✅ 需要錯誤碼（方便前端國際化）
✅ 需要 HTTP Status Code（RESTful API）
✅ 需要額外資訊（Details）
✅ 需要錯誤鏈（Unwrap）

參考：Kubernetes API 也有類似的錯誤設計
```

**設計決策：**

1. **錯誤碼 + HTTP Status 映射**
   ```go
   type AppError struct {
       Code       ErrorCode  // 業務錯誤碼
       HTTPStatus int        // HTTP 狀態碼
       Message    string     // 錯誤訊息
       Details    map[string]interface{}  // 額外資訊
       Err        error      // 原始錯誤（支援 Unwrap）
   }
   ```
   - 參考：gRPC Status Codes
   - 好處：前端可以根據 Code 做不同處理

2. **錯誤包裝 (Error Wrapping)**
   ```go
   errors.Wrap(err, ErrCodeInternal, "Failed to query database")
   ```
   - 支援 Go 1.13+ 的錯誤鏈
   - 可以用 `errors.Is()` 和 `errors.As()` 檢查
   - 參考：Go 官方的 errors 包設計

3. **預定義錯誤構造器**
   ```go
   errors.NewNotFoundError("User not found")
   errors.NewUnauthorizedError("Invalid token")
   ```
   - 減少重複程式碼
   - 統一錯誤訊息格式
   - 參考：Google Cloud Go SDK

**使用範例：**
```go
import "github.com/Leon180/tabelogo-v2/pkg/errors"

// 建立錯誤
if user == nil {
    return errors.NewNotFoundError("User not found").
        WithDetails(map[string]interface{}{
            "user_id": userID,
        })
}

// 包裝錯誤
if err := db.Save(&user).Error; err != nil {
    return errors.Wrap(err, errors.ErrCodeInternal, "Failed to save user")
}

// 在 HTTP Handler 中處理
if err != nil {
    appErr, ok := errors.AsAppError(err)
    if ok {
        c.JSON(appErr.HTTPStatus, gin.H{
            "code": appErr.Code,
            "message": appErr.Message,
        })
    }
}
```

---

### 4. pkg/middleware - HTTP 中介層

**設計理念：**
- 參考 **Gin** 官方中介層設計
- 參考 **Echo** 和 **Chi** 的中介層實踐
- 關注點分離 (Separation of Concerns)

**為什麼需要這些 Middleware？**

#### 4.1 Logger Middleware
```go
middleware.Logger(logger)
```
**設計原因：**
- 統一記錄所有 HTTP 請求
- 包含：狀態碼、延遲、IP、User-Agent
- 根據狀態碼決定日誌層級 (Info/Warn/Error)
- 參考：Kubernetes API Server 的請求日誌

#### 4.2 Recovery Middleware
```go
middleware.Recovery(logger)
```
**設計原因：**
- 捕捉 panic，避免服務崩潰
- 記錄 panic 的 stack trace
- 回傳友善的錯誤訊息給客戶端
- 參考：Gin 官方的 Recovery 中介層

**為什麼要自己實作而不用 Gin 內建？**
```
✅ 整合我們的 logger（zap）
✅ 整合我們的錯誤處理（AppError）
✅ 統一的錯誤回應格式
```

#### 4.3 CORS Middleware
```go
middleware.CORS()
```
**設計原因：**
- 處理跨域請求（前後端分離必備）
- 處理 OPTIONS 預檢請求
- 參考：gin-contrib/cors

#### 4.4 Error Handler Middleware
```go
middleware.ErrorHandler()
```
**設計原因：**
- 統一錯誤回應格式
- 自動將 AppError 轉換為 JSON 回應
- 避免在每個 Handler 重複寫錯誤處理
- 參考：Go-kit 的錯誤處理設計

**使用範例：**
```go
import (
    "github.com/gin-gonic/gin"
    "github.com/Leon180/tabelogo-v2/pkg/middleware"
    "github.com/Leon180/tabelogo-v2/pkg/logger"
)

func main() {
    r := gin.New()

    // 按順序註冊 middleware
    r.Use(middleware.Recovery(logger.GetLogger()))
    r.Use(middleware.Logger(logger.GetLogger()))
    r.Use(middleware.CORS())
    r.Use(middleware.ErrorHandler())

    // 定義路由
    r.GET("/users/:id", getUserHandler)
}
```

---

## 🎯 整體設計原則

### 1. 依賴方向
```
服務層 (cmd/auth-service)
    ↓ 依賴
共用層 (pkg/logger, config, errors)
    ↓ 依賴
第三方庫 (zap, gin, gorm)
```
- pkg 不依賴任何 internal 程式碼
- pkg 可以被所有服務使用
- 參考：Go 標準庫的設計

### 2. 介面隔離
- 每個套件職責單一
- logger 只負責日誌
- config 只負責配置
- errors 只負責錯誤
- 參考：SOLID 原則

### 3. 零依賴原則
- pkg 之間盡量減少依賴
- logger 不依賴 config
- errors 不依賴 logger
- 好處：減少耦合，方便測試

---

## 📖 延伸閱讀

1. [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
2. [Kubernetes Logging](https://kubernetes.io/docs/concepts/cluster-administration/logging/)
3. [12-Factor App](https://12factor.net/)
4. [Go Error Handling](https://go.dev/blog/go1.13-errors)
5. [Effective Go](https://go.dev/doc/effective_go)

---

## 🔄 下一步

當您實作服務時，可以這樣使用這些套件：

```go
package main

import (
    "github.com/Leon180/tabelogo-v2/pkg/config"
    "github.com/Leon180/tabelogo-v2/pkg/logger"
    "github.com/Leon180/tabelogo-v2/pkg/errors"
)

func main() {
    // 1. 載入配置
    cfg, err := config.Load()
    if err != nil {
        panic(err)
    }

    // 2. 初始化 logger
    if cfg.IsDevelopment() {
        logger.InitDevelopment()
    } else {
        logger.Init(cfg.LogLevel)
    }
    defer logger.Sync()

    // 3. 使用
    logger.Info("Service starting",
        zap.String("env", cfg.Environment),
    )

    // 4. 錯誤處理範例
    if err := doSomething(); err != nil {
        logger.Error("Failed", zap.Error(err))
        return errors.Wrap(err, errors.ErrCodeInternal, "Operation failed")
    }
}
```

---

## 📦 Package: migrations

資料庫 migration 版本控制系統。

### 快速使用

```go
import "github.com/Leon180/tabelogo-v2/pkg/migrations"

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

### 詳細文檔

- [完整使用手冊](migrations/README.md)
- [Import 指南](migrations/IMPORT_GUIDE.md)

---

## Module 管理說明

### 為什麼統一在 pkg/ 層級管理 go.mod?

所有 `pkg/` 下的子目錄都屬於同一個 module: `github.com/Leon180/tabelogo-v2/pkg`

**優點**:
1. ✅ 簡化依賴管理 - 只需維護一個 go.mod
2. ✅ 避免循環依賴 - pkg 內的套件可以互相引用
3. ✅ 版本統一 - 所有套件使用相同版本的依賴
4. ✅ 符合 Go 慣例 - 官方推薦做法

### Import 路徑

```go
// ✅ 正確
import "github.com/Leon180/tabelogo-v2/pkg/migrations"
import "github.com/Leon180/tabelogo-v2/pkg/logger"

// ❌ 錯誤
import "pkg/migrations"
import "../logger"
```

### 依賴管理

所有依賴在 `pkg/go.mod` 中統一管理:

```bash
# 添加新依賴
cd pkg
go get github.com/new/package@version

# 清理依賴
go mod tidy
```
