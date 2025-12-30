# 🚀 Tabelogo v2 系統改進指南

> **目標**: 循序漸進提升系統可觀測性、韌性與工程實踐  
> **對象**: 中階工程師技術成長路線圖

---

## 📊 當前系統狀態

### ✅ 已完成項目

| 領域 | 狀態 | 說明 |
|------|------|------|
| 架構 | ✅ | DDD 分層清晰、微服務設計 |
| 依賴注入 | ✅ | 使用 Uber Fx |
| Prometheus | ✅ | 已部署，抓取 3 個服務 |
| Grafana | ✅ | 已部署，2 個 Dashboard |
| Docker Compose | ✅ | 完整配置 |

### ⚠️ 待改進項目

| 領域 | 問題 | 優先級 |
|------|------|--------|
| 日誌聚合 | 無 Loki，日誌未集中 | 🔴 高 |
| Prometheus | Auth Service 未抓取 | 🔴 高 |
| 分佈式追蹤 | `pkg/tracing` 目錄為空 | 🟡 中 |
| 測試覆蓋 | 需要更完整的測試 | 🟡 中 |
| 韌性機制 | 無 Circuit Breaker | 🟡 中 |

---

## 📅 Week 1: 日誌聚合系統 (Loki + Promtail)

### 🎯 目標
建立集中式日誌系統，所有服務日誌可在 Grafana 統一查詢

### Step 1.1: 創建 Loki 配置文件

```bash
# 創建目錄
mkdir -p deployments/loki
```

創建 `deployments/loki/loki-config.yml`:

```yaml
auth_enabled: false

server:
  http_listen_port: 3100

common:
  path_prefix: /loki
  storage:
    filesystem:
      chunks_directory: /loki/chunks
      rules_directory: /loki/rules
  replication_factor: 1
  ring:
    instance_addr: 127.0.0.1
    kvstore:
      store: inmemory

schema_config:
  configs:
    - from: 2020-10-24
      store: boltdb-shipper
      object_store: filesystem
      schema: v11
      index:
        prefix: index_
        period: 24h

storage_config:
  boltdb_shipper:
    active_index_directory: /loki/index
    cache_location: /loki/cache
    shared_store: filesystem
  filesystem:
    directory: /loki/chunks

limits_config:
  reject_old_samples: true
  reject_old_samples_max_age: 168h

chunk_store_config:
  max_look_back_period: 0s

table_manager:
  retention_deletes_enabled: false
  retention_period: 0s
```

### Step 1.2: 創建 Promtail 配置文件

創建 `deployments/loki/promtail-config.yml`:

```yaml
server:
  http_listen_port: 9080
  grpc_listen_port: 0

positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://loki:3100/loki/api/v1/push

scrape_configs:
  - job_name: containers
    docker_sd_configs:
      - host: unix:///var/run/docker.sock
        refresh_interval: 5s
    relabel_configs:
      # 提取容器名稱
      - source_labels: ['__meta_docker_container_name']
        regex: '/(.*)'
        target_label: 'container'
      # 提取服務名稱 (從 tabelogo- 前綴提取)
      - source_labels: ['__meta_docker_container_name']
        regex: '/tabelogo-(.*)'
        target_label: 'service'
      # 只收集 tabelogo 相關容器
      - source_labels: ['__meta_docker_container_name']
        regex: '/tabelogo-.*'
        action: keep
```

### Step 1.3: 更新 docker-compose.yml

在 `deployments/docker-compose/docker-compose.yml` 的 `# Monitoring & Tools` 區塊添加:

```yaml
  # Loki - Log aggregation
  loki:
    image: grafana/loki:3.0.0
    container_name: tabelogo-loki
    ports:
      - "3100:3100"
    volumes:
      - ../loki/loki-config.yml:/etc/loki/local-config.yaml
      - loki-data:/loki
    command: -config.file=/etc/loki/local-config.yaml
    networks:
      - tabelogo-network
    restart: unless-stopped

  # Promtail - Log collector
  promtail:
    image: grafana/promtail:3.0.0
    container_name: tabelogo-promtail
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ../loki/promtail-config.yml:/etc/promtail/config.yml
    command: -config.file=/etc/promtail/config.yml
    depends_on:
      - loki
    networks:
      - tabelogo-network
    restart: unless-stopped
```

