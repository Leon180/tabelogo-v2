package http

import (
	"fmt"
	"net/http"

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
		if err == domainerrors.ErrRestaurantAlreadyExists {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "restaurant_exists",
				Message: "Restaurant already exists",
			})
			return
		}
		if err == domainerrors.ErrInvalidLocation {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "invalid_location",
				Message: "Invalid location coordinates",
			})
			return
		}
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
// @Summary Get restaurant by ID
// @Description Get restaurant details by ID
// @Tags restaurants
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

// SearchRestaurants godoc
// @Summary Search restaurants
// @Description Search restaurants by query string
// @Tags restaurants
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} RestaurantListResponse
// @Failure 400 {object} ErrorResponse
// @Router /restaurants/search [get]
func (h *RestaurantHandler) SearchRestaurants(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "missing_query",
			Message: "Search query is required",
		})
		return
	}

	limit := 20
	offset := 0
	if l := c.Query("limit"); l != "" {
		if _, err := fmt.Sscanf(l, "%d", &limit); err != nil {
			limit = 20
		}
	}
	if o := c.Query("offset"); o != "" {
		if _, err := fmt.Sscanf(o, "%d", &offset); err != nil {
			offset = 0
		}
	}

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
// @Failure 409 {object} ErrorResponse
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
		if err == domainerrors.ErrFavoriteAlreadyExists {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "favorite_exists",
				Message: "Already in favorites",
			})
			return
		}
		h.logger.Error("Failed to add to favorites", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to add to favorites",
		})
		return
	}

	c.JSON(http.StatusCreated, FavoriteResponse{
		Favorite: toFavoriteDTO(favorite),
	})
}

// GetUserFavorites godoc
// @Summary Get user's favorites
// @Description Get all favorites for a user
// @Tags favorites
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
