package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// RequestIDHeader is the header key for request ID
	RequestIDHeader = "X-Request-ID"
	// RequestIDKey is the context key for request ID
	RequestIDKey = "request_id"
)

// RequestID returns a middleware that adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID is already present in header
		requestID := c.GetHeader(RequestIDHeader)

		// Generate new request ID if not present
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set request ID in context
		c.Set(RequestIDKey, requestID)

		// Set request ID in response header
		c.Writer.Header().Set(RequestIDHeader, requestID)

		c.Next()
	}
}

// GetRequestID retrieves request ID from context
func GetRequestID(c *gin.Context) (string, bool) {
	requestID, exists := c.Get(RequestIDKey)
	if !exists {
		return "", false
	}
	id, ok := requestID.(string)
	return id, ok
}