在 `volumes:` 區塊添加:

```yaml
  loki-data:
    driver: local
```

### Step 1.4: 驗證 Loki 部署

```bash
# 重新啟動服務
cd deployments/docker-compose
docker-compose up -d loki promtail

# 驗證 Loki 啟動
curl http://localhost:3100/ready
# 預期輸出: ready

# 查看 Promtail 狀態
docker logs tabelogo-promtail
```

### Step 1.5: 配置 Grafana Loki 數據源

1. 開啟 Grafana: http://localhost:3001
2. 登入 (admin/admin)
3. 進入 **Configuration > Data Sources > Add data source**
4. 選擇 **Loki**
5. 設定 URL: `http://loki:3100`
6. 點擊 **Save & Test**

### Step 1.6: 創建日誌 Dashboard

在 Grafana 創建新 Dashboard，添加以下 Panel:

**Panel 1: 所有服務日誌**
```
{container=~"tabelogo-.*"}
```

**Panel 2: 錯誤日誌**
```
{container=~"tabelogo-.*"} |~ "(?i)error|panic|fatal"
```

**Panel 3: 按服務過濾**
```
{service="auth-service"}
```

### ✅ Week 1 完成檢查清單

- [ ] 創建 `deployments/loki/loki-config.yml`
- [ ] 創建 `deployments/loki/promtail-config.yml`
- [ ] 更新 `docker-compose.yml` 添加 Loki 和 Promtail
- [ ] 驗證 Loki 啟動 (`curl http://localhost:3100/ready`)
- [ ] Grafana 添加 Loki 數據源
- [ ] 創建日誌查詢 Dashboard
- [ ] 測試日誌查詢功能

---

## 📅 Week 2: 完善 Prometheus 指標抓取

### 🎯 目標
啟用所有服務的指標收集，添加關鍵業務指標

### Step 2.1: 啟用 Auth Service 指標抓取

更新 `deployments/docker-compose/prometheus.yml`:

```yaml
  # 取消註解並修正 port
  - job_name: 'auth-service'
    metrics_path: '/metrics'
    static_configs:
      - targets: ['auth-service:8080']
    relabel_configs:
      - source_labels: [__address__]
        target_label: instance
        replacement: 'auth-service'
```

### Step 2.2: 確認 Auth Service 有 Metrics Endpoint

檢查 Auth Service 是否已暴露 `/metrics` endpoint。如果沒有，需要添加:

```go
// 在 auth-service 的 router 中添加
import "github.com/prometheus/client_golang/prometheus/promhttp"

router.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

### Step 2.3: 創建 Auth Service 指標

創建 `pkg/metrics/auth_metrics.go`:

```go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // Auth - Login Metrics
    AuthLoginTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "auth_login_total",
            Help: "Total number of login attempts",
        },
        []string{"status"}, // success, failed, blocked
    )

    AuthLoginDuration = promauto.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "auth_login_duration_seconds",
            Help:    "Login request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
    )

    // Auth - Registration Metrics
    AuthRegisterTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "auth_register_total",
            Help: "Total number of registration attempts",
        },
        []string{"status"}, // success, failed
    )

    // Auth - Token Metrics
    AuthTokenRefreshTotal = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "auth_token_refresh_total",
            Help: "Total number of token refresh operations",
        },
    )

    AuthTokenValidationTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "auth_token_validation_total",
            Help: "Total number of token validation attempts",
        },
        []string{"status"}, // valid, expired, invalid
    )

    // Auth - Session Metrics
    AuthActiveSessionsGauge = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "auth_active_sessions",
            Help: "Current number of active sessions",
        },
    )

    // Auth - Security Metrics
    AuthFailedLoginAttempts = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "auth_failed_login_attempts_total",
            Help: "Total failed login attempts (for security monitoring)",
        },
        []string{"reason"}, // wrong_password, user_not_found, account_locked
    )
)
```

### Step 2.4: 在 Auth Service 中使用指標

在登入邏輯中使用:

```go
import "github.com/Leon180/tabelogo-v2/pkg/metrics"

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
    timer := prometheus.NewTimer(metrics.AuthLoginDuration)
    defer timer.ObserveDuration()

    // ... 登入邏輯 ...
    
    if err != nil {
        metrics.AuthLoginTotal.WithLabelValues("failed").Inc()
        metrics.AuthFailedLoginAttempts.WithLabelValues("wrong_password").Inc()
        return nil, err
    }
    
    metrics.AuthLoginTotal.WithLabelValues("success").Inc()
    return response, nil
}
```

### Step 2.5: 創建資料庫連線指標

創建 `pkg/metrics/db_metrics.go`:

```go
package metrics

