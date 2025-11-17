# 多來源餐廳聚合平台 (Tabelogo V2)

一個基於微服務架構的餐廳資訊聚合平台，整合多個餐廳資訊來源，提供餐廳搜尋、預訂、評論等功能。

## 🏗 架構特色

- **微服務架構**：每個服務獨立開發、部署、擴展
- **Database per Service**：每個微服務擁有獨立的資料庫實例
- **DDD 設計**：領域驅動設計，清晰的分層架構
- **Event-Driven**：使用 Kafka 實現事件驅動架構
- **gRPC 通訊**：服務間使用高效的 gRPC 通訊
- **完整監控**：Prometheus + Grafana + Jaeger 可觀測性

## 🎯 核心服務

| 服務 | 端口 | 資料庫 | 說明 |
|------|------|--------|------|
| API Gateway | 8080 | - | 統一入口、路由、認證 |
| Auth Service | 8081/9081 | auth_db | 使用者認證與授權 |
| Restaurant Service | 8082/9082 | restaurant_db | 餐廳資料管理 |
| Booking Service | 8083/9083 | booking_db | 預訂功能 |
| Spider Service | 8084/9084 | spider_db | 爬蟲服務 |
| Mail Service | 8085/9085 | mail_db | 郵件通知 |
| Map Service | 8086/9086 | - | 地圖與導航 |

## 🚀 快速開始

### 前置需求

- Docker & Docker Compose
- Go 1.21+
- Make

### 本地開發環境設定

```bash
# 1. Clone repository
git clone https://github.com/lileon/tabelogov2.git
cd tabelogov2

# 2. 初始化專案（建立 .env 檔案）
make init

# 3. 啟動所有基礎設施（PostgreSQL, Redis, Kafka等）
make up

# 4. 檢查容器狀態
make ps
```

### 可用的 Make 指令

```bash
make help          # 顯示所有可用指令
make init          # 初始化專案
make up            # 啟動所有 Docker 容器
make down          # 停止所有容器
make restart       # 重啟所有容器
make logs          # 查看容器日誌
make ps            # 查看容器狀態
make clean         # 清理所有容器和 volumes
make build         # 建置所有微服務
make test          # 執行所有測試
make lint          # 執行程式碼檢查
make migrate-up    # 執行資料庫 migrations
make migrate-down  # 回滾資料庫 migrations
```

## 🗄️ 資料庫架構

### Database per Service 原則

每個微服務擁有獨立的 PostgreSQL 資料庫實例：

| 資料庫 | 端口 | 用途 |
|--------|------|------|
| auth_db | 5432 | 使用者認證資料 |
| restaurant_db | 5433 | 餐廳主資料 |
| booking_db | 5434 | 預訂資料 |
| spider_db | 5435 | 爬蟲任務與結果 |
| mail_db | 5436 | 郵件佇列與記錄 |

### Redis 配置

使用不同的 Redis Database Number 區分各服務：

- DB 0: Auth Service (Session, Token Blacklist)
- DB 1: Restaurant Service (Restaurant Cache)
- DB 2: Booking Service (Booking Cache)
- DB 3: Spider Service (Rate Limiting, Distributed Lock)
- DB 4: API Gateway (Rate Limiting, API Cache)

## 📊 監控與可觀測性

- **Kafka UI**: http://localhost:8080
- **Grafana**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090

## 🔧 技術棧

- **語言**: Go 1.21+
- **Web Framework**: Gin
- **gRPC**: Protocol Buffers
- **資料庫**: PostgreSQL 15
- **Cache**: Redis 7
- **Message Queue**: Apache Kafka
- **監控**: Prometheus + Grafana + Jaeger
- **日誌**: Zap + OpenTelemetry
- **容器化**: Docker + Docker Compose

## 📁 專案結構

```
tabelogov2/
├── cmd/                      # 各微服務入口（每個都有獨立的 go.mod）
├── internal/                 # 按服務分離的內部程式碼
├── pkg/                      # 共用套件（獨立的 go.mod）
├── api/proto/                # gRPC Protocol Buffers 定義
├── migrations/               # 各服務的資料庫 migrations
├── deployments/              # Docker & Kubernetes 配置
├── scripts/                  # 建置與部署腳本
├── tests/                    # 測試
└── docs/                     # 文檔
```

詳細架構文檔請參考：[architecture.md](docs/architecture.md)

## 🔐 環境變數

複製 `.env.example` 到 `.env` 並修改相關設定：

```bash
cp .env.example .env
```

重要變數：
- `JWT_SECRET`: JWT 簽名密鑰（生產環境務必更換）
- `GOOGLE_MAPS_API_KEY`: Google Maps API 金鑰
- `SMTP_*`: 郵件服務設定

## 🧪 測試

```bash
# 執行所有服務的測試
make test

# 執行特定服務的測試
cd cmd/auth-service && go test ./... -v
```

## 📝 License

MIT License

## 👥 作者

Leon Li
