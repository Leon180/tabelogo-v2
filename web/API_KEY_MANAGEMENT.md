# Google Maps API Key 管理策略

## 當前狀況

前端已成功啟動,但地圖區域顯示錯誤訊息:
> "Google Maps API key not configured. Please check ENV_CONFIG.md"

![當前頁面狀態](file:///Users/lileon/.gemini/antigravity/brain/ff993abe-b807-4d40-baf6-c790862aebc3/map_interface_no_key_1763979022196.png)

## API Key 安全管理方案

### 方案 1: 後端代理模式 (推薦 ⭐)

**架構**:
```
Frontend → Map Service (Backend) → Google Maps API
                ↑
            API Key 存儲在後端
```

**優點**:
- ✅ API Key 完全不暴露給客戶端
- ✅ 可以在後端實施額外的驗證和限流
- ✅ 可以記錄所有 API 使用情況
- ✅ 符合原版 tabelogo 的架構設計

**實作方式**:
1. 將 Google Maps API Key 存儲在 `google-map` 服務的環境變數中
2. 前端只調用我們的 `map-service` API
3. `map-service` 使用後端的 API Key 調用 Google Maps API
4. 前端地圖顯示使用公開的 Map ID (不需要 API Key)

**需要修改**:
- 前端 `GoogleMap` 組件改為使用 Map ID 而非 API Key
- 確保 `map-service` 已經實作並運行

### 方案 2: 受限的客戶端 API Key

**架構**:
```
Frontend (with restricted API Key) → Google Maps API
```

**限制設定** (在 Google Cloud Console):
1. **HTTP referrers 限制**:
   - `http://localhost:3000/*` (開發環境)
   - `https://yourdomain.com/*` (生產環境)

2. **API 限制**:
   - 只啟用 Maps JavaScript API
   - 不啟用 Places API (改用後端代理)

3. **配額限制**:
   - 設定每日/每月使用上限

**優點**:
- ✅ 實作簡單
- ✅ 地圖載入速度快

**缺點**:
- ❌ API Key 仍然暴露在客戶端程式碼中
- ❌ 有被濫用的風險(即使有限制)

### 方案 3: 混合模式 (最佳平衡 ⭐⭐)

**架構**:
```
Frontend:
  - Map Display: 使用受限的客戶端 API Key
  - Search/Places API: 調用後端 Map Service

Backend (Map Service):
  - 使用不受限的 Server API Key
  - 處理所有敏感的 API 調用
```

**實作步驟**:

1. **建立兩個 API Key**:
   - **Client Key**: 僅用於地圖顯示,高度受限
   - **Server Key**: 用於後端服務,IP 限制

2. **前端配置** (`.env.local`):
   ```bash
   # 受限的客戶端 Key (僅用於地圖顯示)
   NEXT_PUBLIC_GOOGLE_MAPS_API_KEY=AIza...client_key
   
   # 後端服務 URL
   NEXT_PUBLIC_MAP_SERVICE_URL=http://localhost:8080
   ```

3. **後端配置** (Map Service):
   ```bash
   # 不受限的服務器 Key (用於 Places API 等)
   GOOGLE_MAPS_API_KEY=AIza...server_key
   ```

## 推薦方案

**建議採用方案 3 (混合模式)**,原因:
1. 符合微服務架構設計
2. 平衡了安全性和效能
3. 客戶端 Key 即使被提取也無法濫用
4. 敏感操作(搜尋、詳細資訊)都通過後端

## 實作檢查清單

### 立即執行 (開發環境測試)
- [ ] 在 Google Cloud Console 建立受限的客戶端 API Key
- [ ] 設定 HTTP referrer 為 `http://localhost:3000/*`
- [ ] 只啟用 Maps JavaScript API
- [ ] 將 Key 加入 `web/.env.local`
- [ ] 測試地圖是否正常顯示

### 後續執行 (生產環境準備)
- [ ] 建立獨立的 Server API Key 給 Map Service
- [ ] 設定 IP 限制
- [ ] 啟動並配置 Map Service
- [ ] 修改前端搜尋功能改為調用 Map Service
- [ ] 設定 API 使用配額警報

## 下一步

請告訴我您想採用哪個方案,我可以協助:
1. 建立 API Key 的詳細步驟指南
2. 修改程式碼以實作選定的方案
3. 設定環境變數並測試
