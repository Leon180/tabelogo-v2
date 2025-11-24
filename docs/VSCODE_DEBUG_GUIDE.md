# VSCode 調試指南 - Auth Service

本指南說明如何在 VSCode 中本地調試 Auth Service 並訪問 Swagger UI。

## 問題修復總結

### 修復的問題
1. **Swagger UI 路徑錯誤**: 將 `/app/internal/auth/docs/index.html` (Docker 路徑) 改為 `./internal/auth/docs/index.html` (本地相對路徑)
2. **缺少 VSCode 配置**: 創建了 `launch.json` 和 `tasks.json` 文件
3. **Swagger 生成路徑不一致**: 統一使用 `internal/auth/docs/` 作為輸出目錄

### 修改的文件
- `internal/auth/interfaces/http/module.go` - 修復 Swagger UI 靜態文件路徑
- `scripts/start-auth-service.sh` - 更新 Swagger 生成輸出路徑
- `Makefile` - 統一 Swagger 生成命令參數
- `.vscode/launch.json` - 新建 VSCode 調試配置
- `.vscode/tasks.json` - 新建自動 Swagger 生成任務

## 如何使用

### 方法 1: 使用 VSCode 調試器（推薦）

1. **在 VSCode 中啟動調試**:
   - 按 `F5` 或點擊「Run and Debug」
   - 選擇 "Auth Service" 配置
   - 服務會自動生成 Swagger 文檔（通過 preLaunchTask）並啟動

2. **訪問 Swagger UI**:
   ```
   http://localhost:8081/auth-service/swagger/index.html
   ```

3. **API 端點**:
   - HTTP API: `http://localhost:8081/api/v1`
   - gRPC API: `localhost:9091`
   - Health Check: `http://localhost:8081/health`

### 方法 2: 使用 Makefile

```bash
# 生成 Swagger 文檔
make swagger-auth

# 運行 Auth Service（會自動生成 Swagger 並啟動服務）
make auth-dev
```

### 方法 3: 使用啟動腳本

```bash
./scripts/start-auth-service.sh
```

## 環境要求

### 必須運行的依賴服務

在啟動 Auth Service 之前，確保以下服務正在運行：

```bash
# 啟動 PostgreSQL 和 Redis（使用本地開發專用端口）
docker-compose -f deployments/docker-compose/auth-service.yml up -d postgres-auth redis-auth
```

**注意**: 本地開發環境使用不同的端口以避免衝突：
- PostgreSQL: `15432` (而非標準的 5432)
- Redis: `16379` (而非標準的 6379)

或者單獨啟動：

```bash
# PostgreSQL (Port 15432)
docker run -d \
  --name tabelogo-postgres-auth-dev \
  -p 15432:5432 \
  -e POSTGRES_DB=auth_db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  postgres:15-alpine

# Redis (Port 16379)
docker run -d \
  --name tabelogo-redis-auth-dev \
  -p 16379:6379 \
  redis:7-alpine
```

### 環境變量配置

launch.json 中已配置以下環境變量：

```json
{
  "ENVIRONMENT": "development",
  "SERVER_PORT": "8081",
  "GRPC_PORT": "9091",
  "DB_HOST": "localhost",
  "DB_PORT": "15432",
  "DB_NAME": "auth_db",
  "DB_USER": "postgres",
  "DB_PASSWORD": "postgres",
  "REDIS_HOST": "localhost",
  "REDIS_PORT": "16379",
  "REDIS_PASSWORD": "",
  "REDIS_DB": "0",
  "JWT_SECRET": "your-super-secret-jwt-key-change-in-production",
  "JWT_ACCESS_TOKEN_EXPIRE": "15m",
  "JWT_REFRESH_TOKEN_EXPIRE": "168h"
}
```

## Swagger 文檔

### Swagger UI 端點

- **Swagger UI**: `http://localhost:8081/auth-service/swagger/index.html`
- **Swagger JSON**: `http://localhost:8081/auth-service/swagger/doc.json`
- **Quick Access**: `http://localhost:8081/swagger` (redirects to full path)

