package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Leon180/tabelogo-v2/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupError     func(*gin.Context)
		expectedStatus int
		expectedCode   errors.ErrorCode
	}{
		{
			name: "No error",
			setupError: func(c *gin.Context) {
				c.Status(http.StatusOK)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "AppError - Not Found",
			setupError: func(c *gin.Context) {
				appErr := errors.New(errors.ErrCodeNotFound, "Resource not found")
				c.Error(appErr)
			},
			expectedStatus: http.StatusNotFound,
			expectedCode:   errors.ErrCodeNotFound,
		},
		{
			name: "AppError - Validation Error",
			setupError: func(c *gin.Context) {
				appErr := errors.New(errors.ErrCodeValidationFailed, "Validation failed")
				c.Error(appErr)
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   errors.ErrCodeValidationFailed,
		},
		{
			name: "AppError with details",
			setupError: func(c *gin.Context) {
				appErr := errors.New(errors.ErrCodeValidationFailed, "Validation failed").
					WithDetails(map[string]interface{}{
						"field": "email",
						"issue": "invalid format",
					})
				c.Error(appErr)
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   errors.ErrCodeValidationFailed,
		},
		{
			name: "Generic error",
			setupError: func(c *gin.Context) {
				c.Error(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   errors.ErrCodeInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(ErrorHandler())
			router.GET("/test", tt.setupError)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedCode != "" {
				assert.Contains(t, w.Body.String(), string(tt.expectedCode))
			}
		})
	}
}

func TestErrorHandlerMultipleErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		// Add multiple errors
		c.Error(errors.New(errors.ErrCodeValidationFailed, "First error"))
		c.Error(errors.New(errors.ErrCodeNotFound, "Second error"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should handle the last error
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), string(errors.ErrCodeNotFound))
}
