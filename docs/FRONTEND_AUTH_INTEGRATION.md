# Frontend Integration Guide - Authentication

## Overview

This guide explains how the frontend integrates with the authenticated backend services after Phase 2 authentication integration.

---

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                  Frontend (Next.js)                  │
│                                                      │
│  ┌────────────────────────────────────────────┐   │
│  │         AuthContext (React Context)         │   │
│  │  - Manages user state                       │   │
│  │  - Handles login/logout                     │   │
│  │  - Stores JWT tokens in localStorage        │   │
│  └────────────────────────────────────────────┘   │
│                        │                            │
│                        ▼                            │
│  ┌────────────────────────────────────────────┐   │
│  │         API Clients (Axios)                 │   │
│  │  - auth-service.ts                          │   │
│  │  - spider-service.ts                        │   │
│  │  - restaurant-service.ts (future)           │   │
│  │                                              │   │
│  │  Each client:                                │   │
│  │  1. Adds Authorization header               │   │
│  │  2. Handles 401 errors                      │   │
│  │  3. Auto-refreshes tokens                   │   │
│  └────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────┐
│              Backend Services (Go)                   │
│                                                      │
│  Auth Service → Spider Service → Restaurant Service │
│       ↓              ↓                  ↓            │
│  Auth Middleware validates JWT + Session            │
└─────────────────────────────────────────────────────┘
```

---

## Authentication Flow

### 1. User Login

```typescript
// User logs in via AuthContext
const { login } = useAuth();
await login({ email, password });

// AuthContext calls auth-service.ts
const response = await authService.login(data);

// Tokens stored in localStorage
localStorage.setItem('access_token', response.access_token);
localStorage.setItem('refresh_token', response.refresh_token);

// User state updated
setUser(response.user);
```

### 2. Making Authenticated Requests

```typescript
// Example: Spider Service request
import { searchTabelog } from '@/lib/api/spider-service';

// Request interceptor automatically adds Authorization header
const results = await searchTabelog({
  google_id: 'ChIJ...',
  area: 'Tokyo',
  place_name: 'Restaurant Name'
});

// Behind the scenes:
// 1. Interceptor reads token from localStorage
// 2. Adds header: Authorization: Bearer <token>
// 3. Backend validates token + session
// 4. Returns data or 401
```

### 3. Token Refresh (Automatic)

```typescript
// If access token expires (401 response):
// 1. auth-service.ts interceptor catches 401
// 2. Automatically calls /api/v1/auth/refresh
// 3. Gets new access_token and refresh_token
// 4. Updates localStorage
// 5. Retries original request with new token
// 6. User doesn't notice anything!
```

---

## API Client Configuration

### Auth Service Client

**File**: `web/src/lib/api/auth-service.ts`

**Features**:
- ✅ Request interceptor (adds Authorization header)
- ✅ Response interceptor (handles 401, auto-refresh)
- ✅ Token storage (localStorage)
- ✅ Refresh token rotation

**Usage**:
```typescript
import { authService } from '@/lib/api/auth-service';

// Login
const response = await authService.login({ email, password });

// Logout
await authService.logout();

// Validate token
const user = await authService.validateToken();
```

### Spider Service Client

**File**: `web/src/lib/api/spider-service.ts`

**Features**:
- ✅ Request interceptor (adds Authorization header) **[JUST ADDED]**
- ✅ Response interceptor (handles 401) **[JUST ADDED]**
- ✅ Error handling for auth failures

**Usage**:
```typescript
import { searchTabelog, getJobStatus } from '@/lib/api/spider-service';

// Start scraping job (requires auth)
const job = await searchTabelog({
  google_id: 'ChIJ...',
  area: 'Tokyo',
  place_name: 'Restaurant'
});

// Get job status (requires auth)
const status = await getJobStatus(job.job_id);
```

### Restaurant Service Client

**File**: `web/src/lib/api/restaurant-service.ts` (if exists)

**TODO**: Add similar auth interceptors if not already present.

---

## Token Storage

### LocalStorage Keys

| Key | Description | Example |
|-----|-------------|---------|
| `access_token` | JWT for API requests | `eyJhbGci...` |
| `refresh_token` | Token for refreshing access | `eyJhbGci...` |

### Security Considerations

**Current**: Tokens stored in localStorage
- ✅ Pros: Simple, works across tabs
- ⚠️ Cons: Vulnerable to XSS attacks

**Future Improvements**:
- Use httpOnly cookies (more secure)
- Implement Content Security Policy (CSP)
- Add XSS protection headers

---

## Error Handling

### 401 Unauthorized

**Scenario**: Token invalid or expired

**Handling**:
```typescript
// Automatic refresh attempt
try {
  const response = await spiderService.searchTabelog(...);
} catch (error) {
  if (error.message.includes('Authentication required')) {
    // Token refresh failed, redirect to login
    router.push('/auth/login');
  }
}
```

### 403 Forbidden

**Scenario**: User lacks required role/permissions

**Handling**:
```typescript
// Show error message
toast.error('You do not have permission to perform this action');
```

---

## Component Integration

### Using AuthContext

```typescript
'use client';

