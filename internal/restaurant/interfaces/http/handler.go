package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Leon180/tabelogo-v2/internal/restaurant/application"
	domainerrors "github.com/Leon180/tabelogo-v2/internal/restaurant/domain/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type RestaurantHandler struct {
	service application.RestaurantService
	logger  *zap.Logger
}

func NewRestaurantHandler(service application.RestaurantService, logger *zap.Logger) *RestaurantHandler {
	return &RestaurantHandler{
		service: service,
		logger:  logger,
	}
}

// CreateRestaurant godoc
// @Summary Create a new restaurant
// @Description Create a new restaurant with details
// @Tags restaurants
// @Accept json
// @Produce json
// @Param request body CreateRestaurantRequest true "Create restaurant request"
// @Success 201 {object} RestaurantResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /restaurants [post]
func (h *RestaurantHandler) CreateRestaurant(c *gin.Context) {
	var req CreateRestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	appReq := application.CreateRestaurantRequest{
		Name:        req.Name,
		NameJa:      req.NameJa,
		Source:      req.Source,
		ExternalID:  req.ExternalID,
		Address:     req.Address,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Rating:      req.Rating,
		PriceRange:  req.PriceRange,
		CuisineType: req.CuisineType,
		Phone:       req.Phone,
		Website:     req.Website,
	}

	restaurant, err := h.service.CreateRestaurant(c.Request.Context(), appReq)
	if err != nil {
		h.logger.Error("Failed to create restaurant", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create restaurant",
		})
		return
	}

	c.JSON(http.StatusCreated, RestaurantResponse{
		Restaurant: toRestaurantDTO(restaurant),
	})
}

// GetRestaurant godoc
// @Summary Get a restaurant by ID
// @Description Get restaurant details by ID
// @Tags restaurants
// @Accept json
// @Produce json
// @Param id path string true "Restaurant ID"
// @Success 200 {object} RestaurantResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /restaurants/{id} [get]
func (h *RestaurantHandler) GetRestaurant(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid restaurant ID",
		})
		return
	}

	restaurant, err := h.service.GetRestaurant(c.Request.Context(), id)
	if err != nil {
		if err == domainerrors.ErrRestaurantNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "Restaurant not found",
			})
			return
		}

		h.logger.Error("Failed to get restaurant", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get restaurant",
		})
		return
	}

	c.JSON(http.StatusOK, RestaurantResponse{
		Restaurant: toRestaurantDTO(restaurant),
	})
}

// UpdateRestaurant godoc
// @Summary Update a restaurant
// @Description Update restaurant details (currently supports Japanese name)
// @Tags restaurants
// @Accept json
// @Produce json
// @Param id path string true "Restaurant ID"
// @Param request body UpdateRestaurantRequest true "Update request"
// @Success 200 {object} RestaurantResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /restaurants/{id} [patch]
func (h *RestaurantHandler) UpdateRestaurant(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid restaurant ID",
		})
		return
	}

	var req UpdateRestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Map HTTP DTO to Application DTO
	appReq := application.UpdateRestaurantRequest{
		NameJa: req.NameJa,
	}

	restaurant, err := h.service.UpdateRestaurant(c.Request.Context(), id, appReq)
	if err != nil {
		if err == domainerrors.ErrRestaurantNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "Restaurant not found",
			})
			return
		}

		h.logger.Error("Failed to update restaurant", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to update restaurant",
		})
		return
	}

	c.JSON(http.StatusOK, RestaurantResponse{
		Restaurant: toRestaurantDTO(restaurant),
	})
}

// SearchRestaurants godoc
// @Summary Search restaurants
// @Description Search restaurants by query string
// @Tags restaurants
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} RestaurantListResponse
// @Failure 400 {object} ErrorResponse
// @Router /restaurants/search [get]
func (h *RestaurantHandler) SearchRestaurants(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_query",
			Message: "Search query is required",
		})
		return
	}

	limit := 10
	offset := 0

	restaurants, err := h.service.SearchRestaurants(c.Request.Context(), query, limit, offset)
	if err != nil {
		h.logger.Error("Failed to search restaurants", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to search restaurants",
		})
		return
	}

	c.JSON(http.StatusOK, RestaurantListResponse{
		Restaurants: toRestaurantDTOList(restaurants),
		Total:       len(restaurants),
	})
}

