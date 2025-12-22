package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/Leon180/tabelogo-v2/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	// AuthorizationHeader is the header key for authorization
	AuthorizationHeader = "Authorization"
	// BearerPrefix is the prefix for Bearer token
	BearerPrefix = "Bearer "
	// UserIDKey is the context key for user ID
	UserIDKey = "user_id"
	// SessionIDKey is the context key for session ID
	SessionIDKey = "session_id"
	// UserRoleKey is the context key for user role
	UserRoleKey = "user_role"
)

// JWTClaims represents the JWT claims with session support
type JWTClaims struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"` // NEW: Session ID for validation
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware handles authentication with session validation
type AuthMiddleware struct {
	jwtSecret string
	redis     *redis.Client
	logger    *zap.Logger
}

// NewAuthMiddleware creates a new auth middleware with session support
func NewAuthMiddleware(jwtSecret string, redis *redis.Client, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
		redis:     redis,
		logger:    logger,
	}
}

// RequireAuth validates JWT and checks session is active
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extract token from header
		token, err := m.extractToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    errors.ErrCodeUnauthorized,
				"message": "Missing or invalid authorization header",
			})
			c.Abort()
			return
		}

		// 2. Verify JWT signature and expiry
		claims, err := m.verifyJWT(token)
		if err != nil {
			m.logger.Debug("JWT verification failed", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    errors.ErrCodeUnauthorized,
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// 3. Validate session in Redis
		if err := m.validateSession(c.Request.Context(), claims.SessionID); err != nil {
			m.logger.Debug("Session validation failed",
				zap.String("session_id", claims.SessionID),
				zap.Error(err),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    errors.ErrCodeUnauthorized,
				"message": "Session expired or revoked",
			})
			c.Abort()
			return
		}

		// 4. Set user context
		c.Set(UserIDKey, claims.UserID)
		c.Set(SessionIDKey, claims.SessionID)
		c.Set(UserRoleKey, claims.Role)

		c.Next()
	}
}

// Optional validates token if present but doesn't require it
func (m *AuthMiddleware) Optional() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to extract token
		token, err := m.extractToken(c)
		if err != nil {
			// No token or invalid format - continue without auth
			c.Next()
			return
		}

		// Try to verify JWT
		claims, err := m.verifyJWT(token)
		if err != nil {
			// Invalid token - continue without auth
			c.Next()
			return
		}

		// Try to validate session
		if err := m.validateSession(c.Request.Context(), claims.SessionID); err != nil {
			// Invalid session - continue without auth
			c.Next()
			return
		}

		// Valid auth - set context
		c.Set(UserIDKey, claims.UserID)
		c.Set(SessionIDKey, claims.SessionID)
		c.Set(UserRoleKey, claims.Role)

		c.Next()
	}
}

// RequireRole checks if user has required role
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
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
		if !slices.Contains(roles, roleStr) {
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

// extractToken extracts JWT token from Authorization header
func (m *AuthMiddleware) extractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader(AuthorizationHeader)
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	if !strings.HasPrefix(authHeader, BearerPrefix) {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return strings.TrimPrefix(authHeader, BearerPrefix), nil
}

// verifyJWT verifies JWT signature and returns claims
func (m *AuthMiddleware) verifyJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(m.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// validateSession checks if session is active in Redis
func (m *AuthMiddleware) validateSession(ctx context.Context, sessionID string) error {
	sessionKey := "session:" + sessionID

	// Get session as JSON string
	sessionJSON, err := m.redis.Get(ctx, sessionKey).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("session not found")
		}
		return fmt.Errorf("redis error: %w", err)
	}

	// Parse JSON to check is_active field
	var session struct {
		IsActive bool `json:"is_active"`
	}

	if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
		return fmt.Errorf("failed to parse session data: %w", err)
	}

	if !session.IsActive {
		return fmt.Errorf("session is not active")
	}

	return nil
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

// GetSessionID retrieves session ID from context
func GetSessionID(c *gin.Context) (string, bool) {
	sessionID, exists := c.Get(SessionIDKey)
	if !exists {
		return "", false
	}
	id, ok := sessionID.(string)
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
