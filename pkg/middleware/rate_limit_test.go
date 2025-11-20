package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limit := 3
	window := time.Second

	tests := []struct {
		name           string
		requests       int
		delay          time.Duration
		expectedStatus []int
	}{
		{
			name:     "Within limit",
			requests: 2,
			delay:    0,
			expectedStatus: []int{
				http.StatusOK,
				http.StatusOK,
			},
		},
		{
			name:     "Exceed limit",
			requests: 4,
			delay:    0,
			expectedStatus: []int{
				http.StatusOK,
				http.StatusOK,
				http.StatusOK,
				http.StatusTooManyRequests,
			},
		},
		{
			name:     "Limit reset after window",
			requests: 4,
			delay:    window + 100*time.Millisecond,
			expectedStatus: []int{
				http.StatusOK,
				http.StatusOK,
				http.StatusOK,
				http.StatusOK, // After window, limit resets
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(InMemoryRateLimit(limit, window))
			router.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			for i := 0; i < tt.requests; i++ {
				if i == 3 && tt.delay > 0 {
					time.Sleep(tt.delay)
				}

				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				assert.Equal(t, tt.expectedStatus[i], w.Code, "Request %d failed", i+1)

				// Check rate limit headers
				assert.NotEmpty(t, w.Header().Get("X-RateLimit-Limit"))
				assert.NotEmpty(t, w.Header().Get("X-RateLimit-Remaining"))
				assert.NotEmpty(t, w.Header().Get("X-RateLimit-Reset"))
			}
		})
	}
}

func TestInMemoryRateLimitDifferentClients(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limit := 2
	window := time.Second

	router := gin.New()
	router.Use(InMemoryRateLimit(limit, window))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Client 1 makes requests
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// Client 2 should have separate limit
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.2:1234"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRateLimitWithUserContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limit := 2
	window := time.Second

	router := gin.New()

	// Set user ID in context before rate limiting
	router.Use(func(c *gin.Context) {
		c.Set(UserIDKey, "user123")
		c.Next()
	})

	router.Use(InMemoryRateLimit(limit, window))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Make requests
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if i < 2 {
			assert.Equal(t, http.StatusOK, w.Code)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, w.Code)
		}
	}
}

func TestGetClientIdentifier(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupContext   func(*gin.Context)
		expectedPrefix string
	}{
		{
			name: "With user ID",
			setupContext: func(c *gin.Context) {
				c.Set(UserIDKey, "user123")
			},
			expectedPrefix: "user:",
		},
		{
			name: "Without user ID",
			setupContext: func(c *gin.Context) {
				// No user ID set
			},
			expectedPrefix: "ip:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/test", func(c *gin.Context) {
				tt.setupContext(c)
				identifier := getClientIdentifier(c)
				assert.Contains(t, identifier, tt.expectedPrefix)
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}