import { useAuth } from '@/contexts/AuthContext';

export function MyComponent() {
  const { user, isAuthenticated, isLoading, logout } = useAuth();

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (!isAuthenticated) {
    return <div>Please login</div>;
  }

  return (
    <div>
      <p>Welcome, {user.username}!</p>
      <button onClick={logout}>Logout</button>
    </div>
  );
}
```

### Protected Routes

```typescript
// middleware.ts or page component
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

export function ProtectedPage() {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/auth/login');
    }
  }, [isAuthenticated, isLoading, router]);

  if (isLoading) return <div>Loading...</div>;
  if (!isAuthenticated) return null;

  return <div>Protected Content</div>;
}
```

---

## Testing Frontend Integration

### Manual Testing

1. **Start Backend Services**:
   ```bash
   cd /path/to/backend
   make up
   ```

2. **Start Frontend**:
   ```bash
   cd /path/to/frontend
   npm run dev
   ```

3. **Test Flow**:
   - Navigate to `/auth/login`
   - Login with test credentials
   - Check localStorage for tokens
   - Make authenticated request (e.g., search restaurants)
   - Verify request includes Authorization header (DevTools → Network)
   - Logout and verify tokens are cleared

### Browser DevTools

**Check Authorization Header**:
1. Open DevTools (F12)
2. Go to Network tab
3. Make an API request
4. Click on the request
5. Check Headers → Request Headers
6. Should see: `Authorization: Bearer eyJhbGci...`

**Check LocalStorage**:
1. Open DevTools (F12)
2. Go to Application tab
3. Expand Local Storage
4. Check for `access_token` and `refresh_token`

---

## Common Issues & Solutions

### Issue 1: "Authentication required" Error

**Symptom**: All API requests return 401

**Causes**:
1. Not logged in
2. Token expired and refresh failed
3. Backend not running

**Solution**:
```typescript
// Check if token exists
const token = localStorage.getItem('access_token');
console.log('Token:', token ? 'exists' : 'missing');

// Try logging in again
await authService.login({ email, password });
```

### Issue 2: CORS Errors

**Symptom**: "CORS policy" error in console

**Solution**: Backend already configured for CORS
- Check backend is running
- Verify API URLs in frontend match backend ports

### Issue 3: Token Not Sent

**Symptom**: Request doesn't include Authorization header

**Causes**:
1. Interceptor not configured
2. Token not in localStorage

**Solution**:
```typescript
// Verify interceptor is set up
spiderClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  console.log('Adding token:', token ? 'yes' : 'no');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});
```

---

## Next Steps

### Immediate
1. ✅ Update Spider Service client (DONE)
2. ⏳ Test frontend login flow
3. ⏳ Test authenticated Spider requests
4. ⏳ Verify token refresh works

### Future Enhancements
1. Add Restaurant Service client auth
2. Add Map Service client auth
3. Implement httpOnly cookies
4. Add request retry logic
5. Add loading states for auth operations
6. Add toast notifications for auth errors

---

## API Endpoints Reference

### Auth Service (Port 8080)

| Endpoint | Method | Auth Required | Description |
|----------|--------|---------------|-------------|
| `/api/v1/auth/register` | POST | No | Register new user |
| `/api/v1/auth/login` | POST | No | Login user |
| `/api/v1/auth/refresh` | POST | No | Refresh access token |
| `/api/v1/auth/validate` | GET | Yes | Validate current token |

### Spider Service (Port 18084)

| Endpoint | Method | Auth Required | Description |
|----------|--------|---------------|-------------|
| `/api/v1/spider/scrape` | POST | **Yes** | Start scraping job |
| `/api/v1/spider/jobs/:id` | GET | **Yes** | Get job status |
| `/api/v1/spider/jobs/:id/stream` | GET | **Yes** | Stream job updates (SSE) |

### Restaurant Service (Port 18082)

| Endpoint | Method | Auth Required | Description |
|----------|--------|---------------|-------------|
| `/api/v1/restaurants/search` | GET | No (Optional) | Search restaurants |
| `/api/v1/restaurants/:id` | GET | No (Optional) | Get restaurant details |
| `/api/v1/favorites` | POST | **Yes** | Add to favorites |
| `/api/v1/users/:id/favorites` | GET | **Yes** | Get user favorites |

---

## Summary

✅ **What's Working**:
- Auth Service client with auto-refresh
- Spider Service client with auth headers **[JUST ADDED]**
- AuthContext for state management
- Token storage in localStorage

⏳ **What to Test**:
- Frontend login → Spider Service request flow
- Token refresh on expiration
- Error handling for auth failures

🚀 **Ready for Frontend Testing!**
