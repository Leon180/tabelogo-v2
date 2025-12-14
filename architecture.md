# 多來源餐廳聚合平台 - 完整架構設計

## 1. 專案概述

一個整合多個餐廳資訊來源的聚合平台，提供餐廳搜尋、預訂等功能。採用微服務架構，展示分散式系統設計與實作能力。

---

## 2. 核心功能模組

### 2.1 功能服務實現狀態

**最後更新**: 2025-12-11

| 服務 | 狀態 | 完成度 | 主要功能 |
|------|------|--------|------------|
| **Auth Service** | ✅ 已實現 | 100% | HTTP API, gRPC API, Swagger UI, JWT 認證 |
| **Restaurant Service** | ✅ 已實現 | 100% | HTTP API, gRPC API, Map 整合, Cache-First, 98% 測試覆蓋 |
| **Map Service** | ✅ 已實現 | 100% | HTTP API, gRPC API, Swagger UI, Prometheus, Google Maps 整合 |
| **Phase 2 Integration** | ✅ 已完成 | 100% | Restaurant-Map 智能快取整合、80% 成本降低 |
| **Booking Service** | ⏳ 規劃中 | 0% | 餐廳預訂功能（整合 OpenTable API）|
| **Spider Service** | 🚧 開發中 | 40% | HTTP API, gRPC API, Tabelog 爬蟲, Redis 快取, DTO 模式 |
| **Mail Service** | ⏳ 規劃中 | 0% | 郵件通知服務 |
| **API Gateway** | ⏳ 規劃中 | 0% | 統一入口、路由、認證 |

**整體完成度**: **49%** (3.4/7 服務已實現)

### 2.2 已實現服務詳情

#### Auth Service ✅
- ✅ HTTP RESTful API
- ✅ gRPC API (4 RPC methods)
- ✅ Swagger UI 文檔
- ✅ PostgreSQL 資料庫 (port 5435)
- ✅ Redis 快取 (DB 0)
- ✅ JWT 認證與授權
- ✅ Docker 容器化
- ✅ 完整測試覆蓋

#### Restaurant Service ✅ **（功能最完整）**
- ✅ HTTP RESTful API (6 endpoints)
- ✅ gRPC API (10 RPC methods)
- ✅ Swagger UI 文檔（範圍已修復）
- ✅ Prometheus Metrics 監控
- ✅ PostgreSQL 資料庫 (port 5433)
- ✅ Redis 快取 (DB 1)
- ✅ Docker 容器化
- ✅ 98% 測試覆蓋率
- ✅ DDD 架構完整實現

#### Map Service ✅ **（Phase 1 完成）**
- ✅ HTTP RESTful API (3 endpoints)
- ✅ gRPC API (3 RPC methods: QuickSearch, AdvanceSearch, BatchGetPlaces)
- ✅ Swagger UI 文檔
- ✅ Prometheus Metrics 監控
- ✅ Google Maps API 整合
- ✅ Quick Search 功能
- ✅ Advance Search 功能
- ✅ Redis 快取 (DB 5)
- ✅ Docker 容器化 (HTTP: 8081, gRPC: 19083)
- ✅ 健康檢查與日誌
- ✅ **Phase 2**: Restaurant Service 整合（已完成）

#### Spider Service 🚧 **（Phase 1 進行中）**
- ✅ HTTP RESTful API (3 endpoints: scrape, job status, stream)
- ✅ gRPC API (2 RPC methods: SearchSimilarRestaurants, GetRestaurantPhotos)
- ✅ Tabelog 爬蟲實現 (使用 colly)
- ✅ Redis 快取 (DB 2) - DTO 模式實現
- ✅ Domain Model 封裝 (私有字段 + getter 方法)
- ✅ DTO 層處理 JSON 序列化
- ✅ Docker 容器化 (HTTP: 18084, gRPC: 19084)
- ⏳ 非同步任務處理 (Job Queue + SSE)
- ⏳ 錯誤處理與重試機制 (Circuit Breaker)
- ⏳ Rate Limiting (防止被封鎖)
- ⏳ Prometheus Metrics 監控

### 2.3 Phase 2: Map-Restaurant 服務整合 ✅ **（已完成）**

**目標**: 實現智能餐廳搜尋，優先使用本地快取，減少 Google API 調用成本

#### 核心需求

**Quick Search 智能流程**:
1. 使用者發起 `quick_search` 請求（通過 place_id）
2. Restaurant Service 優先查詢本地資料庫
3. 如果找到 → 直接返回（節省 API 成本）
4. 如果未找到 → 調用 Map Service 的 Google API
5. Map Service 返回結果後，自動同步到 Restaurant Service
6. 下次相同查詢直接命中本地快取

#### 架構設計

```
User Request (place_id)
    ↓
Restaurant Service
    ├─ Check Local DB (restaurants table)
    │   ├─ Found → Return (Cache Hit)
    │   └─ Not Found → Continue
    ↓
Call Map Service (gRPC)
    ↓
Map Service
    ├─ Check Redis Cache
    │   ├─ Hit → Return
    │   └─ Miss → Call Google API
    ↓
Return to Restaurant Service
    ↓
Restaurant Service
    ├─ Save to Local DB
    └─ Return to User
```

#### 實現要點

**Restaurant Service 端**:
- ✅ 新增 gRPC Client 調用 Map Service
- ✅ 實現資料新鮮度檢查邏輯
- ✅ 自動同步 Google 資料到本地 DB
- ✅ 處理同步失敗的降級策略

**Map Service 端**:
- ✅ gRPC API 已就緒 (QuickSearch, BatchGetPlaces)
- ✅ Redis 快取機制已實現
- ✅ Google API 整合完成

**資料同步策略**:
- **即時同步**: quick_search 結果立即寫入 Restaurant DB
- **批次同步**: 定期使用 BatchGetPlaces 更新熱門餐廳
- **TTL 管理**: 設定資料過期時間，定期重新驗證

#### 預期效益

- 📉 **降低成本**: 減少 80% Google API 調用
- ⚡ **提升速度**: 本地查詢 < 10ms vs Google API ~200ms
- 🎯 **提高可用性**: Google API 故障時仍可服務
- 📊 **數據積累**: 建立自有餐廳資料庫

---

## 3. 技術架構

### 3.1 架構模式
- **微服務架構 (Microservices)**
  - 服務間通訊：gRPC (內部服務通訊)
  - 對外 API：RESTful API
  - 服務發現：考慮使用 Consul 或 etcd
  - API Gateway：統一入口、路由、認證

- **領域驅動設計 (DDD)**
  - 分層架構：Presentation → Application → Domain → Infrastructure
  - 聚合根 (Aggregate Root) 設計
  - Repository Pattern
  - Value Objects

### 3.2 核心技術棧
- **主要語言**: Go 1.21+
- **依賴注入**: Uber FX
- **Web Framework**: Gin
- **RPC Framework**: gRPC with Protocol Buffers
- **Message Queue**: Apache Kafka
- **資料庫**: PostgreSQL 15+ (使用 GORM)
- **Cache**: Redis 7+
- **並發控制**: Goroutines、Channels、Context

### 3.3 前端架構
- **Framework**: Next.js 16 (App Router)
- **Language**: TypeScript
- **Styling**: TailwindCSS v4 + Shadcn/UI
- **Maps**: @vis.gl/react-google-maps
- **State**: React Query (Server State), React Hooks (Local State)
- **Design Pattern**: Map-First Interface


---

## 4. 資料層設計

### 4.1 Database per Service 原則
遵循微服務架構最佳實踐，**每個微服務擁有獨立的資料庫**，實現真正的服務解耦。

#### 資料庫分配策略 ✅ **已實現 (2025-11-20)**

| 服務 | 資料庫名稱 | 端口 | 主要資料表 | 說明 | 狀態 |
|------|-----------|------|-----------|------|------|
| **Auth Service** | `auth_db` | **15432** ⚠️ | users, refresh_tokens | 使用者認證資料 | ✅ |
| **Restaurant Service** | `restaurant_db` | 5433 | restaurants, user_favorites | 餐廳主資料（來自外部）、使用者收藏 | ✅ |
| **Booking Service** | `booking_db` | 5434 | bookings, booking_history | 預訂資料（Event Sourcing） | ✅ |
| **Spider Service** | `spider_db` | 5435 | crawl_jobs, crawl_results | 爬蟲任務與結果（Google/Tabelog/IG） | ✅ |
| **Mail Service** | `mail_db` | 5436 | email_queue, email_logs | 郵件佇列與追蹤記錄 | ✅ |
| **Map Service** | 無獨立 DB | - | - | 僅作為 Google Maps API 的代理層 | - |

