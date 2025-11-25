# 測試登入功能指南

## 前置條件

1. **後端服務運行中** (Docker)
   ```bash
   docker compose -f deployments/docker-compose/auth-service.yml ps
   ```
   確認 `tabelogo-auth-service` 狀態為 `Up (healthy)`

2. **前端服務運行中**
   ```bash
   cd web
   npm run dev
   ```
   前端應該運行在 http://localhost:3000

## 手動測試步驟

### 1. 訪問登入頁面
在瀏覽器中打開：
```
http://localhost:3000/auth/login
```

### 2. 輸入測試帳號
- **Email**: `test@example.com`
- **Password**: `password123`

### 3. 點擊 "Sign In" 按鈕

### 4. 預期結果
- ✅ 不應該看到 CORS 錯誤
- ✅ 不應該看到 OPTIONS 404 錯誤
- ✅ 成功登入後應該跳轉到首頁
- ✅ 在瀏覽器 Console 中不應該有紅色錯誤
- ✅ localStorage 中應該存有 `access_token` 和 `refresh_token`

## 使用瀏覽器開發者工具檢查

### Chrome/Edge/Arc 開發者工具
1. 按 `F12` 或 `Cmd+Option+I` (Mac) 打開開發者工具
2. 切換到 **Network** 標籤
3. 勾選 "Preserve log"
4. 點擊登入按鈕
5. 查看網路請求：

**預期看到的請求：**

#### OPTIONS 預檢請求
```
Request URL: http://localhost:8080/api/v1/auth/login
Request Method: OPTIONS
Status Code: 204 No Content

Response Headers:
✓ Access-Control-Allow-Origin: *
✓ Access-Control-Allow-Methods: POST, OPTIONS, GET, PUT, DELETE, PATCH
✓ Access-Control-Allow-Headers: Content-Type, ...
✓ Access-Control-Allow-Credentials: true
```

#### POST 登入請求
```
Request URL: http://localhost:8080/api/v1/auth/login
Request Method: POST
Status Code: 200 OK

Response Headers:
✓ Access-Control-Allow-Origin: *
✓ Content-Type: application/json

Response Body:
{
  "access_token": "eyJ...",
  "refresh_token": "eyJ...",
  "user": {
    "id": "...",
    "email": "test@example.com",
    "username": "TestUser",
    ...
  }
}
```

### 檢查 localStorage
在開發者工具的 **Application** (或 **Storage**) 標籤：
1. 展開 **Local Storage**
2. 點擊 `http://localhost:3000`
3. 應該看到：
   - `access_token`: `eyJ...`
   - `refresh_token`: `eyJ...`

## 常見問題排查

### 問題 1: CORS 錯誤
**錯誤訊息：**
```
Access to XMLHttpRequest at 'http://localhost:8080/api/v1/auth/login'
from origin 'http://localhost:3000' has been blocked by CORS policy
```

**解決方案：**
```bash
# 重新構建並重啟 auth-service
cd /Users/lileon/goproject/tabelogov2
docker compose -f deployments/docker-compose/auth-service.yml down
docker compose -f deployments/docker-compose/auth-service.yml up -d --build
```

### 問題 2: OPTIONS 404
**錯誤訊息：** `404 Not Found` on OPTIONS request

**解決方案：**
確認 Docker 容器使用了最新的代碼：
```bash
docker compose -f deployments/docker-compose/auth-service.yml down
docker compose -f deployments/docker-compose/auth-service.yml up -d --force-recreate
```

### 問題 3: 連接被拒絕
**錯誤訊息：** `Failed to fetch` 或 `Network Error`

**檢查：**
```bash
# 確認後端服務運行正常
curl http://localhost:8080/health

# 確認可以登入
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### 問題 4: 前端環境變數未設定
**檢查檔案：** `web/.env.local`
```bash
cat web/.env.local | grep AUTH_SERVICE_URL
```

**應該看到：**
```
NEXT_PUBLIC_AUTH_SERVICE_URL=http://localhost:8080
```

如果沒有，創建或更新該檔案。

## 自動化測試 (使用 curl)

### 測試 CORS 預檢
```bash
curl -v -X OPTIONS http://localhost:8080/api/v1/auth/login \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: content-type"
```

**預期：** `204 No Content` + CORS headers

### 測試登入 API
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "Origin: http://localhost:3000" \
  -d '{"email":"test@example.com","password":"password123"}' \
  | jq
```

**預期：** 返回包含 `access_token`, `refresh_token`, 和 `user` 的 JSON

## 成功登入後的行為

1. **頁面跳轉：** 應該從 `/auth/login` 跳轉到 `/`
2. **Token 儲存：** localStorage 中應該有兩個 token
3. **認證狀態：** AuthContext 應該更新 `user` 和 `isAuthenticated` 狀態
4. **後續請求：** 之後的 API 請求應該自動帶上 `Authorization: Bearer <token>` header

## 進階：使用 Playwright 自動化測試

如果想要自動化測試，可以創建 Playwright 測試腳本：

```typescript
// tests/auth.spec.ts
import { test, expect } from '@playwright/test';

test('user can login successfully', async ({ page }) => {
  await page.goto('http://localhost:3000/auth/login');

  await page.fill('input[type="email"]', 'test@example.com');
  await page.fill('input[type="password"]', 'password123');

  await page.click('button[type="submit"]');

  // 等待跳轉
  await page.waitForURL('http://localhost:3000/');

  // 檢查 localStorage
  const accessToken = await page.evaluate(() => localStorage.getItem('access_token'));
  expect(accessToken).toBeTruthy();
});
```

---

如果遇到任何問題，請查看：
- 前端 Console: `F12` → Console 標籤
- 後端日誌: `docker logs tabelogo-auth-service -f`
- 網路請求: `F12` → Network 標籤
