package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Leon180/tabelogo-v2/internal/map/application/usecases"
	"github.com/Leon180/tabelogo-v2/internal/map/domain/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// MapHandler handles HTTP requests for map service
type MapHandler struct {
	logger          *zap.Logger
	quickSearchUC   *usecases.QuickSearchUseCase
	advanceSearchUC *usecases.AdvanceSearchUseCase
}

// NewMapHandler creates a new MapHandler
func NewMapHandler(
	logger *zap.Logger,
	quickSearchUC *usecases.QuickSearchUseCase,
	advanceSearchUC *usecases.AdvanceSearchUseCase,
) *MapHandler {
	return &MapHandler{
		logger:          logger.With(zap.String("component", "map_handler")),
		quickSearchUC:   quickSearchUC,
		advanceSearchUC: advanceSearchUC,
	}
}

// HealthCheck handles health check requests
// @Summary Health check
// @Description Check if the service is healthy
// @Tags health
// @Produce json
// @Success 200 {object} models.HealthCheckResponse
// @Router /health [get]
func (h *MapHandler) HealthCheck(c *gin.Context) {
	response := models.HealthCheckResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Dependencies: map[string]string{
			"redis": "connected",
		},
	}
	c.JSON(http.StatusOK, response)
}

// QuickSearch handles quick search requests
// @Summary Quick search for a place
// @Description Get place details by place ID
// @Tags map
// @Accept json
// @Produce json
// @Param request body models.QuickSearchRequest true "Quick search request"
// @Success 200 {object} models.QuickSearchResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/map/quick_search [post]
func (h *MapHandler) QuickSearch(c *gin.Context) {
	// 1. Bind and validate request
	var req models.QuickSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:     "invalid_request",
			Message:   fmt.Sprintf("Invalid request: %v", err),
			Code:      http.StatusBadRequest,
			Timestamp: time.Now(),
		})
		return
	}

	// 2. Execute use case
	resp, err := h.quickSearchUC.Execute(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Quick search failed",
			zap.Error(err),
			zap.String("place_id", req.PlaceID),
		)

		// Determine status code based on error
		statusCode := http.StatusInternalServerError
		errorCode := "internal_error"

		// You can add more sophisticated error handling here
		if err.Error() == "place not found" {
			statusCode = http.StatusNotFound
			errorCode = "place_not_found"
		}

		c.JSON(statusCode, models.ErrorResponse{
			Error:     errorCode,
			Message:   err.Error(),
			Code:      statusCode,
			Timestamp: time.Now(),
		})
		return
	}

	// 3. Return success response
	c.JSON(http.StatusOK, resp)
}

// AdvanceSearch handles advance search requests
// @Summary Advance search for places
// @Description Search for places using text query and filters
// @Tags map
// @Accept json
// @Produce json
// @Param request body models.AdvanceSearchRequest true "Advance search request"
// @Success 200 {object} models.AdvanceSearchResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/map/advance_search [post]
func (h *MapHandler) AdvanceSearch(c *gin.Context) {
	// 1. Bind and validate request
	var req models.AdvanceSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:     "invalid_request",
			Message:   fmt.Sprintf("Invalid request: %v", err),
			Code:      http.StatusBadRequest,
			Timestamp: time.Now(),
		})
		return
	}

	// 2. Execute use case
	resp, err := h.advanceSearchUC.Execute(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Advance search failed",
			zap.Error(err),
			zap.String("text_query", req.TextQuery),
		)

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:     "search_failed",
			Message:   err.Error(),
			Code:      http.StatusInternalServerError,
			Timestamp: time.Now(),
		})
		return
	}

	// 3. Return success response
	c.JSON(http.StatusOK, resp)
}
