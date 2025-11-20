package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := zaptest.NewLogger(t)

	router := gin.New()
	router.Use(Recovery(logger))

	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	router.GET("/normal", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	t.Run("Panic recovery", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/panic", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "INTERNAL_ERROR")
		assert.Contains(t, w.Body.String(), "Internal server error")
	})

	t.Run("Normal request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/normal", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestRecoveryWithDifferentPanicTypes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := zaptest.NewLogger(t)

	tests := []struct {
		name      string
		panicVal  interface{}
	}{
		{
			name:     "String panic",
			panicVal: "string error",
		},
		{
			name:     "Integer panic",
			panicVal: 42,
		},
		{
			name:     "Nil panic",
			panicVal: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(Recovery(logger))
			router.GET("/test", func(c *gin.Context) {
				panic(tt.panicVal)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})
	}
}
