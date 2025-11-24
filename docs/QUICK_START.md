# Auth Service 快速啟動指南

## 🚀 一鍵啟動（推薦）

### 1. 啟動依賴服務
```bash
docker-compose -f deployments/docker-compose/auth-service.yml up -d postgres-auth redis-auth
```

### 2. 在 VSCode 中啟動調試
- 按 **F5**
- 選擇 **"Auth Service"**
- 等待服務啟動完成

### 3. 訪問 Swagger UI
```
http://localhost:8081/auth-service/swagger/index.html
```

或者使用快捷方式：
```
http://localhost:8081/swagger
```

就這麼簡單！🎉

---

## 📋 詳細說明

### 服務端口配置

本地開發環境使用專用端口避免衝突：

| 服務 | 標準端口 | 本地開發端口 |
|------|---------|-------------|
| PostgreSQL | 5432 | **15432** |
| Redis | 6379 | **16379** |
| Auth HTTP API | 8080 | **8081** |
| Auth gRPC API | 9090 | **9091** |

### 檢查服務狀態

```bash
# 查看所有運行中的容器
docker ps

# 查看 Auth Service 相關容器
docker ps | grep -E "(postgres-auth|redis-auth)"

# 查看容器日誌
docker logs tabelogo-postgres-auth-dev
docker logs tabelogo-redis-auth-dev
```

### 測試數據庫連接

```bash
# PostgreSQL
docker exec -it tabelogo-postgres-auth-dev psql -U postgres -d auth_db

# Redis
docker exec -it tabelogo-redis-auth-dev redis-cli
```

### 停止服務

```bash
# 停止所有服務
docker-compose -f deployments/docker-compose/auth-service.yml down

# 僅停止 Auth Service（保留數據庫）
# 在 VSCode 中按 Shift+F5 停止調試

# 清理所有數據（包括數據庫數據）
docker-compose -f deployments/docker-compose/auth-service.yml down -v
```

---

## 🔧 其他啟動方式

### 方式 2: 使用 Makefile
```bash
# 啟動 Auth Service（自動生成 Swagger + 啟動服務）
make auth-dev
```

### 方式 3: 使用啟動腳本
```bash
./scripts/start-auth-service.sh
```

---

## 📚 API 測試

### Health Check
```bash
curl http://localhost:8081/health
```

### 註冊用戶
```bash
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 登入
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 驗證 Token
```bash
curl -X GET http://localhost:8081/api/v1/auth/validate \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

## ❓ 常見問題

### 無法訪問 Swagger UI？
確保：
1. ✅ 服務已啟動（檢查 VSCode Debug Console）
2. ✅ Swagger 文檔已生成（自動執行 preLaunchTask）
3. ✅ 訪問正確的 URL: http://localhost:8081/auth-service/swagger/index.html

### 數據庫連接失敗？
```bash
# 檢查 PostgreSQL 是否運行
docker ps | grep postgres-auth

# 如果未運行，啟動它
docker-compose -f deployments/docker-compose/auth-service.yml up -d postgres-auth

# 檢查端口是否正確（應該是 15432）
```

### Redis 連接失敗？
```bash
# 檢查 Redis 是否運行
docker ps | grep redis-auth

# 如果未運行，啟動它
docker-compose -f deployments/docker-compose/auth-service.yml up -d redis-auth

# 檢查端口是否正確（應該是 16379）
```

---

## 📖 更多資訊

詳細的調試指南和故障排除，請參閱：
- [VSCode 調試指南](./VSCODE_DEBUG_GUIDE.md)

---

**提示**: 第一次啟動需要下載 Docker 鏡像，可能需要幾分鐘。後續啟動會很快！⚡
