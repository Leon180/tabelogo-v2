# 環境變量參考文檔

本文檔列出了所有微服務使用的環境變量及其說明。

## Auth Service 環境變量

### 基本配置

| 變量名 | 類型 | 默認值 | 說明 |
|--------|------|--------|------|
| `ENVIRONMENT` | string | `development` | 運行環境：development, staging, production, test |
| `LOG_LEVEL` | string | `info` | 日誌級別：debug, info, warn, error |
| `SERVER_PORT` | int | `8080` | HTTP 服務器端口 |
| `GRPC_PORT` | int | `9090` | gRPC 服務器端口 |

### 數據庫配置 (PostgreSQL)

| 變量名 | 類型 | 默認值 | 說明 |
|--------|------|--------|------|
| `DB_HOST` | string | `localhost` | PostgreSQL 主機地址 |
| `DB_PORT` | int | `5432` | PostgreSQL 端口 |
| `DB_NAME` | string | **必填** | 數據庫名稱 |
| `DB_USER` | string | `postgres` | 數據庫用戶名 |
| `DB_PASSWORD` | string | `postgres` | 數據庫密碼 |
| `DB_SSLMODE` | string | `disable` | SSL 模式：disable, require, verify-ca, verify-full |
| `DB_MAX_OPEN_CONNS` | int | `100` | 最大打開連接數 |
| `DB_MAX_IDLE_CONNS` | int | `10` | 最大空閒連接數 |
| `DB_CONN_MAX_LIFETIME` | duration | `1h` | 連接最大生命週期 |

### Redis 配置

| 變量名 | 類型 | 默認值 | 說明 |
|--------|------|--------|------|
| `REDIS_HOST` | string | `localhost` | Redis 主機地址 |
| `REDIS_PORT` | int | `6379` | Redis 端口 |
| `REDIS_PASSWORD` | string | `""` | Redis 密碼（可選） |
| `REDIS_DB` | int | `0` | Redis 數據庫編號（0-15） |

### JWT 配置

| 變量名 | 類型 | 默認值 | 說明 |
|--------|------|--------|------|
| `JWT_SECRET` | string | **必填** | JWT 簽名密鑰（生產環境必須更改） |
| `JWT_ACCESS_TOKEN_EXPIRE` | duration | `15m` | 訪問令牌過期時間 |
| `JWT_REFRESH_TOKEN_EXPIRE` | duration | `168h` (7天) | 刷新令牌過期時間 |

### Kafka 配置

| 變量名 | 類型 | 默認值 | 說明 |
|--------|------|--------|------|
| `KAFKA_BROKERS` | string | `localhost:9092` | Kafka brokers（逗號分隔） |
| `KAFKA_GROUP_ID` | string | `tabelogo-group` | Kafka 消費者組 ID |

---

## 本地開發環境配置

### Auth Service (VSCode Launch)

本地開發時使用以下配置避免端口衝突：

```bash
ENVIRONMENT=development
SERVER_PORT=8081
GRPC_PORT=9091

# Database (使用特殊端口)
DB_HOST=localhost
DB_PORT=15432
DB_NAME=auth_db
DB_USER=postgres
DB_PASSWORD=postgres

# Redis (使用特殊端口)
REDIS_HOST=localhost
REDIS_PORT=16379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_ACCESS_TOKEN_EXPIRE=15m
JWT_REFRESH_TOKEN_EXPIRE=168h
```

### 為什麼使用不同的端口？

本地開發環境使用特殊端口是為了避免與系統已安裝的服務衝突：

| 服務 | 標準端口 | 本地開發端口 | 原因 |
|------|---------|-------------|------|
| PostgreSQL | 5432 | **15432** | 避免與本機已安裝的 PostgreSQL 衝突 |
| Redis | 6379 | **16379** | 避免與本機已安裝的 Redis 衝突 |
| Auth HTTP | 8080 | **8081** | 8080 常被其他服務佔用 |
| Auth gRPC | 9090 | **9091** | 與 HTTP 端口保持一致性 |

