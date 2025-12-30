# Week 1 實作筆記：日誌聚合系統 + gRPC Keepalive 修復

> **日期**: 2025-12-30  
> **分支**: feat/imp

---

## 📁 變更文件總覽

```
deployments/
├── loki/
│   ├── loki-config.yml        # [NEW] Loki 服務配置
│   └── promtail-config.yml    # [NEW] Promtail 日誌收集配置
├── docker-compose/
│   └── docker-compose.yml     # [MODIFIED] 新增 Loki + Promtail 服務
└── grafana/
    └── dashboards/
        └── logs-overview.json # [NEW] Grafana 日誌儀表板

internal/
├── map/interfaces/grpc/
│   └── module.go              # [MODIFIED] 新增服務端 keepalive 配置
└── restaurant/infrastructure/
    └── module.go              # [MODIFIED] 調整客戶端 keepalive 參數
```

---

## 🔧 Loki 配置解析

**文件**: `deployments/loki/loki-config.yml`

```yaml
auth_enabled: false  # 開發環境關閉認證

server:
  http_listen_port: 3100   # Loki HTTP API 埠
  grpc_listen_port: 9096   # Loki gRPC 埠（內部通訊）
```

### Schema 配置 (重要)
```yaml
schema_config:
  configs:
    - from: 2020-10-24
      store: tsdb           # 使用 TSDB 存儲（v3.0 推薦）
      object_store: filesystem
      schema: v13           # Schema 版本 13（最新）
      index:
        prefix: index_
        period: 24h         # 每 24 小時輪替索引
```

> **注意**: Loki 3.0 使用 `tsdb` 取代舊版的 `boltdb-shipper`

### 儲存配置
```yaml
storage_config:
  filesystem:
    directory: /loki/chunks    # 日誌 chunk 存放位置
  tsdb_shipper:
    active_index_directory: /loki/tsdb-index
    cache_location: /loki/tsdb-cache
```

### 限制配置
```yaml
limits_config:
  reject_old_samples: true
  reject_old_samples_max_age: 168h    # 拒絕 7 天前的日誌
  ingestion_rate_mb: 16               # 每秒最大寫入 16MB
  ingestion_burst_size_mb: 32         # 突發最大 32MB
  allow_structured_metadata: true     # 允許結構化元數據
```

---

## 🔧 Promtail 配置解析

**文件**: `deployments/loki/promtail-config.yml`

### 基本配置
```yaml
server:
  http_listen_port: 9080    # Promtail 監控埠

positions:
  filename: /tmp/positions.yaml  # 記錄讀取位置，重啟後繼續

clients:
  - url: http://loki:3100/loki/api/v1/push  # 推送到 Loki
```

### Docker 服務發現
```yaml
scrape_configs:
  - job_name: containers
    docker_sd_configs:
      - host: unix:///var/run/docker.sock  # 連接 Docker socket
        refresh_interval: 5s               # 每 5 秒刷新容器列表
```

### Relabel 配置（標籤提取）
```yaml
relabel_configs:
  # 1. 提取容器名稱（移除前導斜線）
  - source_labels: ['__meta_docker_container_name']
    regex: '/(.*)'
    target_label: 'container'
  
  # 2. 從 tabelogo- 前綴提取服務名
  - source_labels: ['__meta_docker_container_name']
    regex: '/tabelogo-(.*)'
    target_label: 'service'
  
  # 3. 只收集 tabelogo 容器的日誌
  - source_labels: ['__meta_docker_container_name']
    regex: '/tabelogo-.*'
    action: keep    # 只保留符合的，過濾其他容器
```

### Pipeline 處理（日誌解析）
```yaml
pipeline_stages:
  # 嘗試解析 JSON 格式日誌
  - json:
      expressions:
        level: level    # 提取 level 欄位
        msg: msg
        time: time
        error: error
  
  # 將 level 添加為標籤（可在 Grafana 過濾）
  - labels:
      level:
```

---

## 📊 Grafana Dashboard 解析

**文件**: `deployments/grafana/dashboards/logs-overview.json`

### Dashboard 結構

