# Authentication Implementation - Best Practices

## 概述

本文檔說明前端認證系統的實現，包括 token 管理、自動刷新、安全考量和最佳實踐。

## 架構

### 組件結構

```
web/
├── src/
│   ├── contexts/
│   │   └── AuthContext.tsx          # 全局認證狀態管理
│   ├── lib/
│   │   └── api/
│   │       └── auth-service.ts       # API 客戶端和 interceptors
│   ├── types/
│   │   └── user.ts                   # TypeScript 類型定義
│   └── app/
│       └── auth/
│           └── login/
│               └── page.tsx          # 登入頁面
```

## Token 管理

### 當前實現：localStorage

**位置：** `web/src/lib/api/auth-service.ts`

```typescript
// 儲存 tokens
localStorage.setItem('access_token', response.data.access_token);
localStorage.setItem('refresh_token', response.data.refresh_token);

// 讀取 tokens
const token = localStorage.getItem('access_token');
const refreshToken = localStorage.getItem('refresh_token');
```

### 安全考量

#### ✅ 優點
1. **簡單實現**：容易理解和維護
2. **跨標籤頁共享**：同一域名下的所有標籤頁共享認證狀態
3. **持久化**：關閉瀏覽器後 token 仍然保留
4. **框架無關**：不依賴特定框架或後端配置

#### ⚠️ 缺點和風險
1. **XSS 攻擊風險**：
   - 如果應用存在 XSS 漏洞，攻擊者可以通過 JavaScript 讀取 localStorage
   - 攻擊者可以竊取 access token 和 refresh token

2. **無法設置 HttpOnly**：
   - localStorage 總是可以被 JavaScript 訪問
   - 無法使用 HttpOnly flag 保護

### 🔐 安全最佳實踐

#### 1. **當前實現的保護措施**

**a) Token 自動刷新機制**
```typescript
// 當 access token 過期 (401)，自動使用 refresh token 獲取新 token
authClient.interceptors.response.use(
    (response) => response,
    async (error: AxiosError) => {
        if (error.response?.status === 401) {
            // 自動刷新 token
            const newToken = await refreshToken();
            // 重試原請求
        }
    }
);
```

**b) 短期 Access Token**
- Access Token 有效期：15 分鐘
- Refresh Token 有效期：24 小時
- 減少 token 被濫用的時間窗口

**c) 防止 XSS**
- React 自動轉義輸出，防止 XSS
- 使用 Content Security Policy (CSP)
- 避免使用 `dangerouslySetInnerHTML`
- 驗證所有用戶輸入

#### 2. **推薦的進階實現：HttpOnly Cookies** (未來改進)

**優點：**
- ✅ JavaScript 無法訪問 (HttpOnly)
- ✅ 自動跨域保護 (SameSite)
- ✅ 防止 XSS 攻擊

**實現方式：**

後端需要修改：
```go
// 返回 HttpOnly cookie 而非 JSON
http.SetCookie(w, &http.Cookie{
    Name:     "access_token",
    Value:    token,
    HttpOnly: true,
    Secure:   true,
    SameSite: http.SameSiteStrictMode,
    Path:     "/",
    MaxAge:   900, // 15 minutes
})
```

前端需要修改：
```typescript
// 不再手動管理 token，瀏覽器自動發送 cookie
const authClient = axios.create({
    baseURL: AUTH_SERVICE_URL,
    withCredentials: true, // 允許發送 cookies
});
```

**缺點：**
- ❌ 需要後端大幅修改
- ❌ CORS 配置更複雜
- ❌ 無法輕鬆檢查 token (除非提供專門的 API)
- ❌ 移動應用支援較複雜

## 自動 Token 刷新

### 實現原理

當 API 請求返回 401 (Unauthorized) 時，自動執行以下流程：

```
1. 檢測到 401 錯誤
   ↓
2. 檢查是否已在刷新中
   ↓ 是：將請求加入隊列
   ↓ 否：開始刷新流程
3. 使用 refresh_token 調用 /auth/refresh
   ↓
4. 獲取新的 access_token 和 refresh_token
   ↓
5. 更新 localStorage
   ↓
6. 使用新 token 重試原請求
   ↓
7. 處理隊列中的其他請求
```

### 防止重複刷新

```typescript
let isRefreshing = false;
let failedQueue: Array<{
    resolve: (value?: unknown) => void;
    reject: (reason?: unknown) => void;
}> = [];

// 如果已經在刷新，將請求加入隊列
if (isRefreshing) {
    return new Promise((resolve, reject) => {
        failedQueue.push({ resolve, reject });
    });
}
```

### 處理刷新失敗

如果 refresh token 也過期或無效：
1. 清除所有 tokens
2. 拒絕所有隊列中的請求
3. 用戶需要重新登入

```typescript
catch (refreshError) {
    processQueue(refreshError as Error, null);
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    return Promise.reject(refreshError);
}
```

## 用戶狀態管理

### AuthContext 提供的功能

```typescript
interface AuthContextType {
    user: User | null;           // 當前用戶資訊
    isLoading: boolean;          // 初始化載入狀態
    login: (data: LoginRequest) => Promise<void>;
    register: (data: RegisterRequest) => Promise<void>;
    logout: () => Promise<void>;
}
```

### 初始化流程

應用啟動時自動檢查認證狀態：

```typescript
useEffect(() => {
    checkAuth();
}, []);

const checkAuth = async () => {
    const token = localStorage.getItem('access_token');
    if (token) {
        const userData = await authService.validateToken();
        setUser(userData);
    }
    setIsLoading(false);
};
```

### 登入流程優化

