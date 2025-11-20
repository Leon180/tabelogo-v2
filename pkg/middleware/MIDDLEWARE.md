# Middleware Package

完整的 HTTP 中間件實現，提供認證、日誌、錯誤處理、速率限制等功能。

## 📦 包含的中間件

### 1. 認證中間件 (Authentication)

JWT 基礎的認證與授權中間件。

#### 功能特性

- ✅ JWT Token 驗證
- ✅ Bearer Token 解析
- ✅ 用戶角色檢查 (RBAC)
- ✅ 跳過路徑配置
- ✅ Context 中的用戶資訊

#### 使用範例

```go
package main

import (
    "github.com/Leon180/tabelogo-v2/pkg/middleware"
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    // 配置認證中間件
    authConfig := middleware.AuthConfig{
        JWTSecret: "your-secret-key",
        SkipPaths: []string{"/health", "/api/v1/auth/login"},
    }

    // 應用到所有路由
    router.Use(middleware.Auth(authConfig))

    // 需要特定角色的路由
    admin := router.Group("/admin")
    admin.Use(middleware.RequireRole("admin", "moderator"))
    {
        admin.GET("/users", listUsers)
        admin.DELETE("/users/:id", deleteUser)
    }

    // 獲取當前用戶資訊
    router.GET("/profile", func(c *gin.Context) {
        userID, exists := middleware.GetUserID(c)
        if !exists {
            c.JSON(401, gin.H{"error": "Unauthorized"})
            return
        }

        role, _ := middleware.GetUserRole(c)
        c.JSON(200, gin.H{
            "user_id": userID,
            "role":    role,
        })
    })
}
```

#### JWT Claims 結構

```go
type JWTClaims struct {
    UserID string `json:"user_id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}
```

---

### 2. 速率限制中間件 (Rate Limiting)

提供基於 Redis 的分散式速率限制和記憶體內速率限制。

#### 功能特性

- ✅ Redis 分散式速率限制（生產環境）
- ✅ 記憶體內速率限制（開發環境）
- ✅ 滑動視窗演算法
- ✅ 按 IP 或用戶 ID 限流
- ✅ 自動設定速率限制標頭
- ✅ 跳過路徑配置

#### 使用範例

**生產環境 - Redis 速率限制**

```go
import (
    "time"
    "github.com/Leon180/tabelogo-v2/pkg/middleware"
    "github.com/redis/go-redis/v9"
)

func main() {
    // 初始化 Redis 客戶端
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        DB:   4, // 速率限制專用 DB
    })

    router := gin.Default()

    // 配置速率限制
    rateLimitConfig := middleware.RateLimitConfig{
        RedisClient: redisClient,
        Limit:       100,                    // 每個視窗 100 次請求
        Window:      time.Minute,            // 1 分鐘視窗
        KeyPrefix:   "ratelimit:api",        // Redis key 前綴
        SkipPaths:   []string{"/health"},    // 跳過健康檢查
    }

    router.Use(middleware.RateLimit(rateLimitConfig))
}
```

**開發環境 - 記憶體內速率限制**

```go
func main() {
    router := gin.Default()

    // 簡單的記憶體內速率限制
    router.Use(middleware.InMemoryRateLimit(
        10,            // 限制次數
        time.Minute,   // 時間視窗
    ))
}
```

#### 回應標頭

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1700000000
```

---

### 3. 請求 ID 中間件 (Request ID)

為每個請求添加唯一的追蹤 ID。

#### 功能特性

- ✅ 自動生成 UUID
- ✅ 支援現有 Request ID
- ✅ Context 中的 Request ID
- ✅ 自動設定回應標頭

#### 使用範例

```go
import "github.com/Leon180/tabelogo-v2/pkg/middleware"

func main() {
    router := gin.Default()

    // 應用 Request ID 中間件
    router.Use(middleware.RequestID())

    router.GET("/api/test", func(c *gin.Context) {
        // 獲取 Request ID
        requestID, exists := middleware.GetRequestID(c)
        if exists {
            // 使用 Request ID 進行日誌追蹤
            logger.Info("Processing request",
                zap.String("request_id", requestID))
        }

        c.JSON(200, gin.H{
            "request_id": requestID,
        })
    })
}
```

