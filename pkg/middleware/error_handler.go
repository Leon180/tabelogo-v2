package middleware

import (
	"net/http"

	"github.com/Leon180/tabelogo-v2/pkg/errors"
	"github.com/gin-gonic/gin"
)

// ErrorHandler returns a middleware that handles errors uniformly
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) == 0 {
			return
		}

		// Get the last error
		err := c.Errors.Last().Err

		// Check if it's an AppError
		if appErr, ok := errors.AsAppError(err); ok {
			c.JSON(appErr.HTTPStatus, gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
				"details": appErr.Details,
			})
			return
		}

		// Default to internal server error
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.ErrCodeInternal,
			"message": "Internal server error",
		})
	}
}
