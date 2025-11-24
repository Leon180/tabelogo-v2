# Swagger UI 故障排除指南

## 重定向循環問題 (ERR_TOO_MANY_REDIRECTS)

### 問題描述

瀏覽器顯示 "This page isn't working - redirected you too many times"

### 常見原因

#### 1. 路由衝突

**問題**: 靜態文件路由與重定向路由衝突

```go
// ❌ 錯誤: 使用 StaticFile 可能導致路徑匹配問題
router.StaticFile("/auth-service/swagger/index.html", "./internal/auth/docs/index.html")
router.GET("/auth-service/swagger", func(c *gin.Context) {
    c.Redirect(301, "/auth-service/swagger/index.html")
})
```

**解決方案**: 使用 `c.File()` 代替 `StaticFile`

```go
// ✅ 正確: 使用 GET + File 明確處理路徑
router.GET("/auth-service/swagger/index.html", func(c *gin.Context) {
    c.File("./internal/auth/docs/index.html")
})
```

#### 2. 永久重定向緩存

**問題**: 使用 301 (StatusMovedPermanently) 會被瀏覽器緩存

```go
// ❌ 錯誤: 301 會被緩存，即使修復代碼後仍可能有問題
c.Redirect(http.StatusMovedPermanently, "/target")
```

**解決方案**: 使用 302 (StatusFound) 臨時重定向

```go
// ✅ 正確: 302 不會被緩存
c.Redirect(http.StatusFound, "/target")
```

#### 3. 末尾斜杠匹配

**問題**: `/swagger` 和 `/swagger/` 可能被視為不同路徑

```go
// ❌ 不完整: 只處理一種情況
router.GET("/swagger", func(c *gin.Context) {
    c.Redirect(302, "/auth-service/swagger/index.html")
})
// 訪問 /swagger/ 會 404
```

**解決方案**: 同時處理兩種情況

```go
// ✅ 正確: 處理所有變體
router.GET("/swagger", handler)
router.GET("/swagger/", handler)
router.GET("/auth-service/swagger", handler)
router.GET("/auth-service/swagger/", handler)
```

### 完整的正確實現

```go
// Swagger JSON endpoint
router.GET("/auth-service/swagger/doc.json", func(c *gin.Context) {
    c.String(200, docs.SwaggerInfo.ReadDoc())
})

// Serve Swagger UI HTML file (使用 File 而非 StaticFile)
router.GET("/auth-service/swagger/index.html", func(c *gin.Context) {
    c.File("./internal/auth/docs/index.html")
})

// Redirect shortcuts (使用 302 而非 301)
router.GET("/swagger", func(c *gin.Context) {
    c.Redirect(http.StatusFound, "/auth-service/swagger/index.html")
})
router.GET("/swagger/", func(c *gin.Context) {
    c.Redirect(http.StatusFound, "/auth-service/swagger/index.html")
})
router.GET("/auth-service/swagger", func(c *gin.Context) {
    c.Redirect(http.StatusFound, "/auth-service/swagger/index.html")
})
router.GET("/auth-service/swagger/", func(c *gin.Context) {
    c.Redirect(http.StatusFound, "/auth-service/swagger/index.html")
})
```

## 調試步驟

### 1. 清除瀏覽器緩存

```bash
# Chrome/Edge
1. 打開開發者工具 (F12)
2. 右鍵點擊刷新按鈕
3. 選擇 "清空緩存並硬性重新載入"

# 或使用無痕模式測試
Ctrl+Shift+N (Windows) / Cmd+Shift+N (Mac)
```

### 2. 使用 curl 測試重定向鏈

```bash
# 測試主要端點
curl -I http://localhost:8081/auth-service/swagger/index.html

# 測試重定向
curl -L -I http://localhost:8081/swagger

# 顯示重定向鏈
curl -L -v http://localhost:8081/swagger 2>&1 | grep -E '(< HTTP|< Location)'
```

**預期輸出**:
```
< HTTP/1.1 302 Found
< Location: /auth-service/swagger/index.html
< HTTP/1.1 200 OK
```

**錯誤輸出** (循環):
```
< HTTP/1.1 302 Found
< Location: /auth-service/swagger/index.html
< HTTP/1.1 302 Found
< Location: /auth-service/swagger/index.html
< HTTP/1.1 302 Found
...
```

### 3. 檢查路由註冊

啟動服務後，檢查 Gin 的路由表：

```go
// 在 RegisterRoutes 函數末尾添加（僅用於調試）
for _, route := range router.Routes() {
    logger.Info("Registered route",
        zap.String("method", route.Method),
        zap.String("path", route.Path),
    )
}
```

