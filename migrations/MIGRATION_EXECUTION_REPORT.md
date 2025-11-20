# 🎉 Migration 執行成功報告

執行時間：2025-11-20

---

## ✅ 執行結果總覽

所有 **5 個微服務** 的資料庫 migrations 已成功執行！

| 服務 | 資料庫 | 連接端口 | Migration 版本 | 資料表數量 | 狀態 |
|------|--------|----------|----------------|-----------|------|
| Auth Service | `auth_db` | **15432** | v2 | 2 | ✅ 成功 |
| Restaurant Service | `restaurant_db` | 5433 | v2 | 2 | ✅ 成功 |
| Booking Service | `booking_db` | 5434 | v2 | 2 | ✅ 成功 |
| Spider Service | `spider_db` | 5435 | v2 | 2 | ✅ 成功 |
| Mail Service | `mail_db` | 5436 | v2 | 2 | ✅ 成功 |

**總計：10 個資料表已建立**

---

## 📋 各服務資料表清單

### 1️⃣ Auth Service (`auth_db`)
```
✓ users              - 使用者認證資料
✓ refresh_tokens     - JWT refresh token 管理
✓ schema_migrations  - Migration 版本追蹤
```

### 2️⃣ Restaurant Service (`restaurant_db`)
```
✓ restaurants        - 餐廳主資料（來自 Google Maps/Tabelog/IG）
✓ user_favorites     - 使用者收藏餐廳
✓ schema_migrations  - Migration 版本追蹤
```

### 3️⃣ Booking Service (`booking_db`)
```
✓ bookings           - 訂位資料（支援外部 API 同步）
✓ booking_history    - 訂位歷史（Event Sourcing）
✓ schema_migrations  - Migration 版本追蹤
```

### 4️⃣ Spider Service (`spider_db`)
```
✓ crawl_jobs         - 爬蟲任務管理
✓ crawl_results      - 爬蟲結果儲存
✓ schema_migrations  - Migration 版本追蹤
```

### 5️⃣ Mail Service (`mail_db`)
```
✓ email_queue        - 郵件發送佇列
✓ email_logs         - 郵件追蹤日誌
✓ schema_migrations  - Migration 版本追蹤
```

---

## 🐳 Docker 容器狀態

```bash
NAMES                   STATUS      PORTS
tabelogo-auth-db        Up 4 hours  0.0.0.0:15432->5432/tcp  ⚠️  注意端口
tabelogo-restaurant-db  Up 3 hours  0.0.0.0:5433->5432/tcp
tabelogo-booking-db     Up 3 hours  0.0.0.0:5434->5432/tcp
tabelogo-spider-db      Up 4 hours  0.0.0.0:5435->5432/tcp
tabelogo-mail-db        Up 3 hours  0.0.0.0:5436->5432/tcp
```

### ⚠️ 重要提醒

**Auth DB 端口變更：**
- 原計劃端口：5432
- 實際端口：**15432**（因為 5432 被佔用）
- 連接字串：`postgresql://postgres:postgres@localhost:15432/auth_db`

---

## 🔧 連接資訊

### 連接字串

```bash
# Auth Service
postgresql://postgres:postgres@localhost:15432/auth_db?sslmode=disable

# Restaurant Service
postgresql://postgres:postgres@localhost:5433/restaurant_db?sslmode=disable

# Booking Service
postgresql://postgres:postgres@localhost:5434/booking_db?sslmode=disable

# Spider Service
postgresql://postgres:postgres@localhost:5435/spider_db?sslmode=disable

# Mail Service
postgresql://postgres:postgres@localhost:5436/mail_db?sslmode=disable
```

### 使用 psql 連接

```bash
# 連接到 auth_db
docker exec -it tabelogo-auth-db psql -U postgres -d auth_db

# 連接到 restaurant_db
docker exec -it tabelogo-restaurant-db psql -U postgres -d restaurant_db

# 連接到 booking_db
docker exec -it tabelogo-booking-db psql -U postgres -d booking_db

# 連接到 spider_db
docker exec -it tabelogo-spider-db psql -U postgres -d spider_db

# 連接到 mail_db
docker exec -it tabelogo-mail-db psql -U postgres -d mail_db
```

---

## 📊 資料表結構驗證

### 檢查資料表

```sql
-- 列出所有資料表
\dt

-- 查看資料表結構
\d users
\d restaurants
\d bookings

-- 查看索引
\di

-- 查看觸發器
\dft
```

### 範例查詢

```sql
-- Auth DB
SELECT * FROM users LIMIT 5;
SELECT * FROM refresh_tokens LIMIT 5;

-- Restaurant DB
SELECT * FROM restaurants LIMIT 5;
SELECT * FROM user_favorites LIMIT 5;

-- Booking DB
SELECT * FROM bookings LIMIT 5;
SELECT * FROM booking_history LIMIT 5;

-- Spider DB
SELECT * FROM crawl_jobs LIMIT 5;
SELECT * FROM crawl_results LIMIT 5;

-- Mail DB
SELECT * FROM email_queue LIMIT 5;
SELECT * FROM email_logs LIMIT 5;
```