#### 回應標頭

```
X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
```

---

### 4. 日誌中間件 (Logger)

結構化 HTTP 請求日誌。

#### 功能特性

- ✅ Uber Zap 結構化日誌
- ✅ 請求/回應詳細資訊
- ✅ 延遲時間追蹤
- ✅ 錯誤日誌
- ✅ 不同狀態碼的日誌等級

#### 使用範例

```go
import (
    "github.com/Leon180/tabelogo-v2/pkg/logger"
    "github.com/Leon180/tabelogo-v2/pkg/middleware"
)

func main() {
    // 初始化 logger
    logger.Init("production", "json")
    log := logger.GetLogger()

    router := gin.Default()

    // 應用日誌中間件
    router.Use(middleware.Logger(log))

    router.GET("/api/test", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "success"})
    })
}
```

#### 日誌欄位

```json
{
  "level": "info",
  "ts": "2025-11-20T10:00:00.000Z",
  "msg": "Request completed",
  "status": 200,
  "method": "GET",
  "path": "/api/test",
  "query": "param1=value1",
  "ip": "192.168.1.1",
  "latency": "15ms",
  "user_agent": "Mozilla/5.0..."
}
```

---

### 5. 錯誤處理中間件 (Error Handler)

統一的錯誤處理與回應格式化。

#### 功能特性

- ✅ AppError 類型識別
- ✅ HTTP 狀態碼自動映射
- ✅ 統一錯誤回應格式
- ✅ 詳細錯誤資訊

#### 使用範例

```go
import (
    "github.com/Leon180/tabelogo-v2/pkg/errors"
    "github.com/Leon180/tabelogo-v2/pkg/middleware"
)

func main() {
    router := gin.Default()

    // 應用錯誤處理中間件
    router.Use(middleware.ErrorHandler())

    router.GET("/api/user/:id", func(c *gin.Context) {
        user, err := getUserByID(c.Param("id"))
        if err != nil {
            // 使用 AppError
            appErr := errors.New(errors.ErrCodeNotFound, "User not found").
                WithDetails(map[string]interface{}{
                    "user_id": c.Param("id"),
                })
            c.Error(appErr)
            return
        }

        c.JSON(200, user)
    })
}
```

#### 錯誤回應格式

```json
{
  "code": "NOT_FOUND",
  "message": "User not found",
  "details": {
    "user_id": "12345"
  }
}
```

---

### 6. 恢復中間件 (Recovery)

捕獲 panic 並記錄錯誤。

#### 功能特性

- ✅ Panic 捕獲
- ✅ 錯誤日誌記錄
- ✅ 優雅的錯誤回應
- ✅ 防止服務崩潰

#### 使用範例

```go
import (
    "github.com/Leon180/tabelogo-v2/pkg/logger"
    "github.com/Leon180/tabelogo-v2/pkg/middleware"
)

func main() {
    logger.Init("production", "json")
    log := logger.GetLogger()

    router := gin.Default()

    // 應用恢復中間件（應該最先註冊）
    router.Use(middleware.Recovery(log))

    router.GET("/api/panic", func(c *gin.Context) {
        panic("Something went wrong!")
    })
}
```

---

### 7. CORS 中間件 (CORS)

跨來源資源共享配置。

#### 功能特性

- ✅ 允許所有來源（可配置）
- ✅ OPTIONS 請求處理
- ✅ 憑證支援
- ✅ 自定義標頭和方法

#### 使用範例

```go
import "github.com/Leon180/tabelogo-v2/pkg/middleware"

func main() {
    router := gin.Default()

    // 應用 CORS 中間件
    router.Use(middleware.CORS())

    router.GET("/api/test", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "CORS enabled"})
    })
}
```