**預期路由**:
```
GET /auth-service/swagger/doc.json
GET /auth-service/swagger/index.html
GET /swagger
GET /swagger/
GET /auth-service/swagger
GET /auth-service/swagger/
```

### 4. 檢查文件是否存在

```bash
# 確認 index.html 存在
ls -la ./internal/auth/docs/index.html

# 檢查內容
head -20 ./internal/auth/docs/index.html

# 確認 JSON URL 配置正確
grep "url:" ./internal/auth/docs/index.html
```

**應該看到**:
```javascript
url: "/auth-service/swagger/doc.json",
```

## 常見錯誤模式

### 模式 1: 無限重定向

```
/swagger → /swagger/index.html → /swagger → ...
```

**原因**: `/swagger/index.html` 被錯誤地匹配到 `/swagger` 路由

**修復**: 使用精確路徑匹配，避免使用 `StaticFile`

### 模式 2: 404 後重定向

```
/swagger → /auth-service/swagger/index.html → 404 → /swagger → ...
```

**原因**: 文件路徑不正確或文件不存在

**修復**:
```bash
# 檢查當前工作目錄
pwd

# 確認相對路徑正確
ls -la ./internal/auth/docs/index.html
```

### 模式 3: CORS 錯誤導致重定向

```
/auth-service/swagger/index.html → CORS error → redirect → ...
```

**原因**: 靜態文件和 API 的 CORS 配置不一致

**修復**: 確保靜態文件和 JSON 端點有相同的 CORS 設置

## 最佳實踐

### 1. 使用 File() 而非 StaticFile()

```go
// ✅ 推薦
router.GET("/path/to/file.html", func(c *gin.Context) {
    c.File("./path/to/file.html")
})

// ❌ 避免（可能導致路徑匹配問題）
router.StaticFile("/path/to/file.html", "./path/to/file.html")
```

### 2. 使用 302 而非 301

```go
// ✅ 開發環境推薦
c.Redirect(http.StatusFound, "/target")  // 302

// ⚠️ 僅生產環境使用
c.Redirect(http.StatusMovedPermanently, "/target")  // 301
```

### 3. 處理所有路徑變體

```go
// ✅ 完整
router.GET("/swagger", handler)
router.GET("/swagger/", handler)

// ❌ 不完整
router.GET("/swagger", handler)
```

### 4. 添加調試日誌

```go
router.GET("/auth-service/swagger/index.html", func(c *gin.Context) {
    logger.Debug("Serving Swagger UI", zap.String("path", c.Request.URL.Path))
    c.File("./internal/auth/docs/index.html")
})
```

### 5. 使用絕對路徑（生產環境）

```go
// 開發環境
c.File("./internal/auth/docs/index.html")

// 生產環境（Docker）
c.File("/app/internal/auth/docs/index.html")

// 或使用環境變量
docsPath := os.Getenv("SWAGGER_DOCS_PATH")
if docsPath == "" {
    docsPath = "./internal/auth/docs"
}
c.File(filepath.Join(docsPath, "index.html"))
```

## 測試清單

在修復後，確保測試所有路徑：

- [ ] `http://localhost:8081/auth-service/swagger/index.html` - 直接訪問
- [ ] `http://localhost:8081/auth-service/swagger/doc.json` - JSON API
- [ ] `http://localhost:8081/swagger` - 快捷重定向
- [ ] `http://localhost:8081/swagger/` - 帶斜杠的快捷重定向
- [ ] `http://localhost:8081/auth-service/swagger` - 服務重定向
- [ ] `http://localhost:8081/auth-service/swagger/` - 帶斜杠的服務重定向
- [ ] 清除緩存後重新測試
- [ ] 使用無痕模式測試
- [ ] 使用 curl 測試重定向鏈

## 快速修復命令

```bash
# 1. 停止服務
# 在 VSCode 中按 Shift+F5

# 2. 清理可能的緩存
rm -rf ./internal/auth/docs/docs.go.bak

# 3. 重新生成 Swagger
make swagger-auth

# 4. 檢查文件
ls -la ./internal/auth/docs/

# 5. 重啟服務
# 在 VSCode 中按 F5

# 6. 測試
curl -L -v http://localhost:8081/swagger 2>&1 | grep -E '(< HTTP|< Location)'
```

## 參考資料

- [Gin Web Framework - Static Files](https://gin-gonic.com/docs/examples/serving-static-files/)
- [HTTP Status Codes](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status)
- [Chrome DevTools Network](https://developer.chrome.com/docs/devtools/network/)
