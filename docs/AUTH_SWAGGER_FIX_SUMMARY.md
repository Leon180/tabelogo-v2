# Auth Service Swagger UI 完整修復總結

## 📋 概述

本次修復解決了 Auth Service 在本地開發（VSCode）和 Docker 環境中的 Swagger UI 訪問問題。

**修復日期**: 2025-11-25
**影響範圍**: Auth Service Swagger UI 文檔訪問
**環境**: 本地開發 + Docker 部署

## 🐛 問題描述

### 1. 本地開發環境問題

**症狀**:
- VSCode F5 啟動後，訪問 Swagger UI 出現 `ERR_TOO_MANY_REDIRECTS` 錯誤
- 瀏覽器顯示 "redirected you too many times"

**根本原因**:
1. **錯誤路徑配置**: module.go 使用 Docker 絕對路徑 `/app/internal/auth/docs/index.html`
2. **環境變數錯誤**: launch.json 使用錯誤的變數名稱（`REDIS_ADDR` 而非 `REDIS_HOST`+`REDIS_PORT`）
3. **重定向循環**:
   - `/swagger` → 302 重定向到 `/auth-service/swagger/index.html`
   - Gin 的 `c.File()` 使用 `http.ServeFile` 自動返回 301 到 `./`
   - 形成無限重定向循環

### 2. Docker 環境問題

**症狀**:
- Docker 容器啟動後訪問 Swagger UI 返回 404 Not Found
- 日誌顯示: `GET '/auth-service/swagger/index.html' [GIN] | 404`

**根本原因**:
```dockerfile
# ❌ Dockerfile 錯誤配置
COPY --from=builder /app/cmd/auth-service/docs ./docs
# 複製到了 ./docs，但代碼期望在 ./internal/auth/docs/
```

## ✅ 解決方案

### 1. 路徑配置統一化

**Swagger 生成路徑統一為**: `internal/auth/docs/`

**修改文件**:
- Makefile: 更新 `swagger-auth` 目標
- Dockerfile: 更新 Swagger 生成和複製路徑
- .vscode/tasks.json: 新增 `swag-init-auth` 任務

### 2. 修復重定向循環

**核心修改**: [internal/auth/interfaces/http/module.go](../internal/auth/interfaces/http/module.go)

```go
// ✅ 解決方案 1: 禁用 Gin 自動重定向
router.RedirectTrailingSlash = false
router.RedirectFixedPath = false

// ✅ 解決方案 2: 直接讀取文件而非使用 c.File()
router.GET("/auth-service/swagger/index.html", func(c *gin.Context) {
    absPath, err := filepath.Abs("./internal/auth/docs/index.html")
    if err != nil {
        logger.Error("Failed to resolve Swagger UI path", zap.Error(err))
        c.String(http.StatusInternalServerError, "Internal server error")
        return
    }

    content, err := os.ReadFile(absPath)
    if err != nil {
        logger.Error("Failed to read Swagger UI file", zap.Error(err))
        c.String(http.StatusNotFound, "Swagger UI not found")
        return
    }

    c.Data(http.StatusOK, "text/html; charset=utf-8", content)
})
```

**為什麼這樣做**:
- `http.ServeFile` 會對嵌套路徑自動返回 301 重定向
- 直接用 `os.ReadFile` + `c.Data` 避免了 HTTP 文件服務的自動行為
- 使用 302 (Found) 而非 301 (Permanent) 避免瀏覽器緩存

### 3. 服務特定路徑

**URL 架構變更**:
```
舊路徑: /swagger/*
新路徑: /auth-service/swagger/*
```

**好處**:
- 支援多服務架構（未來可能有 restaurant-service, booking-service 等）
- 避免路徑衝突
- 更清晰的服務邊界

### 4. Docker 路徑修復