#### CORS 標頭

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
Access-Control-Allow-Headers: Content-Type, Authorization, ...
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, PATCH, OPTIONS
```

---

## 🔗 中間件鏈範例

完整的中間件配置範例：

```go
package main

import (
    "time"

    "github.com/Leon180/tabelogo-v2/pkg/config"
    "github.com/Leon180/tabelogo-v2/pkg/logger"
    "github.com/Leon180/tabelogo-v2/pkg/middleware"
    "github.com/gin-gonic/gin"
    "github.com/redis/go-redis/v9"
)

func main() {
    // 載入配置
    cfg, _ := config.Load()

    // 初始化 logger
    logger.Init(cfg.Environment, "json")
    log := logger.GetLogger()

    // 初始化 Redis
    redisClient := redis.NewClient(&redis.Options{
        Addr: cfg.Redis.Host + ":" + cfg.Redis.Port,
        DB:   4,
    })

    // 建立 router
    router := gin.New() // 使用 New() 而不是 Default()

    // 1. Recovery（最先，捕獲所有 panic）
    router.Use(middleware.Recovery(log))

    // 2. Request ID（為請求添加追蹤 ID）
    router.Use(middleware.RequestID())

    // 3. Logger（記錄所有請求）
    router.Use(middleware.Logger(log))

    // 4. CORS（處理跨來源請求）
    router.Use(middleware.CORS())

    // 5. Rate Limiting（限制請求速率）
    rateLimitConfig := middleware.RateLimitConfig{
        RedisClient: redisClient,
        Limit:       100,
        Window:      time.Minute,
        KeyPrefix:   "ratelimit:api",
        SkipPaths:   []string{"/health", "/metrics"},
    }
    router.Use(middleware.RateLimit(rateLimitConfig))

    // 6. Auth（保護需要認證的路由）
    authConfig := middleware.AuthConfig{
        JWTSecret: cfg.JWT.Secret,
        SkipPaths: []string{
            "/health",
            "/metrics",
            "/api/v1/auth/login",
            "/api/v1/auth/register",
        },
    }
    router.Use(middleware.Auth(authConfig))

    // 7. Error Handler（統一錯誤處理）
    router.Use(middleware.ErrorHandler())

    // 註冊路由
    router.GET("/health", healthCheck)

    // API 路由
    api := router.Group("/api/v1")
    {
        api.GET("/restaurants", listRestaurants)
        api.GET("/restaurants/:id", getRestaurant)

        // 需要管理員權限的路由
        admin := api.Group("/admin")
        admin.Use(middleware.RequireRole("admin"))
        {
            admin.POST("/restaurants", createRestaurant)
            admin.DELETE("/restaurants/:id", deleteRestaurant)
        }
    }

    router.Run(":8080")
}
```

---

## 🧪 測試

所有中間件都包含完整的單元測試。

### 執行測試

```bash
# 測試所有中間件
go test ./pkg/middleware -v

# 測試特定中間件
go test -run TestAuth ./pkg/middleware -v
go test -run TestRateLimit ./pkg/middleware -v

