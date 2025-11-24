# Auth Service API Specification

## Base URL
- **Development**: `http://localhost:8081`
- **Production**: TBD

## Authentication
Most endpoints require a Bearer token in the Authorization header:
```
Authorization: Bearer <access_token>
```

---

## Endpoints

### 1. Register User
Create a new user account.

**Endpoint**: `POST /auth/register`

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123",
  "username": "johndoe"
}
```

**Validation**:
- `email`: Required, valid email format
- `password`: Required, minimum 8 characters
- `username`: Required, minimum 3 characters

**Success Response** (201 Created):
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "johndoe",
    "role": "user",
    "email_verified": false,
    "created_at": "2025-11-24T10:00:00Z"
  }
}
```

**Error Responses**:
- `400 Bad Request`: Invalid request body
  ```json
  {
    "error": "invalid_request",
    "message": "validation error details"
  }
  ```
- `409 Conflict`: Email already exists
  ```json
  {
    "error": "email_exists",
    "message": "Email already registered"
  }
  ```

---

### 2. Login
Authenticate user and receive tokens.

**Endpoint**: `POST /auth/login`

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Success Response** (200 OK):
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "johndoe",
    "role": "user",
    "email_verified": false,
    "created_at": "2025-11-24T10:00:00Z"
  }
}
```

**Error Responses**:
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Invalid credentials
  ```json
  {
    "error": "invalid_credentials",
    "message": "Invalid email or password"
  }
  ```

---

### 3. Refresh Token
Get a new access token using refresh token.

**Endpoint**: `POST /auth/refresh`

**Request Body**:
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Success Response** (200 OK):
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Error Responses**:
- `400 Bad Request`: Missing refresh token
- `401 Unauthorized`: Invalid or expired refresh token
  ```json
  {
    "error": "invalid_token",
    "message": "Invalid or expired refresh token"
  }
  ```

---

### 4. Validate Token
Validate an access token and get user info.

**Endpoint**: `GET /auth/validate`

**Headers**:
```
Authorization: Bearer <access_token>
```

**Success Response** (200 OK):
```json
{
  "valid": true,
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "johndoe",
    "role": "user",
    "email_verified": false,
    "created_at": "2025-11-24T10:00:00Z"
  }
}
```

**Error Responses**:
- `401 Unauthorized`: Missing or invalid token
  ```json
  {
    "valid": false
  }
  ```

---

## Data Models

### User
```typescript
{
  id: string;              // UUID
  email: string;           // Valid email
  username: string;        // 3+ characters
  role: string;            // "admin" | "user" | "guest"
  email_verified: boolean;
  created_at: string;      // ISO 8601 timestamp
}
```

### Tokens
- **Access Token**: Short-lived JWT (typically 15-60 minutes)
- **Refresh Token**: Long-lived JWT (typically 7-30 days)

---

## Error Codes

| Code | Description |
|------|-------------|
| `invalid_request` | Request validation failed |
| `email_exists` | Email already registered |
| `invalid_credentials` | Wrong email or password |
| `invalid_token` | Token is invalid or expired |
| `missing_token` | Authorization header missing |
| `internal_error` | Server error |

---

## Frontend Integration Notes

1. **Token Storage**: Store tokens in `localStorage`
   - `access_token`: Used for API requests
   - `refresh_token`: Used to get new access token

2. **Token Refresh**: Implement automatic token refresh when access token expires

3. **Protected Routes**: Check token validity before accessing protected pages

4. **Logout**: Clear tokens from localStorage
