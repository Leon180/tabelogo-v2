package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testSecret := "test-secret"
	config := AuthConfig{
		JWTSecret: testSecret,
		SkipPaths: []string{"/health"},
	}

	tests := []struct {
		name           string
		path           string
		setupToken     func() string
		expectedStatus int
		expectedInContext bool
	}{
		{
			name:           "Valid token",
			path:           "/api/test",
			setupToken: func() string {
				return createTestToken(testSecret, "user123", "user", time.Hour)
			},
			expectedStatus:    http.StatusOK,
			expectedInContext: true,
		},
		{
			name: "Missing token",
			path: "/api/test",
			setupToken: func() string {
				return ""
			},
			expectedStatus:    http.StatusUnauthorized,
			expectedInContext: false,
		},
		{
			name: "Invalid token format",
			path: "/api/test",
			setupToken: func() string {
				return "InvalidToken"
			},
			expectedStatus:    http.StatusUnauthorized,
			expectedInContext: false,
		},
		{
			name: "Expired token",
			path: "/api/test",
			setupToken: func() string {
				return createTestToken(testSecret, "user123", "user", -time.Hour)
			},
			expectedStatus:    http.StatusUnauthorized,
			expectedInContext: false,
		},
		{
			name: "Skip path",
			path: "/health",
			setupToken: func() string {
				return ""
			},
			expectedStatus:    http.StatusOK,
			expectedInContext: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(Auth(config))

			var userID string
			var userRole string
			router.GET(tt.path, func(c *gin.Context) {
				if id, exists := c.Get(UserIDKey); exists {
					userID = id.(string)
				}
				if role, exists := c.Get(UserRoleKey); exists {
					userRole = role.(string)
				}
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			token := tt.setupToken()
			if token != "" {
				req.Header.Set(AuthorizationHeader, BearerPrefix+token)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedInContext {
				assert.NotEmpty(t, userID)
				assert.NotEmpty(t, userRole)
			}
		})
	}
}

func TestRequireRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userRole       string
		requiredRoles  []string
		expectedStatus int
	}{
		{
			name:           "Allowed role",
			userRole:       "admin",
			requiredRoles:  []string{"admin", "moderator"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Forbidden role",
			userRole:       "user",
			requiredRoles:  []string{"admin", "moderator"},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Missing role",
			userRole:       "",
			requiredRoles:  []string{"admin"},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()

			// Setup context with user role
			router.Use(func(c *gin.Context) {
				if tt.userRole != "" {
					c.Set(UserRoleKey, tt.userRole)
				}
				c.Next()
			})

			router.Use(RequireRole(tt.requiredRoles...))
			router.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestGetUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		c.Set(UserIDKey, "user123")
		userID, exists := GetUserID(c)
		assert.True(t, exists)
		assert.Equal(t, "user123", userID)
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUserRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		c.Set(UserRoleKey, "admin")
		userRole, exists := GetUserRole(c)
		assert.True(t, exists)
		assert.Equal(t, "admin", userRole)
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// createTestToken creates a JWT token for testing
func createTestToken(secret, userID, role string, expiry time.Duration) string {
	claims := JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}