**改進前：**
```typescript
// ❌ 登入後再調用 validateToken (多一次 API 請求)
const response = await authService.login(data);
const userData = await authService.validateToken();
setUser(userData);
```

**改進後：**
```typescript
// ✅ 直接使用 login response 中的 user 資料
const response = await authService.login(data);
setUser(response.user);
```

## 用戶資訊顯示

### 首頁顯示用戶名

**位置：** `web/src/app/page.tsx`

```tsx
const { user, logout } = useAuth();

{user ? (
    <div className="flex items-center gap-4">
        <span className="text-zinc-300">Hi, {user.username}</span>
        <Button onClick={() => logout()}>Logout</Button>
    </div>
) : (
    <Link href="/auth/login">
        <Button>Login</Button>
    </Link>
)}
```

### 用戶資料結構

```typescript
interface User {
    id: string;
    email: string;
    username: string;
    role: 'admin' | 'user' | 'guest';
    email_verified: boolean;
    created_at: string;
}
```

## API 類型定義

### 與後端完全匹配

**LoginResponse (後端):**
```go
type LoginResponse struct {
    AccessToken  string        `json:"access_token"`
    RefreshToken string        `json:"refresh_token"`
    User         *UserResponse `json:"user"`
}
```

**AuthResponse (前端):**
```typescript
interface AuthResponse {
    access_token: string;
    refresh_token: string;
    user: User;
}
```

## 錯誤處理

### 登入錯誤

```typescript
try {
    await login(data);
} catch (err: any) {
    setError(err.response?.data?.error || 'Failed to login');
}
```

### Token 過期處理

自動處理，用戶無感知：
- Access token 過期 → 自動刷新 → 重試請求
- Refresh token 過期 → 清除 tokens → 需要重新登入

### 網路錯誤

```typescript
authClient.interceptors.response.use(
    (response) => response,
    (error) => {
        if (!error.response) {
            // 網路錯誤 (無法連接到伺服器)
            console.error('Network error:', error.message);
        }
        return Promise.reject(error);
    }
);
```

## 測試指南

### 手動測試步驟

1. **正常登入**
   ```
   訪問 http://localhost:3000/auth/login
   輸入: test@example.com / password123
   預期: 跳轉到首頁，顯示 "Hi, TestUser"
   ```

2. **檢查 Token**
   ```
   F12 → Application → Local Storage
   預期: 看到 access_token 和 refresh_token
   ```

3. **Token 自動刷新**
   ```
   等待 15 分鐘 (access token 過期)
   執行需要認證的操作
   預期: 自動刷新，操作成功
   ```

4. **登出**
   ```
   點擊 Logout 按鈕
   預期: tokens 被清除，跳轉到登入頁
   ```

### 使用開發者工具測試

#### 模擬 Token 過期

```javascript
// 在 Console 中執行
localStorage.setItem('access_token', 'invalid_token');

// 然後執行需要認證的操作，應該會自動刷新
```

#### 檢查 API 請求

```
F12 → Network → 篩選 XHR
查看:
- OPTIONS /api/v1/auth/login (預檢)
- POST /api/v1/auth/login (登入)
- GET /api/v1/auth/validate (驗證)
```

## 最佳實踐總結

### ✅ 已實現

1. **類型安全**: TypeScript 完整類型定義
2. **自動刷新**: 401 時自動刷新 token
3. **防止重複**: 刷新時隊列化請求
4. **錯誤處理**: 完整的錯誤處理流程
5. **用戶體驗**: 登入後立即顯示用戶資訊
6. **短期 Token**: 15 分鐘 access token
7. **CORS 支援**: 正確的 CORS headers

### 🔄 可選改進

1. **HttpOnly Cookies**: 更安全但需要後端配置
2. **Token 加密**: 加密 localStorage 中的 tokens
3. **會話管理**: 實現 "記住我" 功能
4. **多裝置登出**: 實現全局登出功能
5. **Token 黑名單**: 後端維護失效 token 列表

### ⚠️ 安全注意事項

1. **永遠使用 HTTPS**: 生產環境必須使用 HTTPS
2. **CSP Headers**: 配置 Content Security Policy
3. **定期更新依賴**: 防止已知的安全漏洞
4. **輸入驗證**: 後端驗證所有輸入
5. **Rate Limiting**: 防止暴力破解
6. **審計日誌**: 記錄登入和敏感操作

## 故障排查

### 問題 1: Token 沒有自動刷新

**檢查：**
```typescript
// 確認 interceptor 已註冊
console.log('Interceptors:', authClient.interceptors.response);

// 確認 refresh token 存在
console.log('Refresh Token:', localStorage.getItem('refresh_token'));
```

### 問題 2: 用戶資訊沒有顯示

**檢查：**
```typescript
// 在 AuthContext 中添加 debug
console.log('User:', user);
console.log('Loading:', isLoading);

// 確認 API 返回正確的結構
console.log('Login Response:', response.data);
```

### 問題 3: CORS 錯誤

**檢查：**
```bash
# 確認後端 CORS 配置
curl -v -X OPTIONS http://localhost:8080/api/v1/auth/login \
  -H "Origin: http://localhost:3000"
```

## 相關文檔

- [TESTING_LOGIN.md](./TESTING_LOGIN.md) - 登入功能測試指南
- [FRONTEND_SUMMARY.md](../FRONTEND_SUMMARY.md) - 前端架構總覽
- [Backend Auth Service](../../internal/auth/README.md) - 後端認證服務

## 更新日誌

### 2025-11-25
- ✅ 修復 TypeScript 類型定義匹配後端
- ✅ 實現自動 token 刷新機制
- ✅ 優化登入流程，減少不必要的 API 調用
- ✅ 添加完整的用戶資訊顯示
- ✅ 文檔完善