import (
    "database/sql"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    DBConnectionsActive = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "db_connections_active",
            Help: "Number of active database connections",
        },
        []string{"service"},
    )

    DBConnectionsIdle = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "db_connections_idle",
            Help: "Number of idle database connections",
        },
        []string{"service"},
    )

    DBQueryDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "db_query_duration_seconds",
            Help:    "Database query duration in seconds",
            Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
        },
        []string{"service", "operation"},
    )

    DBSlowQueriesTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "db_slow_queries_total",
            Help: "Total number of slow queries (>100ms)",
        },
        []string{"service"},
    )
)

// RecordDBStats 定期記錄資料庫連線池狀態
func RecordDBStats(db *sql.DB, serviceName string) {
    stats := db.Stats()
    DBConnectionsActive.WithLabelValues(serviceName).Set(float64(stats.InUse))
    DBConnectionsIdle.WithLabelValues(serviceName).Set(float64(stats.Idle))
}
```

### Step 2.6: 驗證指標抓取

```bash
# 重新載入 Prometheus 配置
docker-compose restart prometheus

# 驗證 targets
curl http://localhost:9095/api/v1/targets | jq '.data.activeTargets[] | {job: .labels.job, health: .health}'

# 查詢 auth 指標
curl 'http://localhost:9095/api/v1/query?query=auth_login_total'
```

### ✅ Week 2 完成檢查清單

- [ ] 更新 `prometheus.yml` 啟用 Auth Service
- [ ] 確認 Auth Service 有 `/metrics` endpoint
- [ ] 創建 `pkg/metrics/auth_metrics.go`
- [ ] 創建 `pkg/metrics/db_metrics.go`
- [ ] 在 Auth Service 業務邏輯中埋點
- [ ] 驗證 Prometheus targets 全部 UP
- [ ] 更新 Grafana Dashboard

---

## 📅 Week 3-4: 分佈式追蹤 (OpenTelemetry + Jaeger)

### 🎯 目標
實現跨服務請求追蹤，每個請求有唯一 TraceID

### Step 3.1: 添加 Jaeger 到 Docker Compose

```yaml
  # Jaeger - Distributed tracing
  jaeger:
    image: jaegertracing/all-in-one:1.52
    container_name: tabelogo-jaeger
    environment:
      - COLLECTOR_OTLP_ENABLED=true
    ports:
      - "16686:16686"  # Jaeger UI
      - "4317:4317"    # OTLP gRPC
      - "4318:4318"    # OTLP HTTP
    networks:
      - tabelogo-network
    restart: unless-stopped
```

### Step 3.2: 安裝 OpenTelemetry 依賴

```bash
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc
go get go.opentelemetry.io/otel/sdk/trace
go get go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin
```

### Step 3.3: 創建 Tracer 初始化

創建 `pkg/tracing/tracer.go`:

```go
package tracing

import (
    "context"
    "log"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// InitTracer 初始化 OpenTelemetry tracer
func InitTracer(ctx context.Context, serviceName, jaegerEndpoint string) (*sdktrace.TracerProvider, error) {
    // 創建 OTLP exporter
    exporter, err := otlptracegrpc.New(ctx,
        otlptracegrpc.WithEndpoint(jaegerEndpoint),
        otlptracegrpc.WithInsecure(),
    )
    if err != nil {
        return nil, err
    }

    // 創建 resource
    res, err := resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceName(serviceName),
        ),
    )
    if err != nil {
        return nil, err
    }

    // 創建 TracerProvider
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(res),
        sdktrace.WithSampler(sdktrace.AlwaysSample()),
    )

    // 設置全局 TracerProvider
    otel.SetTracerProvider(tp)
    otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
        propagation.TraceContext{},
        propagation.Baggage{},
    ))

    log.Printf("Tracer initialized for service: %s", serviceName)
    return tp, nil
}