**設計調整說明**：
- ❌ 移除 `reviews` 表 - 評論來自外部資料源（Google/Tabelog），不需本地儲存
- ✅ 新增 `user_favorites` 表 - 使用者僅能收藏、查詢餐廳，不能編輯資料
- ✅ `bookings` 支援外部 API 同步 - 增加 `external_booking_id`, `external_service`, `last_synced_at`
- ✅ 完整 Event Sourcing - `booking_history` 記錄所有狀態變更

#### 獨立 Redis 配置

每個服務使用不同的 Redis Database 或獨立 Redis instance：

```yaml
Auth Service:     redis://redis:6379/0  (Session, Token Blacklist)
Restaurant Service: redis://redis:6379/1  (Restaurant Cache)
Booking Service:   redis://redis:6379/2  (Booking Cache)
Spider Service:    redis://redis:6379/3  (Rate Limiting, Distributed Lock)
API Gateway:       redis://redis:6379/4  (Rate Limiting, API Cache)
```

### 4.2 跨服務資料查詢策略

#### 4.2.1 API Composition Pattern
當需要組合多個服務的資料時，由 API Gateway 或 BFF (Backend for Frontend) 負責：

**範例：查詢使用者的預訂記錄（包含餐廳資訊）**
```
1. API Gateway 收到請求 GET /api/v1/users/{userId}/bookings
2. 調用 Booking Service → 取得 booking 列表 (含 restaurant_id)
3. 調用 Restaurant Service → 根據 restaurant_ids 批次查詢餐廳資訊
4. API Gateway 組合資料後回傳
```

#### 4.2.2 CQRS Pattern (Command Query Responsibility Segregation)
針對複雜查詢，建立 **Read Model**：

- **寫入端 (Command)**：各微服務寫入自己的資料庫
- **讀取端 (Query)**：透過事件同步到專門的查詢資料庫
- **實作方式**：
  - 使用 Kafka 發送資料變更事件
  - Query Service 訂閱事件並更新 Read Model (可使用 Elasticsearch)
  - 複雜查詢直接從 Read Model 讀取

**範例架構：**
```
Restaurant Service (寫) → Kafka (restaurant-events) → Query Service → Elasticsearch (讀)
Booking Service (寫)    → Kafka (booking-events)    → Query Service → Elasticsearch (讀)
```

### 4.3 資料一致性處理

#### 4.3.1 Saga Pattern (分散式交易)
使用 **Choreography-based Saga** 處理跨服務交易：

**範例：使用者建立預訂流程**
```
1. Booking Service 建立預訂 (status: pending)
   ├─ 成功 → 發送事件: BookingCreated
   └─ 失敗 → 回傳錯誤

2. Restaurant Service 監聽 BookingCreated
   ├─ 檢查餐廳可用性與容量
   ├─ 成功 → 發送事件: RestaurantConfirmed
   └─ 失敗 → 發送事件: RestaurantRejected

3. Booking Service 監聽 RestaurantConfirmed/Rejected
   ├─ Confirmed → 更新 status: confirmed, 發送 BookingConfirmed
   └─ Rejected → 補償交易: 取消預訂, 發送 BookingCancelled

4. Mail Service 監聽 BookingConfirmed
   └─ 發送確認信給使用者
```

**補償交易 (Compensating Transaction)**：
- 每個步驟必須設計對應的回滾操作
- 使用 Outbox Pattern 確保事件發送可靠性

#### 4.3.2 Eventual Consistency (最終一致性)
- 接受短暫的資料不一致
- 透過事件驅動最終達成一致
- 適用場景：瀏覽量、評論數等非關鍵資料

### 4.4 資料庫技術細節

#### 4.4.1 PostgreSQL 設計規範
- **Schema 設計**：第三正規化 (3NF)
- **主鍵策略**：使用 UUID v4 (分散式友善)
- **軟刪除**：deleted_at TIMESTAMP NULL
- **Audit 欄位**：created_at, updated_at, created_by, updated_by
- **索引策略**：
  - B-tree index：一般查詢
  - GIN index：JSONB、全文檢索
  - Partial index：WHERE deleted_at IS NULL
  - Covering index：避免回表查詢

#### 4.4.2 Migration 管理

✅ **已完成 (2025-11-20)**

```bash
# 每個服務獨立的 migration 目錄
migrations/
├── auth/                                      # ✅ 已完成
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_refresh_tokens_table.up.sql
│   └── 000002_create_refresh_tokens_table.down.sql
├── restaurant/                                # ✅ 已完成
│   ├── 000001_create_restaurants_table.up.sql
│   ├── 000001_create_restaurants_table.down.sql
│   ├── 000002_create_user_favorites_table.up.sql
│   └── 000002_create_user_favorites_table.down.sql
├── booking/                                   # ✅ 已完成
│   ├── 000001_create_bookings_table.up.sql
│   ├── 000001_create_bookings_table.down.sql
│   ├── 000002_create_booking_history_table.up.sql
│   └── 000002_create_booking_history_table.down.sql
├── spider/                                    # ✅ 已完成
│   ├── 000001_create_crawl_jobs_table.up.sql
│   ├── 000001_create_crawl_jobs_table.down.sql
│   ├── 000002_create_crawl_results_table.up.sql
│   └── 000002_create_crawl_results_table.down.sql
├── mail/                                      # ✅ 已完成
│   ├── 000001_create_email_queue_table.up.sql
│   ├── 000001_create_email_queue_table.down.sql
│   ├── 000002_create_email_logs_table.up.sql
│   └── 000002_create_email_logs_table.down.sql
├── MIGRATIONS_SUMMARY.md                     # 完整文檔
└── MIGRATION_EXECUTION_REPORT.md             # 執行報告
```

**工具**：`golang-migrate/migrate`

**執行 Migrations**：
```bash
# 使用自動化腳本（推薦）
./scripts/run_migrations.sh

# 或手動執行
migrate -path migrations/auth -database "postgresql://postgres:postgres@localhost:15432/auth_db?sslmode=disable" up
migrate -path migrations/restaurant -database "postgresql://postgres:postgres@localhost:5433/restaurant_db?sslmode=disable" up
migrate -path migrations/booking -database "postgresql://postgres:postgres@localhost:5434/booking_db?sslmode=disable" up
migrate -path migrations/spider -database "postgresql://postgres:postgres@localhost:5435/spider_db?sslmode=disable" up
migrate -path migrations/mail -database "postgresql://postgres:postgres@localhost:5436/mail_db?sslmode=disable" up
```

**重要提醒**：Auth DB 端口為 **15432**（非標準 5432）

#### 4.4.3 讀寫分離 (可選)
針對讀取量大的服務（如 Restaurant Service）：
- **Master**：處理所有寫入
- **Slave (Read Replicas)**：處理查詢
- **實作方式**：
  - GORM 支援多個 DB 連線
  - 寫入使用 `db.Master()`
  - 查詢使用 `db.Slave()`

### 4.5 Cache 策略

#### 4.5.1 快取層級
```
L1: Application Memory Cache (sync.Map, go-cache)
     ↓ miss
L2: Redis Cache (分散式快取)
     ↓ miss
L3: Database (PostgreSQL)
```

#### 4.5.2 快取模式

**Cache-Aside Pattern** (最常用)
```go
// 讀取流程
data := cache.Get(key)
if data == nil {
    data = db.Query()
    cache.Set(key, data, ttl)
}
return data

// 更新流程
db.Update(data)
cache.Delete(key)  // 刪除快取，下次讀取時重建
```

**Write-Through Pattern** (強一致性)
```go
cache.Set(key, data)
db.Update(data)
```

