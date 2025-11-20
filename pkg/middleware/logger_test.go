package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		path           string
		method         string
		expectedStatus int
		setupHandler   func(*gin.Context)
	}{
		{
			name:           "Successful request",
			path:           "/test",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			setupHandler: func(c *gin.Context) {
				c.Status(http.StatusOK)
			},
		},
		{
			name:           "Client error",
			path:           "/test",
			method:         http.MethodGet,
			expectedStatus: http.StatusBadRequest,
			setupHandler: func(c *gin.Context) {
				c.Status(http.StatusBadRequest)
			},
		},
		{
			name:           "Server error",
			path:           "/test",
			method:         http.MethodGet,
			expectedStatus: http.StatusInternalServerError,
			setupHandler: func(c *gin.Context) {
				c.Status(http.StatusInternalServerError)
			},
		},
		{
			name:           "Request with query parameters",
			path:           "/test?param1=value1&param2=value2",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			setupHandler: func(c *gin.Context) {
				c.Status(http.StatusOK)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test logger
			logger := zaptest.NewLogger(t)

			router := gin.New()
			router.Use(Logger(logger))
			router.Handle(tt.method, "/test", tt.setupHandler)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestLoggerWithErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := zaptest.NewLogger(t)

	router := gin.New()
	router.Use(Logger(logger))
	router.GET("/test", func(c *gin.Context) {
		c.Error(assert.AnError)
		c.Status(http.StatusInternalServerError)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestLoggerWithUserAgent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := zaptest.NewLogger(t)

	router := gin.New()
	router.Use(Logger(logger))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("User-Agent", "TestAgent/1.0")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