### 重新生成 Swagger 文檔

如果修改了 API 代碼，需要重新生成文檔：

```bash
# 方法 1: 使用 Makefile
make swagger-auth

# 方法 2: 直接使用 swag 命令
swag init --generalInfo cmd/auth-service/main.go --output internal/auth/docs --parseDependency --parseInternal
```

## 調試技巧

### 設置斷點
1. 在代碼中點擊行號左側設置斷點
2. 啟動調試模式（F5）
3. 發送 API 請求，程序會在斷點處暫停

### 查看日誌
- VSCode Debug Console 會顯示應用程序的輸出
- 使用結構化日誌記錄 (zap logger)

### 常見問題

#### 1. 無法訪問 Swagger UI (404)
**原因**: index.html 文件路徑錯誤或不存在

**解決方法**:
```bash
# 確認文件存在
ls -la internal/auth/docs/index.html

# 重新生成 Swagger 文檔
make swagger-auth
```

#### 2. 數據庫連接失敗
**原因**: PostgreSQL 未啟動或連接參數錯誤

**解決方法**:
```bash
# 檢查 PostgreSQL 狀態
docker ps | grep postgres-auth

# 啟動 PostgreSQL
docker start tabelogo-postgres-auth-dev

# 或使用 docker-compose
docker-compose -f deployments/docker-compose/auth-service.yml up -d postgres-auth

# 測試連接（注意端口是 15432）
docker exec -it tabelogo-postgres-auth-dev psql -U postgres -d auth_db
```

#### 3. Redis 連接失敗
**原因**: Redis 未啟動

**解決方法**:
```bash
# 檢查 Redis 狀態
docker ps | grep redis-auth

# 啟動 Redis
docker start tabelogo-redis-auth-dev

# 或使用 docker-compose
docker-compose -f deployments/docker-compose/auth-service.yml up -d redis-auth

# 測試連接（注意端口是 16379）
docker exec -it tabelogo-redis-auth-dev redis-cli
```

#### 4. Swagger 文檔內容過時
**原因**: 修改了代碼但未重新生成文檔

**解決方法**:
```bash
make swagger-auth
```

## 其他微服務

launch.json 中也配置了其他微服務的調試配置：

- Restaurant Service (Port 8082, gRPC 9092)
- Booking Service (Port 8083, gRPC 9093)
- Mail Service (Port 8084, gRPC 9094)
- Spider Service (Port 8085, gRPC 9095)
- Map Service (Port 8086, gRPC 9096)
- API Gateway (Port 8080)

可以在 VSCode 的「Run and Debug」面板中選擇對應的配置來調試。

## 技術細節

### Swagger 文檔生成位置
- **生成目錄**: `internal/auth/docs/`
- **包含文件**:
  - `docs.go` - Swagger 元數據和文檔模板
  - `swagger.json` - OpenAPI JSON 規範
  - `swagger.yaml` - OpenAPI YAML 規範
  - `index.html` - Swagger UI 界面

### 路徑配置
- **開發環境**: 使用相對路徑 `./internal/auth/docs/index.html`
- **Docker 環境**: 使用絕對路徑 `/app/internal/auth/docs/index.html` (通過 volume mount)

### Gin Router 配置
```go
// Serve swagger.json
router.GET("/swagger/doc.json", func(c *gin.Context) {
    c.String(200, docs.SwaggerInfo.ReadDoc())
})

// Serve Swagger UI
router.StaticFile("/swagger/index.html", "./internal/auth/docs/index.html")

// Redirect /swagger to /swagger/index.html
router.GET("/swagger", func(c *gin.Context) {
    c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
})
```

## 總結

現在您可以：
1. ✅ 在 VSCode 中使用 F5 快速啟動調試
2. ✅ 自動生成 Swagger 文檔
3. ✅ 訪問 Swagger UI: http://localhost:8081/auth-service/swagger/index.html
4. ✅ 設置斷點進行代碼調試
5. ✅ 輕鬆切換調試不同的微服務

祝您調試愉快！🎉
