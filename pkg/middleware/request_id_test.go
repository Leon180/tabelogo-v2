package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		existingRequestID  string
		expectGenerated    bool
	}{
		{
			name:              "Generate new request ID",
			existingRequestID: "",
			expectGenerated:   true,
		},
		{
			name:              "Use existing request ID",
			existingRequestID: "existing-id-123",
			expectGenerated:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(RequestID())

			var requestID string
			router.GET("/test", func(c *gin.Context) {
				if id, exists := c.Get(RequestIDKey); exists {
					requestID = id.(string)
				}
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.existingRequestID != "" {
				req.Header.Set(RequestIDHeader, tt.existingRequestID)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.NotEmpty(t, requestID)

			// Check response header
			responseID := w.Header().Get(RequestIDHeader)
			assert.NotEmpty(t, responseID)
			assert.Equal(t, requestID, responseID)

			if tt.expectGenerated {
				// Should be a UUID format
				assert.Len(t, requestID, 36) // UUID length with hyphens
			} else {
				// Should match existing ID
				assert.Equal(t, tt.existingRequestID, requestID)
			}
		})
	}
}

func TestGetRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		c.Set(RequestIDKey, "test-request-id")
		requestID, exists := GetRequestID(c)
		assert.True(t, exists)
		assert.Equal(t, "test-request-id", requestID)
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequestIDNotExists(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		requestID, exists := GetRequestID(c)
		assert.False(t, exists)
		assert.Empty(t, requestID)
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
