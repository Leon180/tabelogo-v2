# Tabelogo v2 - Docker Compose 架構說明

## 架構概述

本專案採用微服務架構，使用 Docker Compose 進行容器編排。有兩種 docker-compose 配置：

### 1. 根目錄 `docker-compose.yml` - 完整系統

**用途：** 啟動整個 Tabelogo v2 系統，包含所有微服務和基礎設施

**包含服務：**
- **基礎設施**
  - PostgreSQL (Auth, Restaurant, Booking 各自獨立)
  - Redis (共享快取)
  - Kafka + Zookeeper (訊息佇列)
  - Prometheus (監控)
  - Grafana (視覺化)

- **微服務**
  - Auth Service (8080/HTTP, 9090/gRPC)
  - Restaurant Service (待實作)
  - Booking Service (待實作)
  - API Gateway (待實作)

**啟動方式：**
```bash
# 在專案根目錄
make up              # 啟動所有服務
make down            # 停止所有服務
make ps              # 查看服務狀態
make logs            # 查看所有日誌
```

### 2. 各服務目錄 `cmd/*/docker-compose.yml` - 單服務開發

**用途：** 僅用於單一服務的本地開發和測試

**特點：**
- 使用不同的端口避免衝突
- 只啟動該服務及其直接依賴
- 適合快速迭代開發

**範例 - Auth Service 本地開發：**
```bash
cd cmd/auth-service
docker-compose up -d    # 啟動 Auth Service (端口 18080/19090)
docker-compose down     # 停止
```

或使用 Makefile：
```bash
make auth-up            # 啟動 Auth Service (本地開發模式)
make auth-down          # 停止
make auth-logs          # 查看日誌
```

## 端口分配

### 根目錄 docker-compose (生產模式)
| 服務 | HTTP | gRPC | 其他 |
|------|------|------|------|
| Auth Service | 8080 | 9090 | - |
| Restaurant Service | 8081 | 9091 | - |
| Booking Service | 8082 | 9092 | - |
| API Gateway | 8000 | - | - |
| PostgreSQL (Auth) | 5432 | - | - |
| PostgreSQL (Restaurant) | 5433 | - | - |
| PostgreSQL (Booking) | 5434 | - | - |
| Redis | 6379 | - | - |
| Kafka | - | - | 9092 |
| Prometheus | - | - | 9090 |
| Grafana | - | - | 3000 |

### 服務目錄 docker-compose (開發模式)
| 服務 | HTTP | gRPC | DB | Redis |
|------|------|------|-----|-------|
| Auth Service | 18080 | 19090 | 15432 | 16379 |

## 使用場景

### 場景 1: 完整系統測試
```bash
# 啟動整個系統
make up

# 測試服務間通訊
curl http://localhost:8080/health
```

### 場景 2: 單服務開發
```bash
# 只開發 Auth Service
cd cmd/auth-service
docker-compose up -d

# 或使用 Makefile
make auth-up
```

### 場景 3: 新增微服務
1. 在 `cmd/new-service/` 建立服務
2. 建立 `cmd/new-service/Dockerfile`
3. 建立 `cmd/new-service/docker-compose.yml` (開發用)
4. 在根目錄 `docker-compose.yml` 加入服務定義
5. 更新 `Makefile` 加入相關指令

## 網路架構

所有服務都在同一個 Docker 網路 `tabelogo-network` 中，可以通過服務名稱互相通訊：

```yaml
# 範例：Restaurant Service 呼叫 Auth Service
AUTH_SERVICE_URL: http://auth-service:8080
AUTH_SERVICE_GRPC: auth-service:9090
```

## 資料持久化

所有資料庫和快取都使用 Docker Volumes 持久化：

```bash
# 查看 volumes
docker volume ls | grep tabelogo

# 清理所有資料 (危險！)
make clean
```

## 最佳實踐

1. **開發時**：使用服務目錄的 docker-compose
2. **整合測試**：使用根目錄的 docker-compose
3. **生產部署**：使用 Kubernetes (未來)
4. **端口衝突**：確保開發模式使用不同端口

## 故障排除

### 端口已被佔用
```bash
# 檢查端口佔用
lsof -i :8080

# 使用開發模式（不同端口）
make auth-up
```

### 服務無法啟動
```bash
# 查看日誌
make logs

# 或特定服務
make auth-logs
```

### 清理並重新開始
```bash
# 停止並刪除所有容器和資料
make clean

# 重新啟動
make up
```
