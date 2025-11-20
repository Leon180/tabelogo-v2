# Database Migrations Summary

完整的微服務資料庫 Migration 設計，遵循 **Database per Service** 原則。

## 📊 服務與資料庫對應

| 服務 | 資料庫 | 連接端口 | 主要資料表 |
|------|--------|----------|-----------|
| Auth Service | `auth_db` | 5432 | users, refresh_tokens |
| Restaurant Service | `restaurant_db` | 5433 | restaurants, user_favorites |
| Booking Service | `booking_db` | 5434 | bookings, booking_history |
| Spider Service | `spider_db` | 5435 | crawl_jobs, crawl_results |
| Mail Service | `mail_db` | 5436 | email_queue, email_logs |

---

## 🗂️ Migration 檔案結構

```
migrations/
├── auth/
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_refresh_tokens_table.up.sql
│   └── 000002_create_refresh_tokens_table.down.sql
├── restaurant/
│   ├── 000001_create_restaurants_table.up.sql
│   ├── 000001_create_restaurants_table.down.sql
│   ├── 000002_create_user_favorites_table.up.sql
│   └── 000002_create_user_favorites_table.down.sql
├── booking/
│   ├── 000001_create_bookings_table.up.sql
│   ├── 000001_create_bookings_table.down.sql
│   ├── 000002_create_booking_history_table.up.sql
│   └── 000002_create_booking_history_table.down.sql
├── spider/
│   ├── 000001_create_crawl_jobs_table.up.sql
│   ├── 000001_create_crawl_jobs_table.down.sql
│   ├── 000002_create_crawl_results_table.up.sql
│   └── 000002_create_crawl_results_table.down.sql
└── mail/
    ├── 000001_create_email_queue_table.up.sql
    ├── 000001_create_email_queue_table.down.sql
    ├── 000002_create_email_logs_table.up.sql
    └── 000002_create_email_logs_table.down.sql
```

---

## 📋 各服務資料表詳細說明

### 1️⃣ Auth Service (`auth_db`)

#### **users** - 使用者資料表
- 儲存使用者認證資訊
- 使用 bcrypt 密碼雜湊
- 支援軟刪除 (soft delete)
- RBAC 角色管理 (admin, user, guest)

**主要欄位：**
- `id` (UUID) - 主鍵
- `email` - 唯一 email
- `password_hash` - bcrypt 密碼雜湊
- `role` - 使用者角色
- `is_active` - 帳號啟用狀態
- `email_verified` - Email 驗證狀態

#### **refresh_tokens** - Refresh Token 管理
- 儲存 JWT Refresh Token
- 支援 Token 撤銷 (revocation)
- 自動過期機制

**主要欄位：**
- `user_id` - 關聯使用者 (無 FK，微服務原則)
- `token_hash` - Token 雜湊值
- `expires_at` - 過期時間
- `revoked_at` - 撤銷時間

---

### 2️⃣ Restaurant Service (`restaurant_db`)

#### **restaurants** - 餐廳主資料表
- 聚合多個外部來源的餐廳資料
- 支援地理位置查詢
- JSONB 儲存營業時間與 metadata

**主要欄位：**
- `id` (UUID) - 主鍵
- `name` - 餐廳名稱
- `source` - 資料來源 (tabelog, google, instagram)
- `external_id` - 外部來源 ID
- `latitude`, `longitude` - 地理座標
- `rating` - 評分
- `opening_hours` (JSONB) - 營業時間
- `metadata` (JSONB) - 額外資訊

**重要索引：**
- Unique index on `(source, external_id)` - 防止重複資料
- GIN index on JSONB 欄位 - 加速 JSON 查詢

#### **user_favorites** - 使用者收藏餐廳
- 使用者可收藏餐廳
- 支援私人筆記與標籤
- 記錄造訪次數

**主要欄位：**
- `user_id` - 使用者 ID (來自 auth_db)
- `restaurant_id` - 餐廳 ID
- `notes` - 私人筆記
- `tags` - 使用者標籤 (Array)
- `visit_count` - 造訪次數
- `last_visited_at` - 最後造訪時間

---

### 3️⃣ Booking Service (`booking_db`)

#### **bookings** - 訂位資料表
- 儲存所有訂位記錄
- 與外部服務 (OpenTable) 同步
- 支援多種訂位狀態

**主要欄位：**
- `user_id` - 使用者 ID (來自 auth_db)
- `restaurant_id` - 餐廳 ID (來自 restaurant_db)
- `booking_date` - 訂位日期時間
- `party_size` - 人數 (1-50)
- `status` - 狀態 (pending, confirmed, cancelled, completed, no_show)
- `external_booking_id` - 外部服務訂位 ID
- `external_service` - 外部服務名稱 (opentable, tabelog)
- `confirmation_code` - 確認碼
- `last_synced_at` - 最後同步時間