#### 4.5.3 快取策略細節
- **TTL 設定**：
  - 熱門餐廳資訊：60 分鐘
  - 搜尋結果：15 分鐘
  - 使用者 Session：30 分鐘
  - Rate Limit Counter：1 分鐘
- **Cache Stampede 防護**：使用 `singleflight` 避免同時查詢 DB
- **快取預熱**：系統啟動時預載熱門資料
- **快取淘汰**：LRU (Least Recently Used)

#### 4.5.4 分散式鎖 (Redlock)
```go
// 使用場景：爬蟲去重、庫存扣減
lock := redislock.Obtain(ctx, "lock:crawl:tabelog:tokyo", 30*time.Second)
if lock != nil {
    defer lock.Release()
    // 執行爬蟲任務
}
```

---

## 5. 訊息佇列架構

### 5.1 Kafka 使用場景
- **爬蟲結果處理**
  - Topic: `spider-results`
  - Partition 策略：按餐廳來源分區
  - Consumer Group：資料處理、索引更新

- **事件驅動架構**
  - Topic: `restaurant-events` (新增、更新、刪除)
  - Topic: `booking-events` (預訂成功、取消)
  - Topic: `user-events` (註冊、登入)

### 5.2 訊息處理保證
- At-least-once delivery
- Idempotent consumer 設計
- Dead Letter Queue (DLQ) 處理失敗訊息

---

## 6. 爬蟲服務設計

### 6.1 爬蟲架構
- **並發控制**
  - Worker Pool Pattern (固定數量 goroutines)
  - Rate Limiting (避免被封鎖)
  - Context timeout 控制 (暫定每個請求 30s timeout)

- **爬蟲策略**
  - 使用 colly
  - User-Agent 輪替
  - Proxy 輪替 (可選)
  - 增量爬取 (只爬新增或更新的餐廳)

- **資料回傳**
  - 透過 gRPC streaming 回傳結果到主服務
  - 或發送到 Kafka topic

### 6.2 錯誤處理與重試
- Exponential backoff 重試機制
- 最多重試 3 次
- Circuit Breaker Pattern (防止雪崩)

---

## 7. API 設計

### 7.1 RESTful API 規範
- **版本控制**: `/api/v1/...`
- **HTTP Methods**:
  - GET: 查詢
  - POST: 新增
  - PUT/PATCH: 更新
  - DELETE: 刪除

- **統一回應格式**
```json
{
  "success": true,
  "data": {},
  "error": null,
  "meta": {
    "timestamp": "2025-11-17T10:00:00Z",
    "request_id": "uuid"
  }
}
```

### 7.2 API 文檔
- 使用 Swagger/OpenAPI 3.0
- 自動生成文檔 (swaggo/swag)
- 提供 Postman Collection

### 7.3 API 安全與限流
- JWT Authentication (Auth Service 簽發)
- API Key for third-party integrations
- Rate Limiting:
  - 每個 IP: 100 requests/min
  - 已認證用戶: 1000 requests/min
- CORS 設定

---

## 8. 認證與授權

### 8.1 Authentication
- **JWT (JSON Web Token)**
  - Access Token (15 分鐘有效)
  - Refresh Token (7 天有效，存在 Redis)
  - Token 黑名單機制 (登出時加入)

### 8.2 Authorization
- RBAC (Role-Based Access Control)
- Roles: Admin, User, Guest
- Permissions 檢查 middleware

### 8.3 安全措施
- 密碼使用 bcrypt hash (cost=12)
- HTTPS/TLS 1.3
- SQL Injection 防護 (使用 prepared statements)
- XSS 防護 (輸入驗證、輸出編碼)
- CSRF Token
- Secrets 管理 (使用 Vault 或 AWS Secrets Manager)

---

## 9. 測試策略

### 9.1 測試層級
- **Unit Tests**
  - 使用 Go 標準 testing package
  - Mock: testify/mock 或 gomock
  - 覆蓋率目標: 80%+

- **Integration Tests**
  - 測試服務間整合 (gRPC、Kafka)
  - 使用 testcontainers-go 啟動測試資料庫/Redis

- **E2E Tests**
  - API 端對端測試
  - 使用 httptest

### 9.2 測試工具
- Testing Framework: Go testing
- Assertion: testify/assert
- Mock: testify/mock
- HTTP Testing: httptest
- Database Testing: testcontainers-go

---

## 10. 監控與可觀測性

### 10.1 Metrics (Prometheus)
- **應用層指標**
  - HTTP 請求數、延遲、錯誤率
  - gRPC 調用統計
  - Database query 時間
  - Cache hit/miss rate
  - Kafka 消費 lag

- **系統層指標**
  - CPU、Memory、Disk I/O
  - Goroutine 數量
  - GC 統計

### 10.2 Visualization (Grafana)
- 預建 Dashboard
- 告警規則設定 (延遲 > 1s、錯誤率 > 5%)

### 10.3 Logging (OpenTelemetry)
- **結構化日誌**
  - 使用 zap 或 zerolog
  - JSON 格式
  - 包含 trace_id, span_id (分散式追蹤)

- **日誌等級**
  - DEBUG: 開發環境
  - INFO: 正常流程
  - WARN: 潛在問題
  - ERROR: 錯誤但可恢復
  - FATAL: 致命錯誤

### 10.4 Distributed Tracing
- 使用 OpenTelemetry + Jaeger
- 追蹤請求在微服務間的流轉
- 效能瓶頸分析

---

## 11. 錯誤處理

### 11.1 統一錯誤處理
- 自定義錯誤類型
- Error Code 機制
- 錯誤包裝 (使用 pkg/errors 或 Go 1.13+ errors)

### 11.2 彈性設計
- **Circuit Breaker**
  - 使用 sony/gobreaker
  - 失敗率 > 50% 時熔斷
  - 半開狀態測試恢復

- **重試機制**
  - Exponential Backoff (2^n * 100ms)
  - Jitter (避免同時重試)
  - 最大重試次數: 3

- **Timeout 控制**
  - HTTP Request: 10s
  - gRPC Call: 5s
  - Database Query: 3s
  - Context 傳遞 timeout

---

## 12. 效能優化

### 12.1 Database 優化
- Connection Pool 設定
  - MaxOpenConns: 100
  - MaxIdleConns: 10
  - ConnMaxLifetime: 1 hour

- Query 優化
  - 使用 EXPLAIN ANALYZE
  - N+1 query 問題解決 (Eager Loading)
  - Batch Insert/Update

### 12.2 Cache 優化
- 多層 Cache (Memory → Redis → DB)
- Cache Stampede 防護 (Singleflight)
- Cache Preloading

### 12.3 並發優化
- Goroutine Pool (避免無限制建立)
- Channel Buffering
- Sync.Pool 重用物件

---

## 13. DevOps 與部署

### 13.1 容器化 (Docker)
- Multi-stage build (減少 image 大小)
- 使用 alpine base image
- 每個服務獨立 Dockerfile

### 13.2 編排 (Kubernetes)
- Deployment、Service、Ingress 配置
- ConfigMap 管理配置
- Secret 管理敏感資料
- HPA (Horizontal Pod Autoscaler)
- Liveness & Readiness Probes

### 13.3 CI/CD (GitHub Actions)
- **CI Pipeline**
  - Lint (golangci-lint)
  - Unit Tests
  - Integration Tests
  - Build Docker Image
  - Security Scan (Trivy)

- **CD Pipeline**
  - 自動部署到 dev/test 環境
  - 手動 approval 到 staging/production
  - Rollback ���制

### 13.4 環境管理
- **環境分離**
  - Development
  - Testing
  - Staging
  - Production

- **配置管理**
  - 環境變數 (12-factor app)
  - .env 檔案 (本地開發)
  - ConfigMap/Secret (Kubernetes)
  - Vault (生產環境敏感資料)

### 13.5 健康檢查
- `/health` endpoint (整體健康)
- `/readiness` endpoint (是否可接受流量)
- 檢查項目: DB 連線、Redis 連線、Kafka 連線

---

## 14. 開發工作流程 (Git Flow)

