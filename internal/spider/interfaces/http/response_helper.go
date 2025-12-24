package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RespondOK sends a successful JSON response
func RespondOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// RespondCreated sends a 201 Created response
func RespondCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

// RespondAccepted sends a 202 Accepted response
func RespondAccepted(c *gin.Context, data interface{}) {
	c.JSON(http.StatusAccepted, data)
}

// RespondError sends an error response with the given status code
func RespondError(c *gin.Context, code int, err error) {
	c.JSON(code, ErrorResponse{Error: err.Error()})
}

// RespondBadRequest sends a 400 Bad Request response
func RespondBadRequest(c *gin.Context, err error) {
	RespondError(c, http.StatusBadRequest, err)
}

// RespondNotFound sends a 404 Not Found response
func RespondNotFound(c *gin.Context, err error) {
	RespondError(c, http.StatusNotFound, err)
}

// RespondInternalError sends a 500 Internal Server Error response
func RespondInternalError(c *gin.Context, err error) {
	RespondError(c, http.StatusInternalServerError, err)
}

// RespondValidationError sends a 400 Bad Request response for validation errors
func RespondValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error":   "validation failed",
		"details": err.Error(),
	})
}