**重要索引：**
- Index on `(user_id, status)` - 查詢使用者訂位
- Index on `(restaurant_id, booking_date)` - 查詢餐廳訂位
- Partial index on upcoming bookings - 查詢即將到來的訂位

#### **booking_history** - 訂位歷史 (Event Sourcing)
- 記錄所有訂位狀態變更
- Event Sourcing 模式實作
- 完整審計追蹤

**主要欄位：**
- `booking_id` - 訂位 ID
- `change_type` - 變更類型 (created, updated, confirmed, cancelled, synced)
- `previous_value` (JSONB) - 變更前狀態
- `new_value` (JSONB) - 變更後狀態
- `changed_by` - 變更者 user_id
- `metadata` (JSONB) - 額外資訊 (IP, user agent, sync source)

---

### 4️⃣ Spider Service (`spider_db`)

#### **crawl_jobs** - 爬蟲任務管理
- 管理爬蟲任務執行
- 支援優先權排程
- 追蹤任務進度與錯誤

**主要欄位：**
- `source` - 來源 (tabelog, google_maps, instagram)
- `region` - 爬取區域
- `job_type` - 任務類型 (full, incremental, update)
- `status` - 狀態 (pending, running, completed, failed, cancelled, paused)
- `priority` - 優先權 (1-10)
- `total_pages`, `completed_pages` - 進度追蹤
- `config` (JSONB) - 爬蟲設定 (rate limit, proxy, user agent)
- `next_run_at` - 下次執行時間 (定期任務)

#### **crawl_results** - 爬蟲結果儲存
- 儲存爬取的原始資料
- 支援去重 (checksum)
- 分離 raw data 與 parsed data

**主要欄位：**
- `job_id` - 爬蟲任務 ID
- `external_id` - 外部來源 ID
- `source` - 資料來源
- `url` - 來源 URL
- `raw_data` (JSONB) - 原始爬取資料 (reviews, ratings, photos, hours)
- `parsed_data` (JSONB) - 解析後的標準化資料
- `checksum` - 資料雜湊值 (去重用)
- `status` - 處理狀態 (pending, processed, failed, duplicate)

**重要特性：**
- Unique index on `(source, external_id)` - 防止重複爬取
- GIN index on JSONB - 支援複雜查詢

---

### 5️⃣ Mail Service (`mail_db`)

#### **email_queue** - 郵件佇列
- 非同步郵件發送佇列
- 支援排程與優先權
- 重試機制

**主要欄位：**
- `recipient_email`, `recipient_name` - 收件人資訊
- `subject`, `body`, `html_body` - 郵件內容
- `template_name` - 郵件模板名稱 (welcome, booking_confirmation, password_reset)
- `template_data` (JSONB) - 模板變數
- `priority` - 優先權 (1-10)
- `status` - 狀態 (pending, sending, sent, failed, cancelled)
- `retry_count`, `max_retries` - 重試機制
- `scheduled_at` - 排程發送時間
- `external_id` - 外部服務 ID (SendGrid, AWS SES)

**重要索引：**
- Index on pending emails - 快速找到待發送郵件
- Index on failed emails with retries - 重試機制

#### **email_logs** - 郵件發送日誌
- 記錄所有郵件事件
- 支援 webhook 資料儲存
- 追蹤開信、點擊等行為

**主要欄位：**
- `email_queue_id` - 郵件佇列 ID
- `status` - 狀態 (sent, failed, bounced, opened, clicked)
- `event_type` - 事件類型 (delivered, opened, clicked, bounced, spam_report, unsubscribed)
- `webhook_data` (JSONB) - Webhook 原始資料
- `metadata` (JSONB) - 額外追蹤資訊 (IP, user agent)

---

## 🔑 關鍵設計原則

### 1. Database per Service
- ✅ 每個微服務擁有獨立資料庫
- ✅ 避免跨資料庫外鍵約束 (Foreign Key)
- ✅ 使用 UUID 作為主鍵 (分散式友善)

### 2. 跨服務資料引用
```sql
-- ❌ 不使用 FOREIGN KEY (跨資料庫)
user_id UUID NOT NULL REFERENCES auth_db.users(id)

-- ✅ 只儲存 ID，不建立約束
user_id UUID NOT NULL  -- Reference to auth_db.users (no FK)
```

### 3. 審計欄位 (Audit Fields)
所有主要資料表都包含：
- `created_at` - 建立時間
- `updated_at` - 更新時間 (自動觸發器)
- `deleted_at` - 軟刪除時間 (可選)

### 4. JSONB 欄位使用
- `metadata`, `config`, `template_data` 等彈性欄位
- 使用 GIN index 加速查詢
- 範例：
  ```sql
  opening_hours JSONB,  -- 營業時間
  metadata JSONB,       -- 額外資訊
  ```

