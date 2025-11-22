# Auth Service Deployment Guide

## 📦 Docker 部署

### 方式 1: 使用 Docker Compose（推薦用於開發）

```bash
cd cmd/auth-service

# 複製環境配置
cp .env.example .env

# 編輯 .env 設定 JWT_SECRET
# JWT_SECRET 必須至少 32 字元

# 啟動所有服務（PostgreSQL + Redis + Auth Service）
docker-compose up -d

# 查看日誌
docker-compose logs -f auth-service

# 停止服務
docker-compose down

# 停止並刪除資料
docker-compose down -v
```

### 方式 2: 單獨構建 Docker Image

```bash
# 在專案根目錄執行
docker build -f cmd/auth-service/Dockerfile -t tabelogo-auth-service:latest .

# 運行容器（需要先啟動 PostgreSQL 和 Redis）
docker run -d \
  --name auth-service \
  -p 8080:8080 \
  -p 9090:9090 \
  -e DB_HOST=postgres \
  -e DB_NAME=auth_db \
  -e REDIS_HOST=redis \
  -e JWT_SECRET=your-secret-key-min-32-chars \
  tabelogo-auth-service:latest
```

## 🚀 本地開發

### 前置需求

- Go 1.23+
- PostgreSQL 15+
- Redis 7+

### 步驟

```bash
cd cmd/auth-service

# 1. 複製環境配置
cp .env.example .env

# 2. 編輯 .env
# 設定 DB_NAME, JWT_SECRET 等

# 3. 啟動資料庫（使用 Docker）
docker-compose up -d postgres-auth redis-auth

# 4. 執行 Migration（如果有）
# make migrate-up

# 5. 運行服務
go run main.go

# 或編譯後運行
GOWORK=off go build -o ../../bin/auth-service .
../../bin/auth-service
```

## ☸️ Kubernetes 部署

### ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: auth-service-config
data:
  ENVIRONMENT: "production"
  LOG_LEVEL: "info"
  SERVER_PORT: "8080"
  GRPC_PORT: "9090"
  DB_HOST: "postgres-service"
  DB_PORT: "5432"
  DB_NAME: "auth_db"
  REDIS_HOST: "redis-service"
  REDIS_PORT: "6379"
```

### Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: auth-service-secret
type: Opaque
stringData:
  DB_PASSWORD: "your-db-password"
  REDIS_PASSWORD: "your-redis-password"
  JWT_SECRET: "your-jwt-secret-min-32-characters"
```

### Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: auth-service
  template:
    metadata:
      labels:
        app: auth-service
    spec:
      containers:
      - name: auth-service
        image: tabelogo-auth-service:latest
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 9090
          name: grpc
        envFrom:
        - configMapRef:
            name: auth-service-config
        - secretRef:
            name: auth-service-secret
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

### Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: auth-service
spec:
  selector:
    app: auth-service
  ports:
  - name: http
    port: 8080
    targetPort: 8080
  - name: grpc
    port: 9090
    targetPort: 9090
  type: ClusterIP
```

## 🔧 環境變數說明

### 必需變數

| 變數名 | 說明 | 範例 |
|--------|------|------|
| `DB_NAME` | 資料庫名稱 | `auth_db` |
| `JWT_SECRET` | JWT 簽名密鑰（最少 32 字元） | `your-secret-key-min-32-chars` |

### 可選變數

| 變數名 | 說明 | 預設值 |
|--------|------|--------|
| `ENVIRONMENT` | 環境模式 | `development` |
| `LOG_LEVEL` | 日誌級別 | `info` |
| `SERVER_PORT` | HTTP 端口 | `8080` |
| `GRPC_PORT` | gRPC 端口 | `9090` |
| `DB_HOST` | PostgreSQL 主機 | `localhost` |
| `DB_PORT` | PostgreSQL 端口 | `5432` |
| `DB_USER` | 資料庫用戶 | `postgres` |
| `DB_PASSWORD` | 資料庫密碼 | `postgres` |
| `REDIS_HOST` | Redis 主機 | `localhost` |
| `REDIS_PORT` | Redis 端口 | `6379` |
| `JWT_ACCESS_TOKEN_EXPIRE` | Access Token 過期時間 | `15m` |
| `JWT_REFRESH_TOKEN_EXPIRE` | Refresh Token 過期時間 | `168h` |

## 🧪 健康檢查

```bash
# HTTP Health Check
curl http://localhost:8080/health

# 預期回應
{"status":"ok"}
```

## 📊 監控

### Prometheus Metrics

服務暴露 Prometheus metrics（如果已配置）：

```
http://localhost:8080/metrics
```

### 日誌

服務使用結構化日誌（zap），輸出 JSON 格式：

```bash
# 查看 Docker 日誌
docker-compose logs -f auth-service

# 查看 Kubernetes 日誌
kubectl logs -f deployment/auth-service
```

## 🔒 安全建議

1. **生產環境必須更改 JWT_SECRET**
   - 使用強隨機字串（至少 32 字元）
   - 定期輪換密鑰

2. **使用 Secret Management**
   - Kubernetes Secrets
   - HashiCorp Vault
   - AWS Secrets Manager

3. **啟用 TLS/SSL**
   - 資料庫連線使用 SSL
   - 使用 HTTPS/gRPC TLS

4. **限制網路訪問**
   - 使用防火牆規則
   - 配置 Network Policies

## 📝 故障排除

### 服務無法啟動

```bash
# 檢查日誌
docker-compose logs auth-service

# 常見問題：
# 1. 資料庫連線失敗 -> 檢查 DB_HOST, DB_NAME
# 2. JWT_SECRET 太短 -> 至少 32 字元
# 3. 端口被佔用 -> 修改 SERVER_PORT, GRPC_PORT
```

### 資料庫連線問題

```bash
# 測試資料庫連線
docker-compose exec postgres-auth psql -U postgres -d auth_db

# 檢查資料庫是否存在
\l

# 檢查表是否存在
\dt
```

### Redis 連線問題

```bash
# 測試 Redis 連線
docker-compose exec redis-auth redis-cli ping

# 預期回應: PONG
```
