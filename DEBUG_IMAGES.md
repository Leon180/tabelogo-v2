# 圖片顯示問題診斷指南

## 請在瀏覽器中執行以下步驟：

### 步驟 1: 打開餐廳詳情
1. 訪問 http://localhost:3000
2. 搜索並點擊任意餐廳
3. 打開瀏覽器開發者工具 (F12)

### 步驟 2: 查看控制台輸出
在 Console 標籤中，查找 `🖼️ Place data:` 的輸出

**請完整複製整個對象**，包括：
```javascript
{
  hasRestaurantData: true,  // 或 false
  hasMapData: true,         // 或 false
  photos: [...],            // 數組內容
  photosLength: 5,          // 數字
  firstPhoto: {...}         // 對象內容
}
```

### 步驟 3: 檢查 Network 標籤
1. 切換到 Network 標籤
2. 過濾 "quick-search"
3. 查看請求：
   - 是否有 `/api/v1/restaurants/quick-search/` 請求？
   - 是否有 `/api/v1/map/quick-search` 請求？
   - 響應狀態碼是什麼？(200, 404, 500?)

### 步驟 4: 查看圖片請求
1. 在 Network 標籤中過濾 "places.googleapis.com"
2. 是否有圖片請求？
3. 如果有，狀態碼是什麼？

### 步驟 5: 檢查錯誤
在 Console 標籤中，是否有任何紅色錯誤信息？

---

## 可能的問題和解決方案

### 情況 A: photos 是 undefined
**原因**: Map Service 沒有被調用或沒有返回 photos
**解決**: 需要確保 Map Service fallback 正常工作

### 情況 B: photos 是空數組 []
**原因**: Google Places API 沒有返回照片
**解決**: 這是正常的，某些地點可能沒有照片

### 情況 C: photos 有數據但圖片不顯示
**原因**: Next.js 圖片配置或 API key 問題
**解決**: 檢查 next.config.ts 和 GOOGLE_MAPS_API_KEY

### 情況 D: Network 錯誤
**原因**: CORS 或認證問題
**解決**: 檢查請求 headers 和響應

---

## 請提供以下信息：

1. **控制台完整輸出** (🖼️ Place data: 後面的完整對象)
2. **Network 請求列表** (quick-search 相關的所有請求)
3. **任何錯誤信息** (紅色的錯誤)
4. **圖片區域顯示什麼？** ("No Image Available" 還是完全空白？)

有了這些信息，我就能準確診斷並修復問題！