---

## 🎯 關鍵設計特點（已實現）

### ✅ 已實現的架構原則

1. **Database per Service** ✅
   - 每個微服務擁有獨立資料庫
   - 5 個獨立的 PostgreSQL 實例

2. **無跨資料庫 FK** ✅
   - user_id、restaurant_id 等欄位不使用 FOREIGN KEY
   - 保持微服務獨立性

3. **UUID 主鍵** ✅
   - 所有主鍵使用 UUID v4
   - 分散式友善設計

4. **軟刪除** ✅
   - deleted_at 欄位支援軟刪除
   - Partial index 過濾已刪除資料

5. **審計追蹤** ✅
   - created_at、updated_at 自動維護
   - Trigger 自動更新 updated_at

6. **JSONB 彈性欄位** ✅
   - metadata、config、template_data 等
   - GIN index 支援複雜查詢

7. **完整索引策略** ✅
   - B-tree index（一般查詢）
   - GIN index（JSONB、Array）
   - Partial index（條件索引）
   - Composite index（組合查詢）

8. **Event Sourcing** ✅
   - booking_history 完整記錄狀態變更
   - 支援審計與回溯

---

## 🚀 下一步建議

### 1. 開發 Domain Layer

開始實作各服務的領域層：

```go
// 範例：Auth Service Domain Layer
internal/auth/domain/
├── user/
│   ├── user.go           // User Aggregate Root
│   ├── repository.go     // Repository Interface
│   └── value_object.go   // Email, Password 等
└── token/
    ├── refresh_token.go
    └── repository.go
```

### 2. 實作 Repository Layer

```go
// 範例：User Repository 實作
internal/auth/infrastructure/persistence/
└── postgres_user_repository.go
```

### 3. 建立測試資料

```bash
# 使用 testcontainers-go 進行整合測試
go test ./internal/auth/infrastructure/persistence/...
```

### 4. 配置環境變數

更新 `.env` 檔案：

```env
# Auth Service
AUTH_DB_HOST=localhost
AUTH_DB_PORT=15432  # ⚠️ 注意：不是 5432
AUTH_DB_NAME=auth_db
AUTH_DB_USER=postgres
AUTH_DB_PASSWORD=postgres

# Restaurant Service
RESTAURANT_DB_HOST=localhost
RESTAURANT_DB_PORT=5433
RESTAURANT_DB_NAME=restaurant_db
RESTAURANT_DB_USER=postgres
RESTAURANT_DB_PASSWORD=postgres

# ... 其他服務類似
```

### 5. 實作 Health Check

```go
// 檢查資料庫連線
func (m *Manager) HealthCheck(ctx context.Context) error {
    return m.db.PingContext(ctx)
}
```

---

## 📝 已修復的問題

### 問題 1: Restaurant DB 缺少 trigger function
**錯誤：** `function update_updated_at_column() does not exist`

**解決方案：** 在 `000001_create_restaurants_table.up.sql` 開頭新增 trigger function 定義

### 問題 2: Booking/Mail DB 的 NOW() 索引問題
**錯誤：** `functions in index predicate must be marked IMMUTABLE`

**解決方案：** 移除 WHERE 子句中的 `NOW()` 函數調用，改為應用層過濾

```sql
-- 修改前
CREATE INDEX idx_bookings_upcoming ON bookings(booking_date)
    WHERE status IN ('pending', 'confirmed') AND booking_date > NOW();

-- 修改後
CREATE INDEX idx_bookings_upcoming ON bookings(booking_date, status)
    WHERE status IN ('pending', 'confirmed');
```

---

## 🎓 學習重點

1. **PostgreSQL Partial Index 限制**
   - WHERE 子句不能使用非 IMMUTABLE 函數（如 NOW()）
   - 解決方案：在應用層過濾或使用其他索引策略

2. **Migration 錯誤處理**
   - Dirty state 需要使用 `force` 命令重置
   - 或者直接重建資料庫

3. **Docker 端口衝突處理**
   - 使用 `lsof -i :PORT` 檢查端口佔用
   - 使用替代端口（如 15432 代替 5432）

4. **Trigger Function 共享**
   - 每個資料庫需要獨立定義 trigger function
   - 不能跨資料庫共享 function

---

## ✅ 檢查清單

- [x] 所有 Docker 容器運行正常
- [x] 5 個資料庫已建立
- [x] 10 個資料表已建立（每個服務 2 個表）
- [x] 所有索引已建立
- [x] 所有觸發器已建立
- [x] Migration 版本追蹤正常（schema_migrations 表）
- [x] 連接測試成功
- [x] 資料表結構驗證通過

---

## 🎉 恭喜！

你的微服務資料庫架構已成功建立！

現在可以開始開發各服務的業務邏輯了。記得遵循：
- DDD 分層架構
- Repository Pattern
- Database per Service 原則
- Event Sourcing（booking_history）

祝開發順利！🚀