// Shutdown 關閉 tracer
func Shutdown(tp *sdktrace.TracerProvider) {
    if err := tp.Shutdown(context.Background()); err != nil {
        log.Printf("Error shutting down tracer provider: %v", err)
    }
}
```

### Step 3.4: 整合到 Gin Router

```go
import (
    "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func SetupRouter() *gin.Engine {
    router := gin.Default()
    
    // 添加 OpenTelemetry middleware
    router.Use(otelgin.Middleware("auth-service"))
    
    // ... 其他 routes
    return router
}
```

### Step 3.5: 在日誌中添加 TraceID

```go
import (
    "go.opentelemetry.io/otel/trace"
)

func LogWithTrace(ctx context.Context, message string) {
    span := trace.SpanFromContext(ctx)
    traceID := span.SpanContext().TraceID().String()
    spanID := span.SpanContext().SpanID().String()
    
    log.Printf("[TraceID: %s] [SpanID: %s] %s", traceID, spanID, message)
}
```

### Step 3.6: 驗證追蹤

1. 發送請求到任意服務
2. 開啟 Jaeger UI: http://localhost:16686
3. 選擇服務，查看 traces

### ✅ Week 3-4 完成檢查清單

- [ ] 添加 Jaeger 到 docker-compose.yml
- [ ] 安裝 OpenTelemetry 依賴
- [ ] 創建 `pkg/tracing/tracer.go`
- [ ] 各服務整合 otelgin middleware
- [ ] 日誌添加 TraceID
- [ ] 驗證 Jaeger UI 可查看 traces
- [ ] gRPC 調用也添加 trace propagation

---

## 📅 Month 2: 韌性測試 (Chaos Engineering)

### 🎯 目標
驗證系統在異常情況下的行為

### Redis 韌性測試

#### 測試 1: 緩存穿透測試

```bash
# 查詢大量不存在的 ID
for i in {1..100}; do
  curl "http://localhost:18082/api/v1/restaurants/nonexistent-id-$i"
done

# 觀察:
# 1. 請求是否都打到資料庫?
# 2. 是否有空值緩存機制?
```

**解決方案**: 實現空值緩存

```go
func (r *RedisRepo) GetWithNullCache(ctx context.Context, key string) (string, error) {
    val, err := r.client.Get(ctx, key).Result()
    if err == redis.Nil {
        return "", ErrNotFound
    }
    if val == "NULL" {
        return "", ErrNotFound // 空值緩存
    }
    return val, err
}

func (r *RedisRepo) SetNullCache(ctx context.Context, key string) error {
    return r.client.Set(ctx, key, "NULL", time.Minute*5).Err()
}
```

#### 測試 2: 緩存雪崩測試

```bash
# 模擬 Redis 宕機
docker stop tabelogo-redis

# 發送請求
curl http://localhost:18082/api/v1/restaurants/123

# 觀察:
# 1. 服務是否正常降級?
# 2. 是否有錯誤螢幕?

# 恢復 Redis
docker start tabelogo-redis
```

**解決方案**: 添加降級邏輯

```go
func (s *Service) GetRestaurant(ctx context.Context, id string) (*Restaurant, error) {
    // 嘗試從緩存獲取
    cached, err := s.cache.Get(ctx, id)
    if err == nil {
        return cached, nil
    }
    
    // 緩存失敗，直接查資料庫 (降級)
    if err != redis.Nil {
        metrics.CacheFailuresTotal.Inc()
        log.Warn("Cache unavailable, falling back to database")
    }
    
    return s.repo.GetByID(ctx, id)
}
```

### Slow SQL 測試

#### 測試 3: 慢查詢測試

```sql
-- 在 PostgreSQL 中模擬慢查詢
SELECT pg_sleep(5); -- 5秒延遲
```

```go
// 測試代碼
func TestSlowQuery(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    
    _, err := db.ExecContext(ctx, "SELECT pg_sleep(5)")
    
    // 應該超時
    assert.Error(t, err)
    assert.True(t, errors.Is(ctx.Err(), context.DeadlineExceeded))
}
```

**解決方案**: 查詢超時配置

```go
// 設置連線池參數
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)