### 14.1 分支策略
- **main**: 生產環境程式碼，只接受來自 release 或 hotfix 的 merge
- **develop**: 開發主分支
- **feature/***: 功能開發 (從 develop 分支)
- **fix/***: Bug 修復 (從 develop 分支)
- **hotfix/***: 緊急修復 (從 main 分支)
- **release/***: 發布準備 (從 develop 分支)

### 14.2 版本管理
- 遵循 Semantic Versioning (v1.2.3)
- 自動產生 CHANGELOG
- Git Tag 標記版本

### 14.3 Code Review
- Pull Request 必須經過至少 1 人 review
- 自動檢查: Lint、Tests、Coverage

---

## 15. 文檔管理

### 15.1 Architecture Decision Records (ADR)
- 記錄重要架構決策
- 格式: 背景、決策、後果

### 15.2 技術文檔
- API 文檔 (Swagger)
- 資料庫 Schema 文檔
- 部署文檔
- 開發環境設定文檔

### 15.3 註解規範
- 公開函數必須有註解
- 複雜邏輯需要說明
- TODO/FIXME 標記待處理項目

---

## 16. 資料庫 Schema 範例（按服務分離）

### 16.1 Auth Service Database (`auth_db`) ✅

**實現狀態**：已完成 (2025-11-20)
**Migration 版本**：v2
**資料表**：users, refresh_tokens
**連接端口**：15432 ⚠️

#### Users Table ✅
```sql
-- Database: auth_db
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    username VARCHAR(50) NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    is_active BOOLEAN DEFAULT TRUE,
    email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(role) WHERE deleted_at IS NULL;
```

#### Refresh Tokens Table
```sql
-- Database: auth_db
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,  -- 注意：不使用 REFERENCES，避免跨 DB 外鍵
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    revoked_at TIMESTAMP NULL
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
```

### 16.2 Restaurant Service Database (`restaurant_db`) ✅

**實現狀態**：已完成 (2025-11-20)
**Migration 版本**：v2
**資料表**：restaurants, user_favorites

#### Restaurants Table ✅
```sql
-- Database: restaurant_db
CREATE TABLE restaurants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    source VARCHAR(50) NOT NULL, -- 'tabelog', 'google', etc.
    external_id VARCHAR(255) NOT NULL,
    address TEXT,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    rating DECIMAL(3, 2),
    price_range VARCHAR(10),
    cuisine_type VARCHAR(50),
    phone VARCHAR(20),
    website VARCHAR(500),
    opening_hours JSONB,
    metadata JSONB,
    view_count BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

CREATE UNIQUE INDEX idx_restaurants_source_external_id
    ON restaurants(source, external_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_restaurants_location
    ON restaurants USING GIST(ll_to_earth(latitude, longitude));
CREATE INDEX idx_restaurants_cuisine ON restaurants(cuisine_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_restaurants_rating ON restaurants(rating DESC) WHERE deleted_at IS NULL;
```

#### User Favorites Table ✅ **新增 (2025-11-20)**
```sql
-- Database: restaurant_db
-- 使用者收藏餐廳功能（取代 reviews 表）
CREATE TABLE user_favorites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,  -- 來自 Auth Service，不使用外鍵
    restaurant_id UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    notes TEXT,  -- 使用者私人筆記
    tags VARCHAR(255)[],  -- 使用者自定義標籤
    visit_count INT DEFAULT 0,  -- 造訪次數
    last_visited_at TIMESTAMP,  -- 最後造訪時間
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_user_favorites_user_id ON user_favorites(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_user_favorites_restaurant_id ON user_favorites(restaurant_id) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_user_favorites_unique ON user_favorites(user_id, restaurant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_user_favorites_tags ON user_favorites USING GIN(tags);
```

**設計說明**：
- ❌ **移除 reviews 表** - 評論、評分等資料完全來自外部（Google Maps, Tabelog, Instagram）
- ✅ **新增 user_favorites 表** - 使用者只能收藏餐廳、添加私人筆記、標籤
- 使用者**無法編輯餐廳資料**，僅能查詢和收藏

### 16.3 Booking Service Database (`booking_db`) ✅

**實現狀態**：已完成 (2025-11-20)
**Migration 版本**：v2
**資料表**：bookings, booking_history

#### Bookings Table ✅
```sql
-- Database: booking_db
CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,        -- 來自 Auth Service
    restaurant_id UUID NOT NULL,  -- 來自 Restaurant Service
    booking_date TIMESTAMP NOT NULL,
    party_size INT NOT NULL CHECK (party_size > 0),
    status VARCHAR(20) DEFAULT 'pending', -- pending, confirmed, cancelled, completed
    external_booking_id VARCHAR(255),  -- OpenTable 的預訂 ID
    external_service VARCHAR(50),  -- 外部服務名稱 (opentable, tabelog)
    last_synced_at TIMESTAMP,  -- 最後同步時間
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_restaurant_id ON bookings(restaurant_id);
CREATE INDEX idx_bookings_date ON bookings(booking_date);
CREATE INDEX idx_bookings_status ON bookings(status);
```

#### Booking History Table (Event Sourcing)
```sql
-- Database: booking_db
CREATE TABLE booking_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id),
    status VARCHAR(20) NOT NULL,
    changed_by UUID,  -- user_id
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_booking_history_booking_id ON booking_history(booking_id);
CREATE INDEX idx_booking_history_created_at ON booking_history(created_at DESC);
```

### 16.4 Spider Service Database (`spider_db`) ✅

**實現狀態**：已完成 (2025-11-20)
**Migration 版本**：v2
**資料表**：crawl_jobs, crawl_results

#### Crawl Jobs Table ✅
```sql
-- Database: spider_db
CREATE TABLE crawl_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source VARCHAR(50) NOT NULL,  -- 'tabelog', 'google_maps', 'instagram', etc.
    region VARCHAR(100),
    status VARCHAR(20) DEFAULT 'pending',  -- pending, running, completed, failed
    total_pages INT,
    completed_pages INT DEFAULT 0,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_crawl_jobs_status ON crawl_jobs(status);
CREATE INDEX idx_crawl_jobs_source ON crawl_jobs(source);
```

#### Crawl Results Table
```sql
-- Database: spider_db
CREATE TABLE crawl_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES crawl_jobs(id),
    external_id VARCHAR(255) NOT NULL,
    source VARCHAR(50) NOT NULL,
    raw_data JSONB NOT NULL,
    processed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_crawl_results_job_id ON crawl_results(job_id);
CREATE INDEX idx_crawl_results_processed ON crawl_results(processed) WHERE processed = FALSE;
CREATE UNIQUE INDEX idx_crawl_results_source_external_id
    ON crawl_results(source, external_id);
```

### 16.5 Mail Service Database (`mail_db`) ✅

**實現狀態**：已完成 (2025-11-20)
**Migration 版本**：v2
**資料表**：email_queue, email_logs

#### Email Queue Table ✅
```sql
-- Database: mail_db
CREATE TABLE email_queue (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipient_email VARCHAR(255) NOT NULL,
    recipient_name VARCHAR(100),
    subject VARCHAR(500) NOT NULL,
    body TEXT NOT NULL,
    template_name VARCHAR(100),
    template_data JSONB,
    priority INT DEFAULT 5,  -- 1 (highest) to 10 (lowest)
    status VARCHAR(20) DEFAULT 'pending',  -- pending, sent, failed
    retry_count INT DEFAULT 0,
    max_retries INT DEFAULT 3,
    scheduled_at TIMESTAMP DEFAULT NOW(),
    sent_at TIMESTAMP,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_email_queue_status ON email_queue(status, scheduled_at);
CREATE INDEX idx_email_queue_priority ON email_queue(priority DESC, created_at);
```

### 16.6 跨服務查詢說明

**重要原則：**
1. ❌ **不使用跨資料庫的外鍵約束** (FOREIGN KEY)
2. ✅ 僅儲存關聯 ID (如 user_id, restaurant_id)
3. ✅ 資料完整性透過應用層或事件驅動機制保證
4. ✅ 使用 Saga Pattern 處理分散式交易

**範例：查詢使用者的預訂（包含餐廳資訊）**
```go
// 1. Booking Service 查詢預訂
bookings := bookingRepo.FindByUserID(userID)

// 2. 取得所有 restaurant_ids
restaurantIDs := extractRestaurantIDs(bookings)

// 3. gRPC 調用 Restaurant Service 批次查詢
restaurants := restaurantClient.GetByIDs(restaurantIDs)

// 4. 組合資料
return combineBookingsWithRestaurants(bookings, restaurants)
```

---

## 17. 專案目錄結構（微服務架構）

```
tabelogov2/
├── cmd/                          # 各微服務的主程式入口
│   ├── api-gateway/              # API Gateway
│   │   └── main.go
│   ├── auth-service/             # Auth 微服務
│   │   └── main.go
│   ├── booking-service/
│   │   └── main.go
│   ├── map-service/
│   │   └── main.go
│   ├── spider-service/
│   │   └── main.go
│   ├── mail-service/
│   │   └── main.go
│   └── restaurant-service/
│       └── main.go
│
├── internal/                     # 按服務分離的內部程式碼
│   ├── auth/                     # Auth Service 專屬
│   │   ├── domain/               # Domain 層
│   │   │   ├── user/
│   │   │   └── token/
│   │   ├── application/          # Application 層 (Use Cases)
│   │   │   ├── service/
│   │   │   └── dto/
│   │   ├── infrastructure/       # Infrastructure 層
│   │   │   ├── persistence/      # auth_db 連線與 Repository
│   │   │   ├── cache/            # Redis 操作
│   │   │   └── messaging/        # Kafka Producer
│   │   └── presentation/         # Presentation 層
│   │       ├── grpc/             # gRPC handlers
│   │       └── http/             # HTTP handlers (可選)
│   │
│   ├── restaurant/               # Restaurant Service 專屬
│   │   ├── domain/
│   │   │   ├── restaurant/
│   │   │   └── review/
│   │   ├── application/
│   │   ├── infrastructure/
│   │   │   ├── persistence/      # restaurant_db
│   │   │   ├── cache/
│   │   │   ├── messaging/
│   │   │   └── grpc/             # gRPC clients (呼叫其他服務)
│   │   └── presentation/
│   │
│   ├── booking/                  # Booking Service 專屬
│   │   ├── domain/
│   │   ├── application/
│   │   ├── infrastructure/
│   │   │   ├── persistence/      # booking_db
│   │   │   ├── cache/
│   │   │   ├── messaging/
│   │   │   └── external/         # OpenTable API client
│   │   └── presentation/
│   │
│   ├── spider/                   # Spider Service 專屬
│   │   ├── domain/
│   │   ├── application/
│   │   ├── infrastructure/
│   │   │   ├── persistence/      # spider_db
│   │   │   ├── messaging/
│   │   │   └── crawler/          # 爬蟲實作 (colly)
│   │   └── presentation/
│   │
│   ├── mail/                     # Mail Service 專屬
│   │   ├── domain/
│   │   ├── application/
│   │   ├── infrastructure/
│   │   │   ├── persistence/      # mail_db
│   │   │   ├── messaging/
│   │   │   └── smtp/             # SMTP client
│   │   └── presentation/
│   │
│   ├── map/                      # Map Service 專屬
│   │   ├── application/
│   │   ├── infrastructure/
│   │   │   ├── cache/
│   │   │   └── external/         # Google Maps API client
│   │   └── presentation/
│   │
│   └── gateway/                  # API Gateway 專屬
│       ├── router/
│       ├── middleware/
│       └── handler/
│
├── pkg/                          # 跨服務共用套件
│   ├── logger/                   # 統一日誌套件
│   ├── errors/                   # 錯誤處理
│   ├── config/                   # 配置載入
│   ├── middleware/               # 共用 middleware
│   ├── utils/                    # 工具函數
│   ├── jwt/                      # JWT 驗證工具
│   ├── grpc/                     # gRPC 連線管理
│   ├── kafka/                    # Kafka 連線管理
│   └── tracing/                  # OpenTelemetry 追蹤
│
├── api/
│   ├── proto/                    # Protocol Buffers 定義
│   │   ├── auth/
│   │   │   └── v1/
│   │   │       ├── auth.proto
│   │   │       └── auth.pb.go (generated)
│   │   ├── restaurant/
│   │   │   └── v1/
│   │   ├── booking/
│   │   │   └── v1/
│   │   └── common/               # 共用的 proto messages
│   │       └── v1/
│   │           └── common.proto
│   └── openapi/                  # OpenAPI/Swagger 定義
│       ├── auth.yaml
│       ├── restaurant.yaml
│       └── booking.yaml
│
├── migrations/                   # 各服務的 DB migrations
│   ├── auth/
│   │   ├── 000001_create_users_table.up.sql
│   │   └── 000001_create_users_table.down.sql
│   ├── restaurant/
│   │   ├── 000001_create_restaurants_table.up.sql
│   │   └── 000001_create_restaurants_table.down.sql
│   ├── booking/
│   ├── spider/
│   └── mail/
│
├── deployments/
│   ├── docker/
│   │   ├── api-gateway.Dockerfile
│   │   ├── auth-service.Dockerfile
│   │   ├── restaurant-service.Dockerfile
│   │   ├── booking-service.Dockerfile
│   │   ├── spider-service.Dockerfile
│   │   ├── mail-service.Dockerfile
│   │   └── map-service.Dockerfile
│   ├── docker-compose/
│   │   ├── docker-compose.yml        # 本地開發環境
│   │   ├── docker-compose.dev.yml
│   │   └── docker-compose.test.yml
│   └── k8s/                          # Kubernetes manifests
│       ├── base/                     # 基礎配置
│       │   ├── namespace.yaml
│       │   ├── auth-service/
│       │   │   ├── deployment.yaml
│       │   │   ├── service.yaml
│       │   │   └── configmap.yaml
│       │   ├── restaurant-service/
│       │   ├── booking-service/
│       │   └── ...
│       ├── overlays/                 # Kustomize overlays
│       │   ├── dev/
│       │   ├── staging/
│       │   └── production/
│       └── infrastructure/           # 基礎設施
│           ├── postgres.yaml
│           ├── redis.yaml
│           ├── kafka.yaml
│           ├── prometheus.yaml
│           └── grafana.yaml
│
├── scripts/
│   ├── build.sh                      # 建置所有服務
│   ├── run-dev.sh                    # 啟動開發環境
│   ├── generate-proto.sh             # 生成 gRPC code
│   ├── db-migrate.sh                 # 執行 migrations
│   └── k8s-deploy.sh                 # 部署到 Kubernetes
│
├── tests/
│   ├── integration/                  # 整合測試
│   │   ├── auth_test.go
│   │   ├── booking_flow_test.go
│   │   └── restaurant_search_test.go
│   ├── e2e/                          # E2E 測試
│   │   └── api_test.go
│   └── fixtures/                     # 測試數據
│
├── docs/
│   ├── architecture.md               # 本文檔
│   ├── api/                          # API 文檔
│   ├── deployment/                   # 部署文檔
│   ├── development/                  # 開發指南
│   └── adr/                          # Architecture Decision Records
│       ├── 0001-use-microservices.md
│       ├── 0002-database-per-service.md
│       └── 0003-use-kafka-for-events.md
│
├── .github/
│   └── workflows/
│       ├── ci.yml                    # CI pipeline
│       ├── cd-dev.yml                # CD to dev
│       ├── cd-staging.yml            # CD to staging
│       └── cd-production.yml         # CD to production
│
├── go.mod
├── go.sum
├── Makefile                          # 常用命令
├── .env.example                      # 環境變數範例
├── .gitignore
└── README.md
```

### 17.1 目錄結構說明

#### 服務獨立性
- 每個微服務在 `cmd/` 有獨立的 main.go 入口
- 每個微服務在 `internal/` 有獨立的程式碼目錄
- 每個微服務在 `migrations/` 有獨立的資料庫遷移
- 每個微服務在 `deployments/docker/` 有獨立的 Dockerfile

#### DDD 分層（以 Restaurant Service 為例）
```
internal/restaurant/
├── domain/              # 核心業務邏輯
│   ├── restaurant/
│   │   ├── restaurant.go      # Aggregate Root
│   │   ├── repository.go      # Repository Interface
│   │   └── value_object.go    # Value Objects
│   └── review/
├── application/         # Use Cases
│   ├── service/
│   │   ├── create_restaurant.go
│   │   └── search_restaurant.go
│   └── dto/             # Data Transfer Objects
├── infrastructure/      # 外部依賴實作
│   ├── persistence/
│   │   └── postgres_repository.go  # Repository 實作
│   ├── cache/
│   │   └── redis_cache.go
│   └── messaging/
│       └── kafka_producer.go
└── presentation/        # 對外接口
    ├── grpc/
    │   └── handler.go
    └── http/
        └── handler.go
```

---

## 18. 下一步行動計畫

### Phase 1: 本地開發基礎建設 (Week 1-2) ✅
- [x] 專案初始化、目錄結構建立
- [x] Docker、docker-compose 設定
- [x] PostgreSQL、Redis、Kafka 環境建置
- [x] 開發工具設定（Makefile、.env、.gitignore）
- [x] **基礎 Migrations 建立（每個服務的第一個 migration）** ✅ **2025-11-20 完成**
  - [x] Auth Service: users, refresh_tokens
  - [x] Restaurant Service: restaurants, user_favorites
  - [x] Booking Service: bookings, booking_history
  - [x] Spider Service: crawl_jobs, crawl_results
  - [x] Mail Service: email_queue, email_logs
  - [x] 共用 trigger functions (update_updated_at_column)
  - [x] 完整索引策略 (B-tree, GIN, Partial, Composite)
  - [x] Migration 文檔 (MIGRATIONS_SUMMARY.md, MIGRATION_EXECUTION_REPORT.md)
- [x] **共用套件完整實作** ✅ **2025-11-20 完成**
  - [x] pkg/logger - 統一日誌套件 (Zap + Context 支援)
  - [x] pkg/config - 配置載入與管理
  - [x] pkg/errors - 錯誤處理 (HTTP + gRPC)
  - [x] pkg/middleware - HTTP 中間件 (7 個完整中間件)
  - [x] 完整單元測試與文檔

### Phase 2: 核心服務開發 (Week 3-6) ✅ **90% 完成 (2025-12-02)**
- [x] **Auth Service 開發** ✅ **完成 (2025-11-20)**
  - [x] Domain Layer: User Aggregate, Token Model
  - [x] Application Layer: Service with JWT logic
  - [x] Infrastructure Layer: PostgreSQL + Redis repositories
  - [x] gRPC Server 實作 (Login, Register, ValidateToken)
  - [x] HTTP REST API (Swagger 文檔)
  - [x] JWT 簽發與驗證 (Access + Refresh Token)
  - [x] RBAC 權限管理
  - [x] Uber FX 依賴注入
  - [x] Docker 容器化 (獨立 auth_db on port 15432)
  - [x] 完整單元測試
  - [x] Swagger UI 整合 (http://localhost:18080/swagger)
  
- [/] **Frontend 開發** ✅ **60% 完成 (2025-11-30)**
  - [x] 專案初始化 (Next.js 16, TailwindCSS v4, Shadcn/UI)
  - [x] Map-First 介面實作 (@vis.gl/react-google-maps)
  - [x] 進階搜尋 UI (AdvanceSearchForm component)
  - [x] 地圖標記與互動 (GoogleMap component)
  - [x] 響應式設計 (Dark mode, Sidebar navigation)
  - [x] TypeScript 類型定義 (Place, SearchFilters)
  - [x] React Query 狀態管理準備
  - [/] **API 整合** (進行中)
    - [x] Map Service 整合 (useMapSearch hook)
    - [ ] Auth Service 整合 (AuthContext 已建立但未連接)
    - [ ] Restaurant Service 整合 (待實作)
  - [ ] **頁面開發**
    - [x] 主頁 (Map + Search)
    - [/] Login/Register 頁面 (UI 已建立，待 API 連接)
    - [ ] Restaurant Details 頁面
    - [ ] User Profile 頁面
    - [ ] Booking 頁面
  
- [x] **Restaurant Service 開發** ✅ **完成 (2025-12-02)**
  - [x] Domain Layer: Restaurant & Favorite Aggregates, Location Value Object
  - [x] Application Layer: RestaurantService (19 methods)
  - [x] Infrastructure Layer: PostgreSQL + GORM repositories
  - [x] Repository Pattern: RestaurantRepository, FavoriteRepository
  - [x] HTTP REST API (6 endpoints: Create, Get, Search, AddFavorite, GetFavorites, Health)
  - [x] CRUD API: 完整 CRUD 操作 + 去重邏輯
  - [x] 搜尋功能: Search, FindByLocation, FindByCuisineType
  - [x] 收藏功能: Add/Remove favorites, Tags, Visit tracking
  - [x] 外部 ID 整合: (source, external_id) composite unique index
  - [x] Docker 容器化 (Dockerfile + docker-compose)
  - [x] 整合測試腳本 (test-restaurant-service.sh)
  - [x] Uber FX 依賴注入
  - [x] 完整文檔 (IMPLEMENTATION_SUMMARY.md)
  - [ ] gRPC Server (未來增強)
  - [ ] 單元測試 (待實作)
  
- [ ] **API Gateway 實作** ⚠️ **未開始**
  - [ ] 路由設定 (Gin router)
  - [ ] gRPC 轉 HTTP (grpc-gateway)
  - [ ] 認證 Middleware (JWT 驗證)
  - [ ] Rate Limiting
  - [ ] CORS 設定
  - [ ] 請求日誌與監控

### Phase 3: 整合服務與事件驅動 (Week 7-9)
- [ ] Kafka 整合
  - [ ] Producer/Consumer 基礎設定
  - [ ] Event Schema 定義
  - [ ] Saga Pattern 實作
- [ ] Booking Service 開發
  - [ ] 預訂流程實作
  - [ ] OpenTable API 整合
  - [ ] 事件發送（Kafka）
- [ ] Spider Service 開發
  - [ ] Tabelog 爬蟲實作
  - [ ] Worker Pool 並發控制
  - [ ] 結果發送至 Kafka
- [ ] Mail Service 開發
  - [ ] 監聽 Kafka 事件
  - [ ] SMTP 郵件發送
- [x] Map Service 開發 ✅ **Phase 4 完成 (2025-11-29)**
  - [x] Phase 1: 基礎架構建立
    - [x] DDD 分層架構 (Domain, Application, Infrastructure, Interfaces)
    - [x] Uber FX 依賴注入
    - [x] Redis 連接與配置
    - [x] HTTP Server 設置 (Gin)
    - [x] Health Check 端點
  - [x] Phase 2: Quick Search 功能
    - [x] Google Places API (New) 客戶端
    - [x] Redis 緩存層 (Cache-first 策略)
    - [x] Quick Search Use Case (業務邏輯)
    - [x] HTTP Handler 實作
    - [x] 完整錯誤處理與日誌
    - [x] 性能測試通過 (50倍速度提升)
  - [x] Phase 3: Advance Search 功能
    - [x] Text Search API 整合
    - [x] 搜索結果過濾 (min_rating, open_now)
    - [x] 地理位置搜索 (Rectangle location bias)
    - [x] 排序偏好 (DISTANCE, RELEVANCE)
    - [x] 多語言支援 (en, ja, zh-TW)
    - [x] 性能測試通過 (~400ms)
  - [x] Phase 4: 優化與監控
    - [x] Prometheus Metrics 整合
      - [x] HTTP 請求監控 (計數器、延遲直方圖)
      - [x] Google API 調用追蹤
      - [x] Cache 命中率監控
      - [x] /metrics 端點
    - [x] Rate Limiting
      - [x] Redis-backed 分佈式限流
      - [x] Quick Search: 60 requests/min
      - [x] Advance Search: 30 requests/min
      - [x] 429 錯誤響應與 retry_after
    - [ ] Enhanced Health Check (未來優化)
    - [ ] 生產環境配置 (未來優化)

### Phase 4: 監控、測試與優化 (Week 10-11)
- [ ] 可觀測性建置
  - [ ] Prometheus Metrics 埋點
  - [ ] Grafana Dashboard 建立
  - [ ] OpenTelemetry 分散式追蹤
  - [ ] Jaeger 整合
- [ ] 測試完善
  - [ ] 單元測試（目標 80%+ 覆蓋率）
  - [ ] 整合測試
  - [ ] E2E 測試
- [ ] 效能優化
  - [ ] Cache 策略優化
  - [ ] Database Query 優化
  - [ ] 負載測試

### Phase 5: CI/CD 與文檔 (Week 12)
- [ ] CI/CD Pipeline 建置
  - [ ] GitHub Actions CI 設定
  - [ ] 自動化測試流程
  - [ ] Docker Image 建置
  - [ ] 部署自動化
- [ ] 文檔完善
  - [ ] API 文檔（Swagger）
  - [ ] 架構決策記錄（ADR）
  - [ ] 部署文檔
  - [ ] Demo 準備

---

## 19. 微服務架構圖

```
                                   ┌─────────────────┐
                                   │   Load Balancer │
                                   └────────┬────────┘
                                            │
                    ┌───────────────────────┼───────────────────────┐
                    │                       │                       │
            ┌───────▼────────┐      ┌──────▼───────┐       ┌──────▼───────┐
            │  API Gateway   │      │  API Gateway │       │  API Gateway │
            │  (Multiple)    │      │  (Multiple)  │       │  (Multiple)  │
            └───────┬────────┘      └──────┬───────┘       └──────┬───────┘
                    │                      │                       │
                    └──────────────────────┴───────────────────────┘
                                          │
          ┌───────────┬──────────┬────────┴────────┬─────────┬──────────┐
          │           │          │                 │         │          │
    ┌─────▼────┐ ┌───▼────┐ ┌──▼─────┐   ┌───────▼──┐  ┌───▼────┐ ┌──▼─────┐
    │   Auth   │ │Restaurant│ │Booking│   │  Spider  │  │  Mail  │ │  Map   │
    │ Service  │ │ Service │ │Service│   │ Service  │  │Service │ │Service │
    └────┬─────┘ └────┬────┘ └───┬───┘   └────┬─────┘  └───┬────┘ └───┬────┘
         │            │          │            │            │          │
    ┌────▼─────┐ ┌───▼─────┐ ┌──▼──────┐ ┌──▼───────┐ ┌──▼─────┐    │
    │ auth_db  │ │restaurant│ │booking │ │spider_db │ │mail_db │    │
    │(Postgres)│ │   _db    │ │  _db   │ │(Postgres)│ │(Postgres)   │
    └──────────┘ │(Postgres)│ │(Postgres└──────────┘ └────────┘    │
                 └──────────┘ └─────────┘                            │
                                                                      │
         ┌────────────────────────────────────────────────────────────┘
         │
    ┌────▼─────────┐
    │ Google Maps  │
    │     API      │
    └──────────────┘

    ┌──────────────────────────────────────────────────────────────┐
    │                        Redis Cluster                          │
    │  DB0: Auth Cache    DB1: Restaurant Cache   DB2: Booking...  │
    └──────────────────────────────────────────────────────────────┘

    ┌──────────────────────────────────────────────────────────────┐
    │                        Kafka Cluster                          │
    │  Topics: user-events, restaurant-events, booking-events...   │
    └──────────────────────────────────────────────────────────────┘

    ┌──────────────────────────────────────────────────────────────┐
    │                    Monitoring Stack                           │
    │  Prometheus + Grafana + Jaeger + OpenTelemetry               │
    └──────────────────────────────────────────────────────────────┘
```

### 19.1 服務通訊模式

#### 同步通訊 (gRPC)
```
API Gateway → Auth Service (驗證 Token)
API Gateway → Restaurant Service (查詢餐廳)
Booking Service → Restaurant Service (檢查餐廳可用性)
```

#### 非同步通訊 (Kafka)
```
Auth Service → Kafka (user-events) → Mail Service (發送歡迎信)
Booking Service → Kafka (booking-events) → Mail Service (發送確認信)
Spider Service → Kafka (spider-results) → Restaurant Service (更新餐廳資料)
```

---

## 20. 履歷展示重點總結

此專案完整展示的技術能力：

✅ **微服務架構能力**
- Database per Service 原則實踐
- 服務間通訊設計 (gRPC + Kafka)
- API Gateway 模式
- Service Discovery

✅ **分散式系統經驗**
- Saga Pattern 分散式交易
- Eventual Consistency 最終一致性
- API Composition Pattern
- CQRS 讀寫分離

✅ **資料庫設計**
- 每個服務獨立資料庫
- 無跨資料庫外鍵約束
- 索引優化策略
- Migration 管理

✅ **DDD 領域驅動設計**
- 清晰的分層架構
- Aggregate Root 設計
- Repository Pattern
- Domain Events

✅ **效能優化**
- 多層 Cache 策略
- Connection Pooling
- 並發控制 (Goroutines)
- 批次查詢優化

✅ **DevOps 能力**
- Docker 容器化
- Kubernetes 編排
- CI/CD 自動化
- 多環境管理

✅ **可觀測性**
- Prometheus 監控
- Grafana 視覺化
- OpenTelemetry 分散式追蹤
- 結構化日誌

✅ **測試能力**
- 單元測試 (80%+ 覆蓋率)
- 整合測試
- E2E 測試
- Testcontainers

✅ **安全性**
- JWT + Refresh Token
- RBAC 授權
- API Rate Limiting
- Secrets 管理

---

## 21. 參考資料

### 微服務架構
- [Microservices Patterns - Chris Richardson](https://microservices.io/patterns/index.html)
- [Building Microservices - Sam Newman](https://samnewman.io/books/building_microservices_2nd_edition/)
- [Database per Service Pattern](https://microservices.io/patterns/data/database-per-service.html)
- [Saga Pattern](https://microservices.io/patterns/data/saga.html)

### Go 語言
- [Go 官方文檔](https://go.dev/doc/)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [Effective Go](https://go.dev/doc/effective_go)

### 分散式系統
- [gRPC Go Quick Start](https://grpc.io/docs/languages/go/quickstart/)
- [Apache Kafka Documentation](https://kafka.apache.org/documentation/)
- [Domain-Driven Design - Eric Evans](https://www.domainlanguage.com/ddd/)

### DevOps & 部署
- [12-Factor App](https://12factor.net/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)

### 可觀測性
- [Prometheus Best Practices](https://prometheus.io/docs/practices/)
- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)
- [Grafana Documentation](https://grafana.com/docs/)

---

## 22. 更新記錄 (Change Log)

### 2025-11-20 - Migration 實現完成 ✅

**完成項目：Phase 1 - 基礎 Migrations 建立**

#### 實現內容
- ✅ 建立 5 個微服務的資料庫 migrations（10 個資料表）
  - Auth Service: `users`, `refresh_tokens`
  - Restaurant Service: `restaurants`, `user_favorites`
  - Booking Service: `bookings`, `booking_history`
  - Spider Service: `crawl_jobs`, `crawl_results`
  - Mail Service: `email_queue`, `email_logs`

#### 設計調整
1. **Restaurant Service**
   - ❌ 移除 `reviews` 表 - 評論資料完全來自外部（Google Maps, Tabelog, Instagram）
   - ✅ 新增 `user_favorites` 表 - 使用者僅能收藏、查詢，無法編輯餐廳資料

2. **Booking Service**
   - ✅ 新增外部 API 同步支援
   - 新增欄位：`external_service`, `last_synced_at`
   - 支援與 OpenTable/Tabelog 等外部服務同步

3. **索引優化**
   - 移除 WHERE 子句中的 `NOW()` 函數（PostgreSQL immutability 限制）
   - 改為在應用層過濾或調整索引策略

#### 技術細節
- 所有表格使用 UUID v4 主鍵
- 實現軟刪除（`deleted_at`）
- 自動 `updated_at` 觸發器
- 完整索引策略（B-tree, GIN, Partial, Composite）
- Event Sourcing（`booking_history`）

#### 文檔
- [migrations/MIGRATIONS_SUMMARY.md](migrations/MIGRATIONS_SUMMARY.md) - 完整設計文檔
- [migrations/MIGRATION_EXECUTION_REPORT.md](migrations/MIGRATION_EXECUTION_REPORT.md) - 執行報告
- [scripts/run_migrations.sh](scripts/run_migrations.sh) - 自動化執行腳本

#### 環境配置
- Auth DB: Port **15432** ⚠️（非標準 5432）
- Restaurant DB: Port 5433
- Booking DB: Port 5434
- Spider DB: Port 5435
- Mail DB: Port 5436

### 2025-11-20 - Middleware 套件完成 ✅

**完成項目：Phase 1 - HTTP 中間件完整實作**

#### 實現內容
- ✅ **認證中間件 (Authentication)**
  - JWT Token 驗證
  - Bearer Token 解析
  - 用戶角色檢查 (RBAC)
  - 跳過路徑配置
  - Context 中的用戶資訊提取

- ✅ **速率限制中間件 (Rate Limiting)**
  - Redis 分散式速率限制（生產環境）
  - 記憶體內速率限制（開發環境）
  - 滑動視窗演算法
  - 按 IP 或用戶 ID 限流
  - 自動設定速率限制標頭

- ✅ **請求 ID 中間件 (Request ID)**
  - 自動生成 UUID
  - 支援現有 Request ID
  - Context 中的 Request ID 管理
  - 回應標頭設定

- ✅ **日誌中間件 (Logger)**
  - Uber Zap 結構化日誌
  - 請求/回應詳細資訊
  - 延遲時間追蹤
  - 錯誤日誌

- ✅ **錯誤處理中間件 (Error Handler)**
  - AppError 類型識別
  - HTTP 狀態碼自動映射
  - 統一錯誤回應格式
  - 詳細錯誤資訊

- ✅ **恢復中間件 (Recovery)**
  - Panic 捕獲
  - 錯誤日誌記錄
  - 優雅的錯誤回應

- ✅ **CORS 中間件 (CORS)**
  - 跨來源資源共享配置
  - OPTIONS 請求處理
  - 自定義標頭和方法

#### 新增依賴
```
github.com/gin-gonic/gin v1.11.0
github.com/golang-jwt/jwt/v5 v5.3.0
github.com/redis/go-redis/v9 v9.17.0
github.com/google/uuid v1.6.0
github.com/stretchr/testify v1.11.1
```

#### 測試與文檔
- ✅ 所有中間件包含完整單元測試
- ✅ 測試案例涵蓋正常和異常情況
- ✅ 完整使用文檔 (MIDDLEWARE.md)
- ✅ 程式碼範例與最佳實踐
- ✅ 效能考量與安全性建議

#### 技術細節
- JWT 認證支援自定義 Claims
- 速率限制使用 Redis 滑動視窗演算法
- Request ID 使用 UUID v4
- 日誌使用 Zap 高效能結構化日誌
- 錯誤處理與 pkg/errors 完美整合

#### 文檔
- [pkg/middleware/MIDDLEWARE.md](pkg/middleware/MIDDLEWARE.md) - 完整使用文檔
- 包含 7 個中間件的詳細說明
- 完整的使用範例和最佳實踐
- 中間件鏈配置範例
- 效能優化建議

#### 下一步
- [ ] 開始 Phase 2：核心服務開發
  - Auth Service (JWT 簽發、使用者管理)
  - Restaurant Service (餐廳 CRUD、搜尋)
  - API Gateway (路由、gRPC 轉 HTTP)

---

### 2025-11-20 - 共用套件實現完成 ✅

**完成項目：Phase 1 - 共用套件基礎實作**

#### 實現內容
- ✅ `pkg/logger` - 統一日誌套件
  - 基於 `go.uber.org/zap` 高效能日誌
  - 支援多個日誌等級 (debug, info, warn, error, fatal)
  - Context 支援（日誌追蹤）
  - Uber FX 依賴注入整合

- ✅ `pkg/config` - 配置載入與管理
  - 從環境變數載入配置
  - 支援環境變數前綴（多服務配置）
  - 完整的配置驗證
  - 型別安全的配置存取
  - Uber FX 依賴注入整合

- ✅ `pkg/errors` - 統一錯誤處理
  - 統一的錯誤碼系統
  - HTTP 狀態碼自動映射
  - gRPC 錯誤支援（ToGRPCError, FromGRPCError）
  - 錯誤包裝與追蹤
  - 詳細資訊附加

#### 新增功能
1. **logger 套件增強**
   - ✅ Context 支援 (`WithContext`, `FromContext`, `WithFields`)
   - ✅ 線程安全的 logger 管理
   - ✅ 自動 fallback 機制

2. **errors 套件增強**
   - ✅ gRPC 錯誤轉換 (`ToGRPCError`, `FromGRPCError`)
   - ✅ 錯誤碼與 gRPC codes 映射
   - ✅ HTTP 狀態碼自動推導

3. **測試與文檔**
   - ✅ 所有套件包含完整單元測試
   - ✅ 測試覆蓋率 > 80%
   - ✅ 完整使用文檔 (SHARED_PACKAGES.md)
   - ✅ 程式碼範例與最佳實踐

#### 技術細節
- 使用 Uber Zap（比 logrus 快 4-10x）
- 結構化日誌 (JSON 格式)
- 環境變數配置（遵循 12-Factor App）
- 錯誤追蹤與包裝（保留 stack trace）

#### 文檔
- [pkg/SHARED_PACKAGES.md](pkg/SHARED_PACKAGES.md) - 完整使用文檔
- [pkg/logger/logger_test.go](pkg/logger/logger_test.go) - 單元測試
- [pkg/config/config_test.go](pkg/config/config_test.go) - 單元測試
- [pkg/errors/errors_test.go](pkg/errors/errors_test.go) - 單元測試

#### 下一步
- [ ] pkg/middleware 實作（認證、日誌、錯誤處理 middleware）
- [ ] 開始 Phase 2：核心服務開發

---

## 附錄：關鍵決策記錄 (ADR)

### ADR-001: 採用 Database per Service 模式

**狀態**: 已採用

**背景**：
微服務架構中，服務間的資料存取方式是關鍵決策。可選方案包括：
1. 共享資料庫
2. Database per Service
3. 混合模式

**決策**：
採用 Database per Service 模式，每個微服務擁有獨立的資料庫實例。

**理由**：
- ✅ 真正的服務解耦，可獨立開發與部署
- ✅ 技術棧自由度（可為不同服務選擇不同資料庫）
- ✅ 資料庫 Schema 變更不影響其他服務
- ✅ 更容易進行服務擴展
- ✅ 符合微服務最佳實踐

**代價**：
- ❌ 無法使用跨資料庫的 JOIN 查詢
- ❌ 需要實作分散式交易（Saga Pattern）
- ❌ 資料一致性需要額外處理
- ❌ 查詢複雜度增加

**緩解措施**：
- 使用 API Composition Pattern 組合多服務資料
- 使用 CQRS + Elasticsearch 處理複雜查詢
- 使用 Saga Pattern 保證最終一致性
- 使用 Kafka 事件驅動架構同步資料

---

### ADR-002: 採用 Saga Pattern 處理分散式交易

**狀態**: 已採用

**背景**：
在 Database per Service 模式下，無法使用傳統的 ACID 交易跨多個服務。

**決策**：
採用 Choreography-based Saga Pattern（編排式 Saga）。

**理由**：
- ✅ 去中心化，無單點故障
- ✅ 服務間耦合度低
- ✅ 易於添加新的參與服務
- ✅ 透過事件驅動，天然支援非同步處理

**代價**：
- ❌ 需要設計補償交易
- ❌ 除錯較困難（需要追蹤事件鏈）
- ❌ 測試複雜度較高

---

### ADR-003: 使用 gRPC 作為同步通訊協定

**狀態**: 已採用

**背景**：
微服務間的同步通訊需要選擇合適的協定（REST vs gRPC）。

**決策**：
內部服務間通訊使用 gRPC，對外 API 使用 RESTful。

**理由**：
- ✅ Protocol Buffers 效能優於 JSON
- ✅ 強型別，編譯期檢查錯誤
- ✅ 支援 streaming（適合爬蟲服務）
- ✅ 自動生成客戶端程式碼
- ✅ 內建負載平衡、超時、重試機制

**實作細節**：
- API Gateway 將外部 HTTP/REST 轉換為內部 gRPC
- 所有內部服務間通訊使用 gRPC
- 使用 gRPC-Gateway 可選擇性提供 REST 介面
