# Map Service - 完整設計文檔

**版本：** v2.0
**日期：** 2025-11-25
**狀態：** 📋 設計階段

---

## 📋 目錄

1. [概述](#概述)
2. [核心功能](#核心功能)
3. [API 設計](#api-設計)
4. [數據模型](#數據模型)
5. [Google Maps API 整合](#google-maps-api-整合)
6. [緩存策略](#緩存策略)
7. [實現階段](#實現階段)
8. [技術棧](#技術棧)
9. [安全與配額管理](#安全與配額管理)
10. [測試策略](#測試策略)

---

## 概述

### 目的

Map Service 作為 Google Maps API 的後端代理層，提供：
1. **安全性**：API Key 不暴露給前端
2. **成本控制**：實現緩存減少 API 調用
3. **統一接口**：標準化的 REST API
4. **數據增強**：整合 Tabelog 數據

### 架構位置

```
Frontend (Next.js)
    ↓ HTTP
API Gateway (未來)
    ↓ HTTP
Map Service ← → Redis (緩存)
    ↓ HTTPS
Google Maps API
```

### 第一版 vs 第二版差異

| 特性 | 第一版 (V1) | 第二版 (V2) |
|------|------------|------------|
| **架構** | 單體服務 | 微服務架構 |
| **認證緩存** | 調用 auth service | 使用 Redis |
| **錯誤處理** | 基礎錯誤返回 | 結構化錯誤 + 重試機制 |
| **日誌** | 基本 println | 結構化日誌 (Uber Zap) |
| **依賴注入** | 無 | Uber FX |
| **配置管理** | 環境變數 | Viper + 環境變數 |
| **數據驗證** | Gin binding | 完整驗證 + 自定義錯誤 |
| **監控** | 無 | Prometheus metrics |

---

## 核心功能

### 1. Quick Search（單一餐廳查詢）

**用途：** 用戶點擊地圖標記時獲取餐廳詳情

**流程：**
```
1. 前端發送 place_id
2. 檢查 Redis 緩存
   ├─ 命中 → 返回緩存數據
   └─ 未命中 ↓
3. 調用 Google Places API (Place Details)
4. 存入 Redis (TTL: 1小時)
5. 返回結果給前端
```

**特點：**
- 支援多語言（en, ja, zh-TW）
- 可自定義返回欄位（field mask）
- 緩存策略：1小時過期

### 2. Advance Search（高級搜索）

**用途：** 用戶通過搜索框查找餐廳

**流程：**
```
1. 前端發送搜索條件（文字、位置、過濾器）
2. 構建 Google Places API 請求
3. 調用 Text Search API
4. 過濾結果（評分、營業狀態等）
5. 返回最多 20 條結果
```

**特點：**
- 地理位置矩形範圍搜索
- 支援過濾：最低評分、營業中、排序偏好
- 不緩存（搜索條件變化太大）

### 3. Geocoding（地址轉換）- 未來功能

**用途：** 地址 ↔ 座標轉換

### 4. Distance Matrix（距離計算）- 未來功能

**用途：** 計算用戶到餐廳的距離和時間

---

## API 設計

### Base URL

```
Development: http://localhost:8081
Production:  https://api.tabelogo.com/map
```

### 1. Quick Search API

**Endpoint:** `POST /api/v1/map/quick_search`

**Request:**
```json
{
  "place_id": "ChIJN1t_tDeuEmsRUsoyG83frY4",
  "language_code": "ja",
  "api_mask": "id,displayName,formattedAddress,location,rating,priceLevel,photos"
}
```

**Response (Success):**
```json
{
  "source": "redis",  // or "google"
  "cached_at": "2025-11-25T10:00:00Z",
  "result": {
    "id": "ChIJN1t_tDeuEmsRUsoyG83frY4",
    "displayName": {
      "text": "すきやばし次郎",
      "languageCode": "ja"
    },
    "formattedAddress": "東京都中央区銀座4-2-15",
    "location": {
      "latitude": 35.6708,
      "longitude": 139.7634
    },
    "rating": 4.8,
    "priceLevel": "PRICE_LEVEL_VERY_EXPENSIVE",
    "photos": [...]
  }
}
```

**Response (Error):**
```json
{
  "error": "place_not_found",
  "message": "Place with ID 'ChIJxxx' not found",
  "code": 404,
  "timestamp": "2025-11-25T10:00:00Z"
}
```

**Request Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| place_id | string | ✅ | Google Place ID |
| language_code | string | ✅ | ISO 639-1 (en, ja, zh-TW) |
| api_mask | string | ❌ | 返回欄位列表 (default: 基本欄位) |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| source | string | "redis" 或 "google" |
| cached_at | string | 緩存時間 (ISO 8601) |
| result | object | Place Details |

### 2. Advance Search API

**Endpoint:** `POST /api/v1/map/advance_search`

**Request:**
```json
{
  "text_query": "寿司 東京",
  "location_bias": {
    "rectangle": {
      "low": {
        "latitude": 35.6,
        "longitude": 139.6
      },
      "high": {
        "latitude": 35.7,
        "longitude": 139.8
      }
    }
  },
  "max_result_count": 20,
  "min_rating": 4.0,
  "open_now": true,
  "rank_preference": "DISTANCE",
  "language_code": "ja",
  "api_mask": "places.id,places.displayName,places.formattedAddress,places.location,places.rating"
}
```

**Response:**
```json
{
  "places": [
    {
      "id": "ChIJ...",
      "displayName": {...},
      "formattedAddress": "...",
      "location": {...},
      "rating": 4.8
    }
  ],
  "total_count": 15,
  "search_metadata": {
    "text_query": "寿司 東京",
    "search_time_ms": 234
  }
}
```

**Request Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| text_query | string | ✅ | 搜索文字 |
| location_bias.rectangle | object | ✅ | 搜索範圍矩形 |
| max_result_count | int | ✅ | 最多返回數量 (1-20) |
| min_rating | float | ❌ | 最低評分 (0-5) |
| open_now | bool | ❌ | 只顯示營業中 |
| rank_preference | string | ✅ | DISTANCE 或 RELEVANCE |
| language_code | string | ✅ | ISO 639-1 |
| api_mask | string | ❌ | 返回欄位 |

### 3. Health Check API

**Endpoint:** `GET /health`

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-11-25T10:00:00Z",
  "version": "2.0.0",
  "dependencies": {
    "redis": "connected",
    "google_api": "reachable"
  }
}
```

---

## 數據模型

### Go Structures

```go
// internal/map/domain/models/place.go

package models

import "time"

// QuickSearchRequest - 快速搜索請求
type QuickSearchRequest struct {
    PlaceID      string `json:"place_id" binding:"required"`
    LanguageCode string `json:"language_code" binding:"required,oneof=en ja zh-TW"`
    APIMask      string `json:"api_mask"`
}

// QuickSearchResponse - 快速搜索響應
type QuickSearchResponse struct {
    Source   string      `json:"source"`
    CachedAt *time.Time  `json:"cached_at,omitempty"`
    Result   interface{} `json:"result"`
}

// AdvanceSearchRequest - 高級搜索請求
type AdvanceSearchRequest struct {
    TextQuery      string         `json:"text_query" binding:"required"`
    LocationBias   LocationBias   `json:"location_bias" binding:"required"`
    MaxResultCount int            `json:"max_result_count" binding:"required,min=1,max=20"`
    MinRating      float64        `json:"min_rating" binding:"omitempty,min=0,max=5"`
    OpenNow        bool           `json:"open_now"`
    RankPreference string         `json:"rank_preference" binding:"required,oneof=DISTANCE RELEVANCE"`
    LanguageCode   string         `json:"language_code" binding:"required,oneof=en ja zh-TW"`
    APIMask        string         `json:"api_mask"`
}

// LocationBias - 位置偏好
type LocationBias struct {
    Rectangle Rectangle `json:"rectangle" binding:"required"`
}

// Rectangle - 矩形範圍
type Rectangle struct {
    Low  Coordinates `json:"low" binding:"required"`
    High Coordinates `json:"high" binding:"required"`
}

// Coordinates - 座標
type Coordinates struct {
    Latitude  float64 `json:"latitude" binding:"required,min=-90,max=90"`
    Longitude float64 `json:"longitude" binding:"required,min=-180,max=180"`
}

// AdvanceSearchResponse - 高級搜索響應
type AdvanceSearchResponse struct {
    Places         []interface{}  `json:"places"`
    TotalCount     int            `json:"total_count"`
    SearchMetadata SearchMetadata `json:"search_metadata"`
}

// SearchMetadata - 搜索元數據
type SearchMetadata struct {
    TextQuery     string `json:"text_query"`
    SearchTimeMs  int64  `json:"search_time_ms"`
}

// ErrorResponse - 錯誤響應
type ErrorResponse struct {
    Error     string    `json:"error"`
    Message   string    `json:"message"`
    Code      int       `json:"code"`
    Timestamp time.Time `json:"timestamp"`
}
```

---

## Google Maps API 整合

### Places API (New)

我們使用最新的 **Places API (New)** 而非舊版 Places API。

**官方文檔：** https://developers.google.com/maps/documentation/places/web-service/overview

### 1. Place Details

**用於：** Quick Search

**Endpoint:**
```
GET https://places.googleapis.com/v1/places/{PLACE_ID}
```

**Headers:**
```
X-Goog-Api-Key: YOUR_API_KEY
X-Goog-FieldMask: id,displayName,formattedAddress,location,rating
```

**Go 實現示例：**
```go
func (s *PlacesService) GetPlaceDetails(placeID string, languageCode string, fieldMask string) (interface{}, error) {
    url := fmt.Sprintf("https://places.googleapis.com/v1/places/%s", placeID)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    // Add query parameters
    q := req.URL.Query()
    if fieldMask != "" {
        q.Add("fields", fieldMask)
    }
    q.Add("languageCode", languageCode)
    req.URL.RawQuery = q.Encode()

    // Add headers
    req.Header.Set("X-Goog-Api-Key", s.config.GoogleMapsAPIKey)

    // Execute request
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Parse response
    var result interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return result, nil
}
```

### 2. Text Search

**用於：** Advance Search

**Endpoint:**
```
POST https://places.googleapis.com/v1/places:searchText
```

**Headers:**
```
Content-Type: application/json
X-Goog-Api-Key: YOUR_API_KEY
X-Goog-FieldMask: places.id,places.displayName,places.formattedAddress
```

**Request Body:**
```json
{
  "textQuery": "寿司 東京",
  "locationBias": {
    "rectangle": {
      "low": {"latitude": 35.6, "longitude": 139.6},
      "high": {"latitude": 35.7, "longitude": 139.8}
    }
  },
  "maxResultCount": 20,
  "minRating": 4.0,
  "openNow": true,
  "rankPreference": "DISTANCE",
  "languageCode": "ja"
}
```

**Go 實現示例：**
```go
func (s *PlacesService) SearchText(req *models.AdvanceSearchRequest) (*models.AdvanceSearchResponse, error) {
    // Build request body
    body := map[string]interface{}{
        "textQuery": req.TextQuery,
        "locationBias": map[string]interface{}{
            "rectangle": map[string]interface{}{
                "low":  map[string]float64{"latitude": req.LocationBias.Rectangle.Low.Latitude, "longitude": req.LocationBias.Rectangle.Low.Longitude},
                "high": map[string]float64{"latitude": req.LocationBias.Rectangle.High.Latitude, "longitude": req.LocationBias.Rectangle.High.Longitude},
            },
        },
        "maxResultCount": req.MaxResultCount,
        "minRating":      req.MinRating,
        "openNow":        req.OpenNow,
        "rankPreference": req.RankPreference,
        "languageCode":   req.LanguageCode,
    }

    jsonBody, err := json.Marshal(body)
    if err != nil {
        return nil, err
    }

    // Create HTTP request
    httpReq, err := http.NewRequest("POST", "https://places.googleapis.com/v1/places:searchText", bytes.NewBuffer(jsonBody))
    if err != nil {
        return nil, err
    }

    // Add headers
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("X-Goog-Api-Key", s.config.GoogleMapsAPIKey)
    httpReq.Header.Set("X-Goog-FieldMask", req.APIMask)

    // Execute
    client := &http.Client{Timeout: 15 * time.Second}
    resp, err := client.Do(httpReq)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Parse
    var result struct {
        Places []interface{} `json:"places"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return &models.AdvanceSearchResponse{
        Places:     result.Places,
        TotalCount: len(result.Places),
    }, nil
}
```

### API 配額管理

**每日免費配額（估算）：**
- Place Details: $0.017 per request
- Text Search: $0.032 per request
- 每月 $200 免費額度

**配額計算：**
```
假設每日 1000 次請求：
- Quick Search (70%): 700 × $0.017 = $11.90
- Advance Search (30%): 300 × $0.032 = $9.60
每日成本: $21.50
每月成本: $645 (超出免費額度)

使用緩存後 (80% 命中率):
- Quick Search: 140 × $0.017 = $2.38
- Advance Search: 300 × $0.032 = $9.60
每日成本: $11.98
每月成本: $359.40
```

**限制策略：**
1. Redis 緩存（Quick Search 1小時 TTL）
2. Rate Limiting (每用戶每分鐘最多 10 次請求)
3. 監控告警（接近配額時通知）

---

## 緩存策略

### Redis 緩存設計

**Key 格式：**
```
map:place:{place_id}:{language_code}
```

**TTL 策略：**
- Place Details: 1 小時
- 熱門地點可延長至 24 小時

**緩存數據結構：**
```json
{
  "place_id": "ChIJ...",
  "language_code": "ja",
  "data": {...},  // Google API 原始響應
  "cached_at": "2025-11-25T10:00:00Z",
  "expires_at": "2025-11-25T11:00:00Z"
}
```

**實現示例：**
```go
func (s *CacheService) GetPlace(placeID, languageCode string) (*CachedPlace, error) {
    key := fmt.Sprintf("map:place:%s:%s", placeID, languageCode)

    val, err := s.redis.Get(context.Background(), key).Result()
    if err == redis.Nil {
        return nil, nil // Cache miss
    } else if err != nil {
        return nil, err
    }

    var cached CachedPlace
    if err := json.Unmarshal([]byte(val), &cached); err != nil {
        return nil, err
    }

    return &cached, nil
}

func (s *CacheService) SetPlace(placeID, languageCode string, data interface{}, ttl time.Duration) error {
    key := fmt.Sprintf("map:place:%s:%s", placeID, languageCode)

    cached := CachedPlace{
        PlaceID:      placeID,
        LanguageCode: languageCode,
        Data:         data,
        CachedAt:     time.Now(),
        ExpiresAt:    time.Now().Add(ttl),
    }

    jsonData, err := json.Marshal(cached)
    if err != nil {
        return err
    }

    return s.redis.Set(context.Background(), key, jsonData, ttl).Err()
}
```

### 緩存失效策略

1. **TTL 自動過期**：1小時後自動刪除
2. **主動失效**：餐廳資訊更新時（未來功能）
3. **LRU 淘汰**：Redis 記憶體不足時自動淘汰

---

## 實現階段

### Phase 1: 基礎架構（優先）⭐

**目標：** 建立可運行的 Map Service 骨架

**任務清單：**
- [ ] 創建項目結構（DDD 分層）
- [ ] 設置 Uber FX 依賴注入
- [ ] 配置管理（Viper + 環境變數）
- [ ] 健康檢查 API
- [ ] Docker 配置

**預計時間：** 2-3 小時

**交付物：**
- 可啟動的 HTTP 服務
- Health check 端點
- 基本日誌

### Phase 2: Quick Search 實現

**目標：** 實現單一餐廳查詢功能

**任務清單：**
- [ ] Google Places API 客戶端
- [ ] Place Details API 整合
- [ ] Redis 緩存層
- [ ] Quick Search API 端點
- [ ] 錯誤處理

**預計時間：** 3-4 小時

**交付物：**
- `/api/v1/map/quick_search` 端點
- 緩存機制
- 單元測試

**測試計劃：**
```bash
# 測試 Place Details (不會調用 API，使用 mock)
curl -X POST http://localhost:8081/api/v1/map/quick_search \
  -H "Content-Type: application/json" \
  -d '{
    "place_id": "ChIJN1t_tDeuEmsRUsoyG83frY4",
    "language_code": "ja"
  }'

# 實際調用 Google API 測試（需確認）
# 會消耗配額！
```

### Phase 3: Advance Search 實現

**目標：** 實現高級搜索功能

**任務清單：**
- [ ] Text Search API 整合
- [ ] 請求構建器
- [ ] 結果過濾
- [ ] Advance Search API 端點

**預計時間：** 3-4 小時

**交付物：**
- `/api/v1/map/advance_search` 端點
- 搜索過濾邏輯
- 集成測試

### Phase 4: 優化與監控

**目標：** 生產就緒

**任務清單：**
- [ ] Rate Limiting
- [ ] Prometheus metrics
- [ ] 結構化日誌（Zap）
- [ ] 錯誤重試機制
- [ ] API 配額監控

**預計時間：** 2-3 小時

**交付物：**
- 完整監控
- 告警機制
- 性能優化

### Phase 5: Docker & 部署

**目標：** 容器化部署

**任務清單：**
- [ ] Dockerfile
- [ ] Docker Compose 配置
- [ ] 與其他服務集成
- [ ] CI/CD 配置

**預計時間：** 2 小時

---

## 技術棧

### 核心依賴

```go
// go.mod (estimated)
module github.com/Leon180/tabelogo-v2/cmd/map-service

go 1.24

require (
    github.com/gin-gonic/gin v1.10.0
    github.com/redis/go-redis/v9 v9.6.1
    github.com/spf13/viper v1.19.0
    go.uber.org/fx v1.23.0
    go.uber.org/zap v1.27.0
)
```

### 項目結構（DDD）

```
cmd/map-service/
├── main.go                      # 入口點
├── go.mod
├── go.sum
├── .env                         # 環境變數
├── Dockerfile
└── README.md

internal/map/
├── domain/                      # 領域層
│   ├── models/                  # 數據模型
│   │   ├── place.go
│   │   └── search.go
│   └── services/                # 領域服務
│       └── places_service.go
│
├── application/                 # 應用層
│   └── usecases/
│       ├── quick_search.go
│       └── advance_search.go
│
├── interfaces/                  # 接口層
│   └── http/
│       ├── handler.go           # HTTP handlers
│       ├── routes.go            # 路由註冊
│       ├── dto.go               # 數據傳輸對象
│       └── middleware.go        # 中間件
│
└── infrastructure/              # 基礎設施層
    ├── cache/
    │   └── redis.go             # Redis 客戶端
    ├── external/
    │   └── google_places.go     # Google API 客戶端
    └── config/
        └── config.go            # 配置管理

pkg/
└── errors/                      # 共享錯誤處理
    └── errors.go
```

---

## 安全與配額管理

### 1. API Key 安全

**✅ 安全做法：**
- API Key 存在環境變數 `.env`
- 不提交到 Git（`.gitignore` 包含 `.env`）
- 後端代理，前端不直接訪問

**❌ 避免：**
- 硬編碼 API Key
- 前端直接調用 Google API
- 公開 API Key

### 2. Rate Limiting

**實現策略：**
```go
// 使用 Redis + Token Bucket 算法
// 每個用戶每分鐘最多 10 次請求

func (m *RateLimitMiddleware) LimitByUser() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := getUserID(c) // 從 token 獲取

        key := fmt.Sprintf("rate_limit:map:%s", userID)
        count, _ := m.redis.Incr(c, key).Result()

        if count == 1 {
            m.redis.Expire(c, key, time.Minute)
        }

        if count > 10 {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "rate_limit_exceeded",
                "message": "Too many requests. Please try again later."
            })
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### 3. 配額監控

**Prometheus Metrics：**
```go
var (
    googleAPICallsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "map_google_api_calls_total",
            Help: "Total number of Google API calls",
        },
        []string{"endpoint", "status"},
    )

    cacheHitRate = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "map_cache_hit_rate",
            Help: "Cache hit rate percentage",
        },
        []string{"cache_type"},
    )
)
```

---

## 測試策略

### 1. 單元測試

**測試範圍：**
- 數據模型驗證
- 緩存邏輯
- 請求構建器

**示例：**
```go
func TestQuickSearchRequest_Validate(t *testing.T) {
    tests := []struct {
        name    string
        req     QuickSearchRequest
        wantErr bool
    }{
        {
            name: "valid request",
            req: QuickSearchRequest{
                PlaceID:      "ChIJN1t_tDeuEmsRUsoyG83frY4",
                LanguageCode: "ja",
            },
            wantErr: false,
        },
        {
            name: "missing place_id",
            req: QuickSearchRequest{
                LanguageCode: "ja",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.req.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### 2. 集成測試

**使用 Mock Google API：**
```go
func TestQuickSearchHandler_Integration(t *testing.T) {
    // Setup mock server
    mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Return mock Google API response
        json.NewEncoder(w).Encode(map[string]interface{}{
            "id": "ChIJ...",
            "displayName": map[string]string{"text": "Test Restaurant"},
        })
    }))
    defer mockServer.Close()

    // Test handler with mock
    // ...
}
```

### 3. 手動測試清單

**Quick Search：**
- [ ] 有效 place_id
- [ ] 無效 place_id
- [ ] 不同語言（en, ja, zh-TW）
- [ ] 自定義 field mask
- [ ] 緩存命中
- [ ] 緩存未命中

**Advance Search：**
- [ ] 基本文字搜索
- [ ] 位置範圍搜索
- [ ] 最低評分過濾
- [ ] 營業中過濾
- [ ] 排序偏好（DISTANCE vs RELEVANCE）
- [ ] 結果數量限制

---

## 配置文件

### `.env.example`

```bash
# Server Configuration
PORT=8081
GIN_MODE=release

# Google Maps API
GOOGLE_MAPS_API_KEY=your_api_key_here

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=5
REDIS_PASSWORD=

# Cache Settings
CACHE_PLACE_DETAILS_TTL=3600  # 1 hour in seconds

# Rate Limiting
RATE_LIMIT_PER_MINUTE=10

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

---

## 下一步行動

### 立即開始（推薦）

1. **確認需求**
   - ✅ Google Maps API Key 已設置
   - ✅ Redis 可用（可用 docker-compose）
   - ⏸️ 等待你的確認再調用 Google API

2. **Phase 1 實現**
   - 創建項目結構
   - 設置基礎架構
   - 健康檢查 API

3. **逐步實現**
   - Phase 2: Quick Search
   - Phase 3: Advance Search

### 問題與決策

**需要你決定：**

1. **Google API 調用測試**
   - 何時可以開始調用真實 Google API？
   - 是否需要設置每日配額限制？

2. **優先級**
   - 先完整實現 Quick Search？
   - 還是同時實現兩個功能？

3. **額外功能**
   - 是否需要 Geocoding API？
   - 是否需要 Distance Matrix？

---

**準備好開始了嗎？** 🚀

請告訴我：
1. 你想從 Phase 1 開始嗎？
2. 何時可以測試調用 Google API？
3. 是否有其他需求或調整？