| Panel ID | 名稱 | 用途 | LogQL 查詢 |
|----------|------|------|------------|
| 1 | 🔴 Errors & Panics | 即時錯誤監控 | `{container=~"tabelogo-.*"} \|~ "(?i)error\|panic\|fatal"` |
| 2 | 📊 Log Volume | 日誌量時序圖 | `sum by (service) (count_over_time({...}[$__interval]))` |
| 3 | 🔐 Auth Service | 認證服務日誌 | `{service="auth-service"}` |
| 4 | 🍽️ Restaurant Service | 餐廳服務日誌 | `{service="restaurant-service"}` |
| 5 | 🗺️ Map Service | 地圖服務日誌 | `{service="map-service"}` |
| 6 | 🕷️ Spider Service | 爬蟲服務日誌 | `{service="spider-service"}` |
| 7 | 🐘 PostgreSQL | 資料庫日誌 | `{container=~"tabelogo-postgres-.*"}` |
| 8 | 📦 Redis | 快取日誌 | `{container="tabelogo-redis"}` |

### LogQL 語法說明

```logql
# 基本選擇器 - 使用標籤過濾
{container="tabelogo-auth-service"}

# 正則匹配多個容器
{container=~"tabelogo-.*"}

# 日誌內容過濾（管道符 |）
{service="auth-service"} |= "error"      # 包含 "error"
{service="auth-service"} != "health"     # 不包含 "health"
{service="auth-service"} |~ "(?i)error"  # 正則匹配（不分大小寫）

# 統計查詢
count_over_time({service="auth-service"}[5m])  # 5分鐘內日誌數量
rate({service="auth-service"}[1m])             # 每秒日誌速率
```

---

## 🔌 gRPC Keepalive 配置

### 問題描述
```
ERROR: [transport] Client received GoAway with error code ENHANCE_YOUR_CALM 
and debug data equal to ASCII "too_many_pings"
```

**原因**: Restaurant Service 客戶端每 30 秒發送 keepalive ping，但 Map Service 服務端未配置允許客戶端 ping。gRPC 預設最小間隔為 5 分鐘。

### 解決方案

#### 服務端配置（Map Service）
**文件**: `internal/map/interfaces/grpc/module.go`

```go
grpcServer := grpc.NewServer(
    // 強制策略 - 控制允許的客戶端行為
    grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
        MinTime:             10 * time.Second, // 允許最小 10 秒 ping 間隔
        PermitWithoutStream: true,             // 允許無活動流時 ping
    }),
    // 服務端參數
    grpc.KeepaliveParams(keepalive.ServerParameters{
        MaxConnectionIdle:     5 * time.Minute,  // 閒置 5 分鐘後關閉
        MaxConnectionAge:      30 * time.Minute, // 連線最長 30 分鐘
        MaxConnectionAgeGrace: 10 * time.Second, // 關閉前的優雅等待
        Time:                  1 * time.Minute,  // 服務端每分鐘 ping
        Timeout:               20 * time.Second, // ping 確認超時
    }),
)
```

#### 客戶端配置（Restaurant Service）
**文件**: `internal/restaurant/infrastructure/module.go`

```go
grpcConfig := &grpc.ConnectionConfig{
    KeepAliveTime:    60 * time.Second, // 每 60 秒 ping（> MinTime 10s）
    KeepAliveTimeout: 20 * time.Second, // 與服務端一致
}
```

### 生產環境考量

| 參數 | 值 | 說明 |
|------|------|------|
| `MaxConnectionAge` | 30m | 連線定期輪替，便於負載均衡器分散請求 |
| `MaxConnectionIdle` | 5m | 及早釋放閒置連線，節省資源 |
| `PermitWithoutStream` | true | 允許長時間無請求時保持連線活躍 |
| Client `Time` | 60s | 每分鐘 ping 確保連線存活（適用於 AWS ALB/NLB） |

---

## 🚀 驗證命令

```bash
# 驗證 Loki 狀態
curl http://localhost:3100/ready

# 查看收集的容器
curl -s http://localhost:3100/loki/api/v1/label/container/values | jq

# 查詢日誌
curl -s 'http://localhost:3100/loki/api/v1/query_range' \
  --data-urlencode 'query={service="auth-service"}' \
  --data-urlencode 'limit=10' | jq

# 查看 Promtail 狀態
docker logs tabelogo-promtail --tail 10

# 檢查服務錯誤
docker logs tabelogo-restaurant-service --since 5m 2>&1 | grep -i error
```

---

## 📚 參考資源

- [Loki 官方文檔](https://grafana.com/docs/loki/latest/)
- [Promtail 配置](https://grafana.com/docs/loki/latest/send-data/promtail/)
- [LogQL 查詢語法](https://grafana.com/docs/loki/latest/query/)
- [gRPC Keepalive](https://grpc.io/docs/guides/keepalive/)