**Dockerfile 修改**: [cmd/auth-service/Dockerfile:55](../cmd/auth-service/Dockerfile#L55)

```dockerfile
# ❌ 修復前
COPY --from=builder /app/cmd/auth-service/docs ./docs

# ✅ 修復後
COPY --from=builder /app/internal/auth/docs ./internal/auth/docs
```

**容器內文件結構**:
```
/app/
├── auth-service              # 二進制文件
└── internal/
    └── auth/
        └── docs/
            ├── docs.go       # Swagger 元數據
            ├── index.html    ✅ 更新為使用 /auth-service/swagger/doc.json
            ├── swagger.json  # OpenAPI JSON
            └── swagger.yaml  # OpenAPI YAML
```

### 5. 環境變數修復

**VSCode Launch Configuration**: [.vscode/launch.json](../.vscode/launch.json)

```json
{
  "REDIS_HOST": "localhost",      // ✅ 正確（原為 REDIS_ADDR）
  "REDIS_PORT": "16379",           // ✅ 正確
  "DB_PORT": "15432",              // ✅ 正確（本地開發端口）
  "JWT_ACCESS_TOKEN_EXPIRE": "15m",  // ✅ 正確（原為 DURATION）
  "JWT_REFRESH_TOKEN_EXPIRE": "168h" // ✅ 正確
}
```

## 📁 文件變更清單

### 核心代碼修改

| 文件 | 變更類型 | 說明 |
|------|---------|------|
| `internal/auth/interfaces/http/module.go` | 修改 | 修復重定向循環，更新 Swagger 路徑 |
| `internal/auth/docs/index.html` | 修改 | 更新 Swagger JSON URL |
| `cmd/auth-service/Dockerfile` | 修改 | 修復 Swagger 文檔複製路徑 |
| `cmd/auth-service/main.go` | 修改 | 新增 Swagger 註解 |

### VSCode 配置

| 文件 | 變更類型 | 說明 |
|------|---------|------|
| `.vscode/launch.json` | 新增 | VSCode 調試配置 |
| `.vscode/tasks.json` | 新增 | 自動生成 Swagger 任務 |

### 構建和部署

| 文件 | 變更類型 | 說明 |
|------|---------|------|
| `Makefile` | 修改 | 新增 Docker 和 Swagger 相關命令 |
| `scripts/rebuild-docker-auth.sh` | 新增 | Docker 重建自動化腳本 |
| `scripts/start-auth-service.sh` | 新增 | 本地啟動腳本 |

### 文檔

| 文件 | 變更類型 | 說明 |
|------|---------|------|
| `README_SWAGGER.md` | 新增 | Swagger UI 訪問指南 |
| `docs/VSCODE_DEBUG_GUIDE.md` | 新增 | VSCode 調試完整指南 |
| `docs/QUICK_START.md` | 新增 | 快速啟動指南 |
| `docs/ENVIRONMENT_VARIABLES.md` | 新增 | 環境變數參考 |
| `docs/SWAGGER_URL_CHANGES.md` | 新增 | URL 架構變更說明 |
| `docs/SWAGGER_TROUBLESHOOTING.md` | 新增 | 故障排除指南 |
| `docs/FINAL_SOLUTION_SUMMARY.md` | 新增 | 最終解決方案總結 |
| `docs/DOCKER_SWAGGER_FIX.md` | 新增 | Docker 環境修復文檔 |
| `docs/AUTH_SWAGGER_FIX_SUMMARY.md` | 新增 | 本文檔 |

### 依賴更新

| 文件 | 變更類型 | 說明 |
|------|---------|------|
| `cmd/auth-service/go.mod` | 修改 | 新增 Swagger 相關依賴 |
| `cmd/auth-service/go.sum` | 修改 | 依賴校驗和更新 |
| `go.mod` | 修改 | 更新根模組依賴 |
| `go.sum` | 修改 | 依賴校驗和更新 |

## 🧪 驗證測試

### 本地開發環境測試

```bash
# 1. 啟動依賴服務
docker-compose -f deployments/docker-compose/auth-service.yml up -d postgres-auth redis-auth

# 2. VSCode F5 啟動服務

# 3. 測試端點
curl http://localhost:8081/health                                    # ✅ 200 OK
curl http://localhost:8081/auth-service/swagger/doc.json            # ✅ 200 OK
curl http://localhost:8081/auth-service/swagger/index.html          # ✅ 200 OK
curl -L http://localhost:8081/swagger                                # ✅ 302→200

# 4. 瀏覽器訪問
open http://localhost:8081/auth-service/swagger/index.html
```

**測試結果**: ✅ 所有端點正常，無重定向循環

### Docker 環境測試

```bash
# 1. 重建並啟動服務
make auth-rebuild

# 2. 驗證容器內文件
docker exec tabelogo-auth-service ls -la /app/internal/auth/docs/
# 輸出:
# docs.go         ✅ (10.7 KB)
# index.html      ✅ (1.6 KB)
# swagger.json    ✅ (10.1 KB)
# swagger.yaml    ✅ (5.0 KB)

# 3. 測試端點
curl http://localhost:18080/health                                   # ✅ 200 OK
curl http://localhost:18080/auth-service/swagger/doc.json           # ✅ 200 OK
curl http://localhost:18080/auth-service/swagger/index.html         # ✅ 200 OK
curl -L http://localhost:18080/swagger                               # ✅ 302→200

# 4. 瀏覽器訪問
open http://localhost:18080/auth-service/swagger/index.html
```

**測試結果**: ✅ 所有端點正常，文件存在於正確位置

### 容器日誌檢查

```bash
docker-compose -f deployments/docker-compose/auth-service.yml logs --tail 20 auth-service
```

**正常輸出**:
```
✅ Database connected successfully
✅ Redis connected successfully
✅ Starting gRPC server on port 9090
✅ Starting HTTP server on port 8080
✅ [GIN] 2025/11/24 | 200 | GET "/auth-service/swagger/index.html"
✅ [GIN] 2025/11/24 | 200 | GET "/auth-service/swagger/doc.json"
```

## 📊 URL 對比表

| 環境 | HTTP 端口 | gRPC 端口 | Swagger URL |
|------|----------|----------|-------------|
| **本地開發 (VSCode)** | 8081 | 9091 | http://localhost:8081/auth-service/swagger/index.html |
| **Docker (本地測試)** | 18080 | 19090 | http://localhost:18080/auth-service/swagger/index.html |
| **Docker (生產環境)** | 8080 | 9090 | http://localhost:8080/auth-service/swagger/index.html |

## 🎯 Makefile 新增命令

```bash
# Swagger 文檔生成
make swagger-auth         # 生成 Auth Service Swagger 文檔

# Docker 管理
make auth-build           # 構建 Docker 鏡像
make auth-rebuild         # 完整重建（停止→構建→啟動→測試）
make auth-up              # 啟動服務
make auth-down            # 停止服務
make auth-restart         # 重啟服務
make auth-logs            # 查看日誌
make auth-ps              # 查看狀態
make auth-clean           # 清理容器和數據
make auth-shell           # 進入容器

# 本地開發
make auth-dev             # 本地開發模式（自動生成 Swagger + 啟動）
```

## 🔧 技術細節

### 為什麼 c.File() 會導致重定向循環？

**Go 的 http.ServeFile 行為**:
1. 當請求路徑為 `/auth-service/swagger/index.html`
2. `http.ServeFile` 檢測到這是一個嵌套路徑
3. 自動返回 301 重定向到 `./`（相對路徑）
4. 瀏覽器解析 `./` 為 `/auth-service/swagger/`
5. Gin 的路由匹配失敗，可能觸發其他重定向規則
6. 形成無限循環

**解決方案**:
- 使用 `os.ReadFile()` 直接讀取文件內容
- 使用 `c.Data()` 以 HTTP 響應方式返回內容
- 完全繞過 `http.ServeFile` 的自動行為

### 為什麼要禁用 Gin 的自動重定向？

```go
router.RedirectTrailingSlash = false  // 禁用 /path/ → /path 的自動重定向
router.RedirectFixedPath = false      // 禁用路徑修正的自動重定向
```

**原因**:
- Gin 默認會自動修正 URL（例如添加或移除尾隨斜槓）
- 這可能與我們的 Swagger 路徑規則衝突
- 在重定向場景中可能引發額外的 302/301 響應
- 禁用後路由行為更加可預測

### Docker 多階段構建細節

```dockerfile
# Builder Stage
FROM golang:1.24-alpine AS builder
RUN swag init --output internal/auth/docs ...    # ✅ 生成到正確位置

# Runtime Stage
FROM alpine:3.19
COPY --from=builder /app/internal/auth/docs ./internal/auth/docs  # ✅ 保持結構
```

**關鍵點**:
- Builder 階段生成 Swagger 文檔
- Runtime 階段必須保持相同的目錄結構
- 代碼使用相對路徑 `./internal/auth/docs/`，因此容器內必須匹配

## 📚 相關文檔索引

1. [README_SWAGGER.md](../README_SWAGGER.md) - Swagger UI 訪問指南（主文檔）
2. [VSCODE_DEBUG_GUIDE.md](./VSCODE_DEBUG_GUIDE.md) - VSCode 調試完整指南
3. [QUICK_START.md](./QUICK_START.md) - 快速啟動指南
4. [ENVIRONMENT_VARIABLES.md](./ENVIRONMENT_VARIABLES.md) - 環境變數參考
5. [SWAGGER_TROUBLESHOOTING.md](./SWAGGER_TROUBLESHOOTING.md) - 故障排除指南
6. [DOCKER_SWAGGER_FIX.md](./DOCKER_SWAGGER_FIX.md) - Docker 環境修復詳細說明
7. [FINAL_SOLUTION_SUMMARY.md](./FINAL_SOLUTION_SUMMARY.md) - 最終解決方案（詳細版）

## 🎓 經驗教訓

### 1. 路徑一致性至關重要

**教訓**: 本地開發和 Docker 環境必須使用相同的相對路徑結構
- ✅ 使用相對路徑 `./internal/auth/docs/`
- ❌ 避免硬編碼絕對路徑 `/app/internal/auth/docs/`

### 2. HTTP 文件服務的隱藏行為

**教訓**: `http.ServeFile` 和 `c.File()` 有自動重定向行為
- ✅ 需要精確控制時使用 `os.ReadFile()` + `c.Data()`
- ❌ 避免在可能觸發重定向的場景使用 `c.File()`

### 3. 環境變數命名規範

**教訓**: 必須與代碼中的配置結構體欄位名稱完全匹配
- ✅ 查看 `pkg/config/config.go` 確認正確的變數名
- ❌ 不要憑猜測或慣例命名環境變數

### 4. 服務路徑命名空間

**教訓**: 微服務架構中應使用服務特定的路徑前綴
- ✅ `/auth-service/swagger/` 清晰且不會衝突
- ❌ `/swagger/` 在多服務環境中容易衝突

### 5. Docker 構建驗證

**教訓**: 修改 Dockerfile 後必須重新構建並驗證文件結構
- ✅ 使用 `docker exec` 檢查容器內文件
- ✅ 使用 `--no-cache` 確保完全重建
- ❌ 不要假設文件會自動更新

## 🚀 後續建議

### 1. 自動化測試

建議添加 Swagger 端點的自動化測試：

```go
// internal/auth/interfaces/http/module_test.go
func TestSwaggerEndpoints(t *testing.T) {
    router := setupRouter()

    tests := []struct {
        path       string
        wantStatus int
    }{
        {"/health", http.StatusOK},
        {"/auth-service/swagger/doc.json", http.StatusOK},
        {"/auth-service/swagger/index.html", http.StatusOK},
        {"/swagger", http.StatusFound},
    }

    for _, tt := range tests {
        req := httptest.NewRequest("GET", tt.path, nil)
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)
        assert.Equal(t, tt.wantStatus, w.Code)
    }
}
```

### 2. CI/CD 集成

在 CI/CD pipeline 中添加 Swagger 文檔驗證：

```yaml
# .github/workflows/ci.yml
- name: Generate and Verify Swagger
  run: |
    make swagger-auth
    test -f internal/auth/docs/swagger.json
    test -f internal/auth/docs/index.html
```

### 3. 其他服務應用

將此解決方案應用到其他微服務：
- Restaurant Service
- Booking Service
- Mail Service
- Spider Service

每個服務使用自己的路徑前綴：
- `/restaurant-service/swagger/`
- `/booking-service/swagger/`
- 等等

### 4. API Gateway 整合

考慮在 API Gateway 層面統一 Swagger UI：

```
GET /docs/auth → 代理到 auth-service:8080/auth-service/swagger/index.html
GET /docs/restaurant → 代理到 restaurant-service:8080/restaurant-service/swagger/index.html
```

## ✅ 修復確認清單

- [x] 本地開發環境 Swagger UI 可訪問
- [x] Docker 環境 Swagger UI 可訪問
- [x] 無重定向循環錯誤
- [x] 環境變數配置正確
- [x] 容器內文件結構正確
- [x] VSCode 調試配置正常工作
- [x] 所有端點返回正確的 HTTP 狀態碼
- [x] 服務日誌無錯誤信息
- [x] Makefile 命令全部可用
- [x] 文檔完整且準確
- [x] 快捷重定向功能正常

## 📞 支援

如遇到問題，請參考：
1. [故障排除指南](./SWAGGER_TROUBLESHOOTING.md)
2. [Docker 修復文檔](./DOCKER_SWAGGER_FIX.md)
3. 運行 `make help` 查看所有可用命令

---

**文檔版本**: 1.0.0
**最後更新**: 2025-11-25
**狀態**: ✅ 已完成並驗證
