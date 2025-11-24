# Google Maps API Key 設置指南 (混合模式)

## 步驟 1: 建立 Google Cloud 專案

1. 前往 [Google Cloud Console](https://console.cloud.google.com/)
2. 點擊專案選擇器,建立新專案或選擇現有專案
3. 專案名稱建議: `tabelogo-v2`

## 步驟 2: 啟用必要的 API

在 Google Cloud Console 中:

1. 前往 **APIs & Services** → **Library**
2. 搜尋並啟用以下 API:
   - ✅ **Maps JavaScript API** (用於前端地圖顯示)
   - ✅ **Places API (New)** (用於搜尋和詳細資訊)
   - ✅ **Geocoding API** (用於地址轉換)

## 步驟 3: 建立客戶端 API Key (受限)

### 3.1 建立 Key

1. 前往 **APIs & Services** → **Credentials**
2. 點擊 **+ CREATE CREDENTIALS** → **API key**
3. 複製生成的 API key (暫時保存) 

### 3.2 設定限制

點擊剛建立的 API key 進行編輯:

**名稱**: `Tabelogo-Client-Key`

**Application restrictions**:
- 選擇 **HTTP referrers (web sites)**
- 新增以下 referrers:
  ```
  http://localhost:3000/*
  http://127.0.0.1:3000/*
  https://yourdomain.com/*  (生產環境,稍後設定)
  ```

**API restrictions**:
- 選擇 **Restrict key**
- 只勾選:
  - ✅ Maps JavaScript API

點擊 **SAVE**

## 步驟 4: 建立服務器 API Key (IP 限制)

### 4.1 建立 Key

1. 再次點擊 **+ CREATE CREDENTIALS** → **API key**
2. 複製生成的 API key: AIzaSyBTB8hEYIZ1JCFAg4bV7B6tN2V9ENidO2o

### 4.2 設定限制

**名稱**: `Tabelogo-Server-Key`

**Application restrictions**:
- 選擇 **IP addresses (web servers, cron jobs, etc.)**
- 新增以下 IP (開發環境):
  ```
  127.0.0.1
  ::1
  ```
  > 生產環境時需要加入實際服務器的 IP

**API restrictions**:
- 選擇 **Restrict key**
- 勾選:
  - ✅ Places API (New)
  - ✅ Geocoding API

點擊 **SAVE**

## 步驟 5: 配置前端環境變數

在 `web/` 目錄下創建 `.env.local` 文件:

```bash
# 客戶端 API Key (僅用於地圖顯示)
NEXT_PUBLIC_GOOGLE_MAPS_API_KEY=AIza...你的客戶端Key

# 後端服務 URL
NEXT_PUBLIC_MAP_SERVICE_URL=http://localhost:8080
NEXT_PUBLIC_AUTH_SERVICE_URL=http://localhost:8081
NEXT_PUBLIC_RESTAURANT_SERVICE_URL=http://localhost:8082
NEXT_PUBLIC_BOOKING_SERVICE_URL=http://localhost:8083

# 預設地圖中心 (東京)
NEXT_PUBLIC_DEFAULT_LAT=35.6762
NEXT_PUBLIC_DEFAULT_LNG=139.6503
```

## 步驟 6: 配置後端環境變數

在 `google-map/` 或 `map-service/` 目錄下的 `.env` 文件:

```bash
# 服務器 API Key (用於 Places API)
GOOGLE_MAPS_API_KEY=AIza...你的服務器Key

# 服務端口
PORT=8080
```

## 步驟 7: 測試設置

### 7.1 測試前端地圖顯示

```bash
cd web
npm run dev
```

打開 http://localhost:3000
- 應該能看到 Google Maps 正常載入
- 地圖上應該顯示模擬的餐廳標記

### 7.2 測試後端 API (稍後)

當 Map Service 啟動後:
```bash
curl -X POST http://localhost:8080/advance_search \
  -H "Content-Type: application/json" \
  -d '{
    "text_query": "sushi in tokyo",
    "low_latitude": 35.6,
    "low_longitude": 139.6,
    "high_latitude": 35.7,
    "high_longitude": 139.8,
    "max_result_count": 10,
    "min_rating": 4,
    "open_now": false,
    "rank_preference": "RELEVANCE",
    "language_code": "en"
  }'
```

## 安全檢查清單

- [ ] 客戶端 Key 已設定 HTTP referrer 限制
- [ ] 客戶端 Key 只啟用 Maps JavaScript API
- [ ] 服務器 Key 已設定 IP 限制
- [ ] 服務器 Key 只啟用必要的 API
- [ ] `.env.local` 已加入 `.gitignore`
- [ ] 生產環境使用不同的 API Key

## 費用管理

### 設定配額警報

1. 前往 **APIs & Services** → **Quotas**
2. 選擇 **Maps JavaScript API**
3. 點擊 **EDIT QUOTAS**
4. 設定每日請求上限 (建議: 25,000 次/天)

### 設定計費警報

1. 前往 **Billing** → **Budgets & alerts**
2. 建立預算警報
3. 設定閾值 (例如: $50/月)

## 常見問題

### Q: 地圖顯示空白或錯誤?
A: 檢查:
1. API Key 是否正確複製到 `.env.local`
2. HTTP referrer 是否包含 `http://localhost:3000/*`
3. Maps JavaScript API 是否已啟用
4. 瀏覽器 Console 是否有錯誤訊息

### Q: 搜尋功能不工作?
A: 確認:
1. Map Service 是否正在運行
2. 服務器 API Key 是否正確設定
3. Places API 是否已啟用

### Q: 如何切換到生產環境?
A: 
1. 建立新的 API Key 組
2. 更新 HTTP referrer 為生產域名
3. 更新 IP 限制為生產服務器 IP
4. 使用環境變數管理不同環境的 Key

## 下一步

完成設置後:
1. 重啟開發服務器: `npm run dev`
2. 打開瀏覽器測試地圖是否正常顯示
3. 如有問題,查看瀏覽器 Console 的錯誤訊息
