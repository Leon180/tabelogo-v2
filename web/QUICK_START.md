# Quick Start: 設置 Google Maps API Key

## 最快速的測試方式 (5 分鐘)

如果您只是想快速測試地圖功能,可以先建立一個簡單的 API Key:

### 步驟 1: 建立 API Key

1. 前往 https://console.cloud.google.com/apis/credentials
2. 點擊 **+ CREATE CREDENTIALS** → **API key**
3. 複製生成的 Key

### 步驟 2: 啟用 API

前往 https://console.cloud.google.com/apis/library 並啟用:
- Maps JavaScript API

### 步驟 3: 設定環境變數

在 `web/` 目錄創建 `.env.local` 文件:

```bash
NEXT_PUBLIC_GOOGLE_MAPS_API_KEY=你的API_Key
NEXT_PUBLIC_DEFAULT_LAT=35.6762
NEXT_PUBLIC_DEFAULT_LNG=139.6503
```

### 步驟 4: 重啟服務器

```bash
# 停止當前服務器 (Ctrl+C)
npm run dev
```

### 步驟 5: 測試

打開 http://localhost:3000 - 地圖應該會正常顯示!

---

## 生產環境設置

完整的安全設置請參考 [GOOGLE_MAPS_SETUP.md](file:///Users/lileon/goproject/tabelogov2/web/GOOGLE_MAPS_SETUP.md)