# 檢查測試覆蓋率
go test ./pkg/middleware -cover
```

### 測試範例

```go
func TestAuth(t *testing.T) {
    gin.SetMode(gin.TestMode)

    testSecret := "test-secret"
    config := middleware.AuthConfig{
        JWTSecret: testSecret,
    }

    router := gin.New()
    router.Use(middleware.Auth(config))
    router.GET("/test", func(c *gin.Context) {
        userID, _ := middleware.GetUserID(c)
        c.JSON(200, gin.H{"user_id": userID})
    })

    // 建立測試 token
    token := createTestToken(testSecret, "user123", "user", time.Hour)

    // 發送請求
    req := httptest.NewRequest(http.MethodGet, "/test", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
}
```

---

## 📊 效能考量

### 中間件順序重要性

中間件的執行順序很重要，建議順序：

1. **Recovery** - 最先，捕獲所有 panic
2. **RequestID** - 早期添加，用於後續日誌追蹤
3. **Logger** - 記錄請求詳情
4. **CORS** - 處理 OPTIONS 請求，避免不必要的處理
5. **RateLimit** - 早期限流，減少無效請求
6. **Auth** - 認證檢查
7. **ErrorHandler** - 最後，統一處理錯誤

### 效能優化建議

#### 1. 速率限制

- **生產環境**: 使用 Redis 分散式速率限制
- **開發環境**: 使用記憶體內速率限制
- **考慮**: 為不同端點設定不同限制

```go
// 不同端點不同限制
publicAPI := router.Group("/api/v1")
publicAPI.Use(middleware.InMemoryRateLimit(10, time.Minute))

authenticatedAPI := router.Group("/api/v1")
authenticatedAPI.Use(middleware.Auth(authConfig))
authenticatedAPI.Use(middleware.InMemoryRateLimit(100, time.Minute))
```

#### 2. JWT 驗證

- 緩存公鑰/密鑰
- 使用適當的過期時間
- 考慮使用 Redis 進行 token 黑名單

#### 3. 日誌

- 生產環境使用 JSON 格式
- 避免在 hot path 中進行高成本操作
- 使用適當的日誌等級

---

## 🔒 安全性最佳實踐

### 1. JWT 密鑰管理

```go
// ❌ 錯誤 - 硬編碼密鑰
authConfig := middleware.AuthConfig{
    JWTSecret: "my-secret-key",
}

// ✅ 正確 - 從環境變數載入
authConfig := middleware.AuthConfig{
    JWTSecret: os.Getenv("JWT_SECRET"),
}
```

### 2. CORS 配置

```go
// ❌ 生產環境不建議 - 允許所有來源
router.Use(middleware.CORS())

// ✅ 建議 - 限制特定來源（需要自定義 CORS 中間件）
corsConfig := cors.Config{
    AllowOrigins: []string{"https://example.com"},
    AllowMethods: []string{"GET", "POST"},
}
router.Use(cors.New(corsConfig))
```

### 3. 速率限制

- 為 API 端點設定適當的限制
- 監控速率限制觸發
- 考慮分層限制（IP、用戶、端點）

---

## 🔄 與其他套件整合

### 與 Logger 套件整合

```go
import (
    "github.com/Leon180/tabelogo-v2/pkg/logger"
    "github.com/Leon180/tabelogo-v2/pkg/middleware"
)

// 使用 context logger
router.Use(middleware.RequestID())
router.Use(func(c *gin.Context) {
    requestID, _ := middleware.GetRequestID(c)
    ctx := logger.WithFields(c.Request.Context(),
        zap.String("request_id", requestID))
    c.Request = c.Request.WithContext(ctx)
    c.Next()
})
```

### 與 Errors 套件整合

```go
import (
    "github.com/Leon180/tabelogo-v2/pkg/errors"
    "github.com/Leon180/tabelogo-v2/pkg/middleware"
)

router.Use(middleware.ErrorHandler())

router.GET("/api/user/:id", func(c *gin.Context) {
    user, err := userService.GetByID(c.Param("id"))
    if err != nil {
        // ErrorHandler 會自動處理 AppError
        c.Error(errors.Wrap(err, errors.ErrCodeNotFound, "User not found"))
        return
    }
    c.JSON(200, user)
})
```

---

## 📝 總結

本中間件套件提供了構建安全、可觀測、高效能 HTTP API 所需的所有基礎組件：

- ✅ **完整功能**: 7 個生產級中間件
- ✅ **安全性**: JWT 認證、速率限制、CORS
- ✅ **可觀測性**: 日誌、Request ID、錯誤追蹤
- ✅ **測試**: 完整的單元測試覆蓋
- ✅ **文檔**: 詳細的使用範例和最佳實踐
- ✅ **效能**: 優化的執行順序和配置建議

這些中間件已準備好用於生產環境，並與專案的其他共用套件（logger, config, errors）完美整合。
