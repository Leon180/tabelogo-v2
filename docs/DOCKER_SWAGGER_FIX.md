# Docker 環境 Swagger UI 404 問題修復

## 🐛 問題描述

使用 Docker 啟動 auth-service 後，訪問 Swagger UI 出現 404 Not Found 錯誤。

## 🔍 根本原因

Dockerfile 中 Swagger 文檔的複製路徑錯誤：

```dockerfile
# ❌ 錯誤：複製到 ./docs
COPY --from=builder /app/cmd/auth-service/docs ./docs

# 但代碼中使用的路徑是
filepath.Abs("./internal/auth/docs/index.html")
```

### 問題分析

1. **Swagger 生成位置**: `internal/auth/docs/` (在 builder 階段)
2. **Dockerfile 複製位置**: `./docs/` (錯誤的目標路徑)
3. **代碼期望位置**: `./internal/auth/docs/` (相對於 WORKDIR=/app)
4. **結果**: 文件不在預期位置，返回 404

## ✅ 解決方案

### 修改 Dockerfile

```dockerfile
# ✅ 正確：保持目錄結構
COPY --from=builder /app/internal/auth/docs ./internal/auth/docs
```

這樣在容器中的文件結構為：
```
/app/
├── auth-service (二進制文件)
└── internal/
    └── auth/
        └── docs/
            ├── docs.go
            ├── index.html
            ├── swagger.json
            └── swagger.yaml
```

## 📋 完整修復步驟

### 1. 更新 Dockerfile

已修改文件: `cmd/auth-service/Dockerfile` 第 55 行

### 2. 重新構建 Docker 鏡像

```bash
# 方法 1: 使用 docker-compose
docker-compose -f deployments/docker-compose/auth-service.yml build auth-service

# 方法 2: 使用 Makefile
make auth-build

# 方法 3: 直接使用 docker build
docker build -f cmd/auth-service/Dockerfile -t tabelogo-auth-service:latest .
```

### 3. 重啟服務

```bash
# 停止舊容器
docker-compose -f deployments/docker-compose/auth-service.yml down

# 啟動新容器
docker-compose -f deployments/docker-compose/auth-service.yml up -d
```

### 4. 驗證修復

```bash
# 檢查容器狀態
docker ps | grep auth-service

# 測試 Swagger UI
curl http://localhost:18080/auth-service/swagger/index.html

# 或在瀏覽器訪問
open http://localhost:18080/auth-service/swagger/index.html
```

## 🔍 調試方法

### 檢查容器內的文件結構

```bash
# 進入容器
docker exec -it tabelogo-auth-service sh

# 檢查文件是否存在
ls -la /app/internal/auth/docs/

# 檢查文件內容
cat /app/internal/auth/docs/index.html | head -20

# 退出容器
exit
```

**預期輸出**:
```
/app/internal/auth/docs/
├── docs.go
├── index.html
├── swagger.json
└── swagger.yaml
```

### 檢查服務日誌

```bash
# 查看實時日誌
docker logs -f tabelogo-auth-service

# 查看最後 50 行日誌
docker logs --tail 50 tabelogo-auth-service
```

## 📊 對比：本地開發 vs Docker

| 環境 | 工作目錄 | Swagger 文檔路徑 |
|------|---------|-----------------|
| 本地開發 | `/Users/lileon/goproject/tabelogov2` | `./internal/auth/docs/` |
| Docker 容器 | `/app` | `./internal/auth/docs/` |

**關鍵**: 兩個環境使用相同的**相對路徑**，因此目錄結構必須一致。

## ⚠️ 常見錯誤

### 錯誤 1: 使用絕對路徑

```go
// ❌ 不要在代碼中硬編碼絕對路徑
absPath := "/app/internal/auth/docs/index.html"  // 本地開發無法使用

// ✅ 使用相對路徑 + filepath.Abs
absPath, err := filepath.Abs("./internal/auth/docs/index.html")
```

### 錯誤 2: Dockerfile 中路徑不匹配

```dockerfile
# ❌ 錯誤：複製到錯誤的位置
COPY --from=builder /app/internal/auth/docs ./docs

# ✅ 正確：保持原有結構
COPY --from=builder /app/internal/auth/docs ./internal/auth/docs
```

### 錯誤 3: 忘記重新構建鏡像

```bash
# ❌ 錯誤：只重啟容器，沒有重新構建
docker-compose restart

# ✅ 正確：重新構建並啟動
docker-compose down
docker-compose build
docker-compose up -d
```

## 🧪 測試清單

- [ ] Docker 鏡像成功構建
- [ ] 容器成功啟動
- [ ] 容器健康檢查通過
- [ ] Swagger JSON 可訪問: `http://localhost:18080/auth-service/swagger/doc.json`
- [ ] Swagger UI 可訪問: `http://localhost:18080/auth-service/swagger/index.html`
- [ ] 快捷重定向可用: `http://localhost:18080/swagger`
- [ ] API 端點可正常調用

## 📝 端口對比

| 環境 | HTTP 端口 | Swagger URL |
|------|----------|-------------|
| 本地開發 (VSCode) | 8081 | `http://localhost:8081/auth-service/swagger/index.html` |
| Docker (本地測試) | 18080 | `http://localhost:18080/auth-service/swagger/index.html` |
| Docker (生產環境) | 8080 | `http://localhost:8080/auth-service/swagger/index.html` |

## 🔧 快速重建命令

```bash
# 一鍵重建並啟動
docker-compose -f deployments/docker-compose/auth-service.yml down && \
docker-compose -f deployments/docker-compose/auth-service.yml build --no-cache auth-service && \
docker-compose -f deployments/docker-compose/auth-service.yml up -d

# 查看日誌
docker-compose -f deployments/docker-compose/auth-service.yml logs -f auth-service
```

## 🎯 驗證腳本

```bash
#!/bin/bash

echo "🔍 驗證 Auth Service Docker 部署..."

# 等待服務啟動
echo "⏳ 等待服務啟動 (30秒)..."
sleep 30

# 測試 Health Check
echo "✅ 測試 Health Check..."
curl -f http://localhost:18080/health || echo "❌ Health check failed"

# 測試 Swagger JSON
echo "✅ 測試 Swagger JSON..."
curl -f http://localhost:18080/auth-service/swagger/doc.json > /dev/null || echo "❌ Swagger JSON failed"

# 測試 Swagger UI
echo "✅ 測試 Swagger UI..."
curl -f http://localhost:18080/auth-service/swagger/index.html > /dev/null || echo "❌ Swagger UI failed"

# 測試快捷重定向
echo "✅ 測試快捷重定向..."
curl -f -L http://localhost:18080/swagger > /dev/null || echo "❌ Redirect failed"

echo "🎉 所有測試完成！"
```

保存為 `scripts/verify-docker-swagger.sh` 並執行：

```bash
chmod +x scripts/verify-docker-swagger.sh
./scripts/verify-docker-swagger.sh
```

## 📚 相關文檔

- [Dockerfile](../../cmd/auth-service/Dockerfile)
- [Docker Compose](../../deployments/docker-compose/auth-service.yml)
- [Swagger 故障排除](./SWAGGER_TROUBLESHOOTING.md)
- [最終解決方案總結](./FINAL_SOLUTION_SUMMARY.md)

---

**最後更新**: 2025-11-25
**狀態**: ✅ 已修復
**影響範圍**: Docker 部署環境
