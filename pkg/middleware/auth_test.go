package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const testJWTSecret = "test-secret-key-must-be-at-least-32-characters-long"

func setupTestMiddleware(t *testing.T) (*AuthMiddleware, *redis.Client) {
	// Create Redis client for testing
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15, // Use separate DB for testing
	})

	// Clear test DB
	ctx := context.Background()
	client.FlushDB(ctx)

	logger, _ := zap.NewDevelopment()
	middleware := NewAuthMiddleware(testJWTSecret, client, logger)

	return middleware, client
}

func createTestToken(userID, sessionID, role string, duration time.Duration) (string, error) {
	claims := &JWTClaims{
		UserID:    userID,
		SessionID: sessionID,
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(testJWTSecret))
}

func createTestSession(client *redis.Client, sessionID, userID string) error {
	ctx := context.Background()
	sessionKey := "session:" + sessionID

	return client.HSet(ctx, sessionKey, map[string]interface{}{
		"id":          sessionID,
		"user_id":     userID,
		"device_info": "test-device",
		"ip_address":  "127.0.0.1",
		"is_active":   "true",
	}).Err()
}

func TestAuthMiddleware_RequireAuth_Success(t *testing.T) {
	middleware, client := setupTestMiddleware(t)
	defer client.Close()

	// Create test session
	userID := "test-user-123"
	sessionID := "test-session-456"
	role := "user"

	err := createTestSession(client, sessionID, userID)
	require.NoError(t, err)

	// Create valid token
	token, err := createTestToken(userID, sessionID, role, 15*time.Minute)
	require.NoError(t, err)

	// Setup test request
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	// Test middleware
	middleware.RequireAuth()(c)

	// Verify
	assert.Equal(t, http.StatusOK, w.Code)
	assert.False(t, c.IsAborted())

	// Check context values
	contextUserID, exists := c.Get(UserIDKey)
	assert.True(t, exists)
	assert.Equal(t, userID, contextUserID)

	contextSessionID, exists := c.Get(SessionIDKey)
	assert.True(t, exists)
	assert.Equal(t, sessionID, contextSessionID)

	contextRole, exists := c.Get(UserRoleKey)
	assert.True(t, exists)
	assert.Equal(t, role, contextRole)
}

func TestAuthMiddleware_RequireAuth_MissingToken(t *testing.T) {
	middleware, client := setupTestMiddleware(t)
	defer client.Close()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	middleware.RequireAuth()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.True(t, c.IsAborted())
}

func TestAuthMiddleware_RequireAuth_InvalidToken(t *testing.T) {
	middleware, client := setupTestMiddleware(t)
	defer client.Close()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid-token")

	middleware.RequireAuth()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.True(t, c.IsAborted())
}

func TestAuthMiddleware_RequireAuth_ExpiredToken(t *testing.T) {
	middleware, client := setupTestMiddleware(t)
	defer client.Close()

	// Create expired token
	token, err := createTestToken("user-123", "session-456", "user", -1*time.Hour)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	middleware.RequireAuth()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.True(t, c.IsAborted())
}

func TestAuthMiddleware_RequireAuth_SessionNotFound(t *testing.T) {
	middleware, client := setupTestMiddleware(t)
	defer client.Close()

	// Create token but no session
	token, err := createTestToken("user-123", "nonexistent-session", "user", 15*time.Minute)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	middleware.RequireAuth()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.True(t, c.IsAborted())
}

func TestAuthMiddleware_RequireAuth_InactiveSession(t *testing.T) {
	middleware, client := setupTestMiddleware(t)
	defer client.Close()

	userID := "user-123"
	sessionID := "session-456"

	// Create inactive session
	ctx := context.Background()
	sessionKey := "session:" + sessionID
	err := client.HSet(ctx, sessionKey, map[string]interface{}{
		"id":        sessionID,
		"user_id":   userID,
		"is_active": "false", // Inactive
	}).Err()
	require.NoError(t, err)

	token, err := createTestToken(userID, sessionID, "user", 15*time.Minute)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	middleware.RequireAuth()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.True(t, c.IsAborted())
}

func TestAuthMiddleware_Optional_WithValidToken(t *testing.T) {
	middleware, client := setupTestMiddleware(t)
	defer client.Close()

	userID := "user-123"
	sessionID := "session-456"
	role := "user"

	err := createTestSession(client, sessionID, userID)
	require.NoError(t, err)

	token, err := createTestToken(userID, sessionID, role, 15*time.Minute)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	middleware.Optional()(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.False(t, c.IsAborted())

	// Should have user context
	contextUserID, exists := c.Get(UserIDKey)
	assert.True(t, exists)
	assert.Equal(t, userID, contextUserID)
}

func TestAuthMiddleware_Optional_WithoutToken(t *testing.T) {
	middleware, client := setupTestMiddleware(t)
	defer client.Close()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	middleware.Optional()(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.False(t, c.IsAborted())

	// Should NOT have user context
	_, exists := c.Get(UserIDKey)
	assert.False(t, exists)
}

func TestAuthMiddleware_Optional_WithInvalidToken(t *testing.T) {
	middleware, client := setupTestMiddleware(t)
	defer client.Close()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid-token")

	middleware.Optional()(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.False(t, c.IsAborted())

	// Should NOT have user context
	_, exists := c.Get(UserIDKey)
	assert.False(t, exists)
}

func TestAuthMiddleware_RequireRole_Success(t *testing.T) {
	middleware, client := setupTestMiddleware(t)
	defer client.Close()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(UserRoleKey, "admin")

	middleware.RequireRole("admin", "moderator")(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.False(t, c.IsAborted())
}

func TestAuthMiddleware_RequireRole_InsufficientPermissions(t *testing.T) {
	middleware, client := setupTestMiddleware(t)
	defer client.Close()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(UserRoleKey, "user")

	middleware.RequireRole("admin")(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.True(t, c.IsAborted())
}

func TestAuthMiddleware_RequireRole_MissingRole(t *testing.T) {
	middleware, client := setupTestMiddleware(t)
	defer client.Close()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	middleware.RequireRole("admin")(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.True(t, c.IsAborted())
}

func TestGetUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Test when user ID exists
	c.Set(UserIDKey, "user-123")
	userID, exists := GetUserID(c)
	assert.True(t, exists)
	assert.Equal(t, "user-123", userID)

	// Test when user ID doesn't exist
	c2, _ := gin.CreateTestContext(httptest.NewRecorder())
	userID, exists = GetUserID(c2)
	assert.False(t, exists)
	assert.Empty(t, userID)
}

func TestGetSessionID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Set(SessionIDKey, "session-456")
	sessionID, exists := GetSessionID(c)
	assert.True(t, exists)
	assert.Equal(t, "session-456", sessionID)
}

func TestGetUserRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Set(UserRoleKey, "admin")
	role, exists := GetUserRole(c)
	assert.True(t, exists)
	assert.Equal(t, "admin", role)
}