// AddToFavorites godoc
// @Summary Add restaurant to favorites
// @Description Add a restaurant to user's favorites
// @Tags favorites
// @Accept json
// @Produce json
// @Param request body AddFavoriteRequest true "Add favorite request"
// @Success 201 {object} FavoriteResponse
// @Failure 400 {object} ErrorResponse
// @Router /favorites [post]
func (h *RestaurantHandler) AddToFavorites(c *gin.Context) {
	var req AddFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_user_id",
			Message: "Invalid user ID",
		})
		return
	}

	restaurantID, err := uuid.Parse(req.RestaurantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_restaurant_id",
			Message: "Invalid restaurant ID",
		})
		return
	}

	favorite, err := h.service.AddToFavorites(c.Request.Context(), userID, restaurantID)
	if err != nil {
		h.logger.Error("Failed to add favorite", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to add favorite",
		})
		return
	}

	c.JSON(http.StatusCreated, FavoriteResponse{
		Favorite: toFavoriteDTO(favorite),
	})
}

// GetUserFavorites godoc
// @Summary Get user favorites
// @Description Get all favorites for a user
// @Tags favorites
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} FavoriteListResponse
// @Failure 400 {object} ErrorResponse
// @Router /users/{userId}/favorites [get]
func (h *RestaurantHandler) GetUserFavorites(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_user_id",
			Message: "Invalid user ID",
		})
		return
	}

	favorites, err := h.service.GetUserFavorites(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user favorites", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get favorites",
		})
		return
	}

	c.JSON(http.StatusOK, FavoriteListResponse{
		Favorites: toFavoriteDTOList(favorites),
		Total:     len(favorites),
	})
}

// QuickSearchByPlaceID godoc
// @Summary Quick search restaurant by Google Place ID
// @Description Search for a restaurant using Google Place ID with cache-first strategy.
// @Description This endpoint implements a cache-first approach:
// @Description - Returns cached data if fresh (< 3 days old)
// @Description - Falls back to Map Service if cache miss or stale
// @Description - Returns stale data if Map Service fails (graceful degradation)
// @Description Response headers indicate cache status and data source.
// @Tags restaurants
// @Accept json
// @Produce json
// @Param place_id path string true "Google Place ID" example("ChIJN1t_tDeuEmsRUsoyG83frY4")
// @Success 200 {object} RestaurantResponse "Restaurant found" headers(X-Cache-Status=string,X-Data-Source=string,X-Data-Age=string)
// @Failure 400 {object} ErrorResponse "Invalid place ID"
// @Failure 404 {object} ErrorResponse "Restaurant not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Failure 503 {object} ErrorResponse "Map Service unavailable and no cached data"
// @Router /restaurants/quick-search/{place_id} [get]
func (h *RestaurantHandler) QuickSearchByPlaceID(c *gin.Context) {
	placeID := c.Param("place_id")

	// Validate place_id
	if placeID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_place_id",
			Message: "Place ID is required",
		})
		return
	}

	// Log request
	h.logger.Info("QuickSearchByPlaceID request",
		zap.String("place_id", placeID),
		zap.String("client_ip", c.ClientIP()),
	)

	// Call service
	restaurant, err := h.service.QuickSearchByPlaceID(c.Request.Context(), placeID)
	if err != nil {
		h.logger.Error("QuickSearchByPlaceID failed",
			zap.String("place_id", placeID),
			zap.Error(err),
		)

		// Handle specific errors
		if err == domainerrors.ErrRestaurantNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "Restaurant not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to search restaurant",
		})
		return
	}

	// Calculate data age
	dataAge := time.Since(restaurant.UpdatedAt())
	cacheStatus := "HIT"
	dataSource := "CACHE"

	// If data is very fresh (< 1 second), it's likely from Map Service
	if dataAge < time.Second {
		cacheStatus = "MISS"
		dataSource = "MAP_SERVICE"
	}

	// Add cache status headers
	c.Header("X-Cache-Status", cacheStatus)
	c.Header("X-Data-Source", dataSource)
	c.Header("X-Data-Age", fmt.Sprintf("%.0fs", dataAge.Seconds()))

	// Log success
	h.logger.Info("QuickSearchByPlaceID success",
		zap.String("place_id", placeID),
		zap.String("restaurant_id", restaurant.ID().String()),
		zap.String("cache_status", cacheStatus),
		zap.Float64("data_age_seconds", dataAge.Seconds()),
	)

	c.JSON(http.StatusOK, RestaurantResponse{
		Restaurant: toRestaurantDTO(restaurant),
	})
}