---

## 環境變量類型說明

### Duration 格式

Duration 類型使用 Go 的 time.Duration 格式：

- `s` - 秒（例如：`30s`）
- `m` - 分鐘（例如：`15m`）
- `h` - 小時（例如：`24h`）
- 組合使用（例如：`1h30m`）

### 布爾值

布爾類型可以使用以下值（不區分大小寫）：
- True: `true`, `1`, `yes`, `on`
- False: `false`, `0`, `no`, `off`

---

## 配置優先級

1. 環境變量
2. 配置文件（如果實現）
3. 默認值

---

## 安全建議

### 生產環境

❌ **絕對不要**在生產環境使用以下值：

```bash
JWT_SECRET=change-me-in-production
DB_PASSWORD=postgres
REDIS_PASSWORD=
```

✅ **必須使用強密碼和密鑰**：

```bash
JWT_SECRET=使用至少32字符的隨機字符串
DB_PASSWORD=使用強密碼
REDIS_PASSWORD=使用強密碼（如果 Redis 暴露在網絡中）
```

### 密鑰生成

生成強 JWT 密鑰：

```bash
# 使用 openssl
openssl rand -base64 32

# 使用 Python
python3 -c "import secrets; print(secrets.token_urlsafe(32))"

# 使用 Node.js
node -e "console.log(require('crypto').randomBytes(32).toString('base64'))"
```

### 敏感信息管理

生產環境建議使用密鑰管理服務：
- AWS Secrets Manager
- HashiCorp Vault
- Azure Key Vault
- Google Cloud Secret Manager

---

## 驗證配置

### 檢查當前配置

啟動服務時會自動驗證配置。查看日誌確認：

```bash
# 查看服務日誌
docker logs tabelogo-auth-service

# 或在 VSCode Debug Console 中查看
```

### 配置驗證規則

- ✅ 所有端口必須在 1-65535 範圍內
- ✅ `SERVER_PORT` 和 `GRPC_PORT` 不能相同
- ✅ `DB_MAX_IDLE_CONNS` 不能超過 `DB_MAX_OPEN_CONNS`
- ✅ `REDIS_DB` 必須在 0-15 範圍內
- ✅ `JWT_REFRESH_TOKEN_EXPIRE` 必須大於 `JWT_ACCESS_TOKEN_EXPIRE`
- ✅ 生產環境不能使用默認的 JWT_SECRET

---

## 常見錯誤

### 錯誤 1: 連接被拒絕

```
dial tcp [::1]:6379: connect: connection refused
```

**原因**: Redis 配置錯誤
- 檢查 `REDIS_HOST` 和 `REDIS_PORT` 是否正確
- 確認 Redis 服務正在運行

**解決**:
```bash
# 本地開發應使用
REDIS_HOST=localhost
REDIS_PORT=16379

# 而不是
REDIS_ADDR=localhost:16379  # ❌ 錯誤！此變量不存在
```

### 錯誤 2: 數據庫連接失敗

```
failed to connect to database: connection refused
```

**原因**: PostgreSQL 配置錯誤或服務未啟動

**解決**:
```bash
# 檢查 PostgreSQL 狀態
docker ps | grep postgres-auth

# 確保使用正確端口
DB_PORT=15432  # 本地開發端口
```

### 錯誤 3: JWT 密鑰驗證失敗

```
JWT_SECRET must be changed in production environment
```

**原因**: 生產環境使用了默認密鑰

**解決**: 在生產環境設置強 JWT 密鑰

---

## 其他微服務

其他微服務的環境變量配置類似，主要差異在於：
- 端口號不同
- 數據庫名稱不同
- 服務特定的配置項

待實現的微服務配置文檔：
- Restaurant Service
- Booking Service
- Mail Service
- Spider Service
- Map Service
- API Gateway

---

## 更新歷史

- 2025-11-25: 初始版本，包含 Auth Service 完整配置
