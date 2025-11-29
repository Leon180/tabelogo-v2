# Map Service

Map Service 作為 Google Maps API 的後端代理層，提供安全、高效的地圖服務整合。

## 功能特性

- ✅ **Quick Search**: 單一餐廳查詢（Phase 2）
- ✅ **Advance Search**: 高級搜索功能（Phase 3）
- ✅ **Redis 緩存**: 減少 API 調用成本
- ✅ **API Key 保護**: 後端代理，前端不暴露 API Key

## 技術棧

- **語言**: Go 1.24
- **框架**: Gin (HTTP), Uber FX (DI)
- **緩存**: Redis 7
- **架構**: DDD 分層架構

## 快速開始

### 本地開發

1. **複製環境變數**
```bash
cd cmd/map-service
cp .env.example .env
```

2. **啟動 Redis**
```bash
docker-compose -f deployments/docker-compose/map-service.yml up -d redis-map
```

3. **運行服務**
```bash
cd cmd/map-service
go run main.go
```

4. **測試健康檢查**
```bash
curl http://localhost:8081/health
```

### Docker 部署

```bash
# 啟動完整服務（Redis + Map Service）
docker-compose -f deployments/docker-compose/map-service.yml up --build

# 查看日誌
docker-compose -f deployments/docker-compose/map-service.yml logs -f map-service

# 停止服務
docker-compose -f deployments/docker-compose/map-service.yml down
```

## 環境變數

| 變數 | 說明 | 預設值 |
|------|------|--------|
| `PORT` | HTTP 服務端口 | 8081 |
| `REDIS_HOST` | Redis 主機 | localhost |
| `REDIS_PORT` | Redis 端口 | 6380 |
| `REDIS_DB` | Redis 資料庫編號 | 5 |
| `GOOGLE_MAPS_API_KEY` | Google Maps API Key | - |
| `LOG_LEVEL` | 日誌等級 | info |

## API 端點

### Health Check
```bash
GET /health
```

### Quick Search (Phase 2)
```bash
POST /api/v1/map/quick_search
Content-Type: application/json

{
  "place_id": "ChIJN1t_tDeuEmsRUsoyG83frY4",
  "language_code": "ja"
}
```

### Advance Search (Phase 3)
```bash
POST /api/v1/map/advance_search
Content-Type: application/json

{
  "text_query": "寿司 東京",
  "location_bias": {
    "rectangle": {
      "low": {"latitude": 35.6, "longitude": 139.6},
      "high": {"latitude": 35.7, "longitude": 139.8}
    }
  },
  "max_result_count": 20,
  "language_code": "ja"
}
```

## 開發指令

```bash
# 安裝依賴
go mod tidy

# 運行測試
go test ./...

# 構建
go build -o map-service

# 格式化代碼
go fmt ./...

# Lint 檢查
golangci-lint run
```

## 專案結構

```
cmd/map-service/
├── main.go              # 入口點
├── go.mod               # 依賴管理
├── .env.example         # 環境變數範例
├── Dockerfile           # Docker 配置
└── README.md            # 本文檔

internal/map/
├── domain/              # 領域層
│   └── models/          # 數據模型
├── application/         # 應用層
│   └── usecases/        # 用例
├── interfaces/          # 接口層
│   └── http/            # HTTP 處理器
├── infrastructure/      # 基礎設施層
│   ├── cache/           # Redis 緩存
│   └── external/        # Google API 客戶端
└── module.go            # FX 模組定義
```

## 實現階段

- ✅ **Phase 1**: 基礎架構（已完成）
- ⏳ **Phase 2**: Quick Search 實現
- ⏳ **Phase 3**: Advance Search 實現
- ⏳ **Phase 4**: 優化與監控
- ⏳ **Phase 5**: 生產部署

## 相關文檔

- [MAP_SERVICE_DESIGN.md](../../MAP_SERVICE_DESIGN.md) - 完整設計文檔
- [architecture.md](../../architecture.md) - 整體架構說明
- [pkg/SHARED_PACKAGES.md](../../pkg/SHARED_PACKAGES.md) - 共用套件文檔

## License

MIT