### 5. 索引策略
- **B-tree index**：一般查詢 (id, email, status)
- **GIN index**：JSONB、Array 欄位
- **Partial index**：WHERE 條件索引 (`deleted_at IS NULL`)
- **Composite index**：複合查詢 (`user_id, status`)

---

## 🚀 執行 Migration

### 使用專案的 Migration Manager

```go
// 範例：執行 Auth Service migrations
import (
    "github.com/lileon/tabelogov2/pkg/migrations"
)

manager, err := migrations.NewManager(migrations.Config{
    DB:             db,
    Logger:         logger,
    MigrationsPath: "file://migrations/auth",
    ServiceName:    "auth",
})

// 執行所有 migrations
err = manager.Up(ctx)
```

### 使用 golang-migrate CLI

```bash
# Auth Service
migrate -path migrations/auth \
        -database "postgresql://postgres:postgres@localhost:5432/auth_db?sslmode=disable" \
        up

# Restaurant Service
migrate -path migrations/restaurant \
        -database "postgresql://postgres:postgres@localhost:5433/restaurant_db?sslmode=disable" \
        up

# Booking Service
migrate -path migrations/booking \
        -database "postgresql://postgres:postgres@localhost:5434/booking_db?sslmode=disable" \
        up

# Spider Service
migrate -path migrations/spider \
        -database "postgresql://postgres:postgres@localhost:5435/spider_db?sslmode=disable" \
        up

# Mail Service
migrate -path migrations/mail \
        -database "postgresql://postgres:postgres@localhost:5436/mail_db?sslmode=disable" \
        up
```

### 回滾 Migration

```bash
# 回滾最後一個 migration
migrate -path migrations/auth \
        -database "postgresql://..." \
        down 1

# 回滾到特定版本
migrate -path migrations/auth \
        -database "postgresql://..." \
        migrate 1
```

---

## 📊 資料表關聯圖

```
┌─────────────────┐
│    auth_db      │
├─────────────────┤
│ users           │─┐
│ refresh_tokens  │ │
└─────────────────┘ │
                    │ (user_id, no FK)
┌─────────────────┐ │
│ restaurant_db   │ │
├─────────────────┤ │
│ restaurants     │◄┼─────┐
│ user_favorites  │◄┘     │ (restaurant_id, no FK)
└─────────────────┘       │
                          │
┌─────────────────┐       │
│  booking_db     │       │
├─────────────────┤       │
│ bookings        │───────┤
│ booking_history │       │
└─────────────────┘       │
        ▲                 │
        │ (external_booking_id)
        │                 │
┌─────────────────┐       │
│   spider_db     │       │
├─────────────────┤       │
│ crawl_jobs      │       │
│ crawl_results   │───────┘ (產生餐廳資料)
└─────────────────┘

┌─────────────────┐
│    mail_db      │
├─────────────────┤
│ email_queue     │ (發送 booking confirmation)
│ email_logs      │
└─────────────────┘
```

---

## ✅ Migration 檢查清單

- [x] Auth Service - users, refresh_tokens
- [x] Restaurant Service - restaurants, user_favorites
- [x] Booking Service - bookings, booking_history
- [x] Spider Service - crawl_jobs, crawl_results
- [x] Mail Service - email_queue, email_logs
- [x] 所有表格包含審計欄位 (created_at, updated_at)
- [x] 所有表格使用 UUID 主鍵
- [x] 適當的索引設計 (B-tree, GIN, Partial, Composite)
- [x] 軟刪除支援 (deleted_at)
- [x] 跨服務資料引用不使用 FK
- [x] JSONB 欄位用於彈性資料
- [x] 完整的 up/down migration 檔案

---

## 📝 注意事項

1. **跨服務資料一致性**
   - 使用 Saga Pattern 處理分散式交易
   - 使用 Kafka 事件同步資料
   - 接受 Eventual Consistency

2. **外部服務同步**
   - `bookings.external_booking_id` 用於與 OpenTable 同步
   - `bookings.last_synced_at` 追蹤同步時間
   - 使用 Webhook + Polling 雙重機制

3. **效能考量**
   - JSONB 欄位建立 GIN index
   - 高頻查詢建立 Composite index
   - 使用 Partial index 過濾軟刪除資料

4. **安全性**
   - 密碼使用 bcrypt hash
   - Token 使用 hash 儲存
   - 敏感資訊不明文儲存

---

## 🎯 下一步

1. 啟動 docker-compose 環境
2. 執行所有 migrations
3. 驗證資料表建立成功
4. 開始開發各微服務的 Domain Layer

```bash
# 啟動所有資料庫
cd deployments/docker-compose
docker-compose up -d

# 等待資料庫就緒
sleep 10

# 執行 migrations (使用專案提供的工具)
# TODO: 建立 migration 執行腳本
```