// 使用 context timeout
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()
result, err := db.QueryContext(ctx, query)
```

### ✅ Month 2 完成檢查清單

#### Redis 韌性
- [ ] 穿透測試: 驗證空值緩存機制
- [ ] 雪崩測試: 驗證降級策略
- [ ] 擊穿測試: 驗證互斥鎖
- [ ] 實現布隆過濾器或空值緩存
- [ ] 添加隨機 TTL 防止同時過期

#### SQL 韌性
- [ ] 模擬慢查詢場景
- [ ] 驗證查詢超時機制
- [ ] 測試連線池耗盡情況
- [ ] 添加慢查詢日誌監控
- [ ] 考慮實現 Circuit Breaker

---

## 📅 Month 3+: 進階工程能力

### ☁️ 雲端部署 (AWS/GCP)

- [ ] 學習 Terraform 基礎 (IaC)
- [ ] 部署服務到 AWS ECS/Fargate
- [ ] 配置 RDS PostgreSQL
- [ ] 配置 ElastiCache Redis
- [ ] 設置 CloudWatch 監控
- [ ] 配置 ALB 負載均衡

### ⚓ Kubernetes 入門

- [ ] 學習 K8s 核心概念 (Pod, Deployment, Service)
- [ ] 使用 Minikube 本地測試
- [ ] 編寫 Helm Charts
- [ ] 配置 HPA 自動擴縮
- [ ] 設置 Ingress Controller

### 🔄 CI/CD 流水線

- [ ] GitHub Actions 自動測試
- [ ] Docker 鏡像自動構建
- [ ] 自動部署到 Staging
- [ ] 藍綠/金絲雀部署
- [ ] Semantic Versioning

### 🔐 安全加固

- [ ] 實現 Rate Limiting
- [ ] 添加 CORS 配置
- [ ] 密鑰管理 (Vault/Secrets Manager)
- [ ] SQL 注入防護審查
- [ ] 安全頭設置 (HSTS, CSP)

---

## 🏗️ 架構全景圖

```
┌─────────────────────────────────────────────────────────┐
│                     Grafana (3001)                      │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐                   │
│  │ Metrics │ │  Logs   │ │ Traces  │                   │
│  └────┬────┘ └────┬────┘ └────┬────┘                   │
└───────┼──────────┼──────────┼───────────────────────────┘
        │          │          │
        ▼          ▼          ▼
   Prometheus    Loki      Jaeger
    (9095)      (3100)    (16686)
        ▲          ▲          ▲
        │          │          │
   ┌────┴────┐ ┌───┴───┐ ┌───┴───┐
   │ /metrics│ │Promtail│ │ OTEL  │
   │Endpoint │ │       │ │  SDK  │
   └────┬────┘ └───┬───┘ └───┬───┘
        │          │          │
   ┌────┴──────────┴──────────┴────┐
   │      Microservices            │
   │  Auth │ Restaurant │ Map │... │
   │ (8080)│  (18082)   │(8081)│   │
   └───────────────────────────────┘
```

---

## 💡 學習資源

| 主題 | 資源 |
|------|------|
| Prometheus | [prometheus.io/docs](https://prometheus.io/docs) |
| Loki | [grafana.com/docs/loki](https://grafana.com/docs/loki) |
| OpenTelemetry | [opentelemetry.io/docs](https://opentelemetry.io/docs) |
| Go 可觀測性 | Practical Go: Observability |
| Chaos Engineering | [Principles of Chaos](https://principlesofchaos.org) |

---

## 📋 快速驗證命令

```bash
# 查看待處理的 TODO
grep -rn "TODO\|FIXME" internal/

# 檢查現有指標
curl http://localhost:18082/metrics | head -20

# 驗證 Prometheus targets
curl http://localhost:9095/api/v1/targets | jq '.data.activeTargets[] | {job: .labels.job, health: .health}'

# 驗證 Loki
curl http://localhost:3100/ready

# 查看容器日誌
docker logs tabelogo-auth-service --tail 50
```

---

> **提示**: 建議從 Week 1 開始，完成每週任務後勾選檢查清單。遇到問題時可參考各步驟的驗證命令。
