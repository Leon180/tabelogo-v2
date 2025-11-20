package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/Leon180/tabelogo-v2/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	// AuthorizationHeader is the header key for authorization
	AuthorizationHeader = "Authorization"
	// BearerPrefix is the prefix for Bearer token
	BearerPrefix = "Bearer "
	// UserIDKey is the context key for user ID
	UserIDKey = "user_id"
	// UserRoleKey is the context key for user role
	UserRoleKey = "user_role"
)

// JWTClaims represents the JWT claims
type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthConfig holds the configuration for auth middleware
type AuthConfig struct {
	// JWTSecret is the secret key for JWT validation
	JWTSecret string
	// SkipPaths are paths that skip authentication
	SkipPaths []string
}

// Auth returns a middleware that validates JWT tokens
func Auth(config AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if path should skip authentication
		if slices.Contains(config.SkipPaths, c.Request.URL.Path) {
			c.Next()
			return
		}

		// Get token from header
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    errors.ErrCodeUnauthorized,
				"message": "Missing authorization header",
			})
			c.Abort()
			return
		}

		// Check Bearer prefix
		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    errors.ErrCodeUnauthorized,
				"message": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New(errors.ErrCodeUnauthorized, "Invalid signing method")
			}
			return []byte(config.JWTSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    errors.ErrCodeUnauthorized,
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(*JWTClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    errors.ErrCodeUnauthorized,
				"message": "Invalid token claims",
			})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set(UserIDKey, claims.UserID)
		c.Set(UserRoleKey, claims.Role)

		c.Next()
	}
}

// RequireRole returns a middleware that checks if user has required role
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context
		userRole, exists := c.Get(UserRoleKey)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    errors.ErrCodeForbidden,
				"message": "User role not found",
			})
			c.Abort()
			return
		}

		// Check if user has required role
		roleStr, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    errors.ErrCodeForbidden,
				"message": "Invalid user role",
			})
			c.Abort()
			return
		}

		// Check if role is in allowed roles
		allowed := slices.Contains(roles, roleStr)

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    errors.ErrCodeForbidden,
				"message": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserID retrieves user ID from context
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return "", false
	}
	id, ok := userID.(string)
	return id, ok
}

// GetUserRole retrieves user role from context
func GetUserRole(c *gin.Context) (string, bool) {
	userRole, exists := c.Get(UserRoleKey)
	if !exists {
		return "", false
	}
	role, ok := userRole.(string)
	return role, ok
}
