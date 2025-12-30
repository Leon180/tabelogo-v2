package http

import (
	"net/http"

	"github.com/Leon180/tabelogo-v2/internal/auth/application"
	"github.com/Leon180/tabelogo-v2/internal/auth/domain/errors"
	"github.com/Leon180/tabelogo-v2/internal/auth/domain/model"
	"github.com/Leon180/tabelogo-v2/pkg/metrics"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service application.AuthService
	logger  *zap.Logger
}

func NewAuthHandler(service application.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		service: service,
		logger:  logger,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register request"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.AuthRegisterTotal.WithLabelValues("failed").Inc()
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	user, err := h.service.Register(c.Request.Context(), req.Email, req.Password, req.Username)
	if err != nil {
		if err == errors.ErrEmailAlreadyExists {
			metrics.AuthRegisterTotal.WithLabelValues("email_exists").Inc()
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "email_exists",
				Message: "Email already registered",
			})
			return
		}
		metrics.AuthRegisterTotal.WithLabelValues("failed").Inc()
		h.logger.Error("Failed to register user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to register user",
		})
		return
	}

	metrics.AuthRegisterTotal.WithLabelValues("success").Inc()
	c.JSON(http.StatusCreated, RegisterResponse{
		User: toUserResponse(user),
	})
}

// Login godoc
// @Summary Login user
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	timer := metrics.NewTimer(metrics.AuthLoginDuration)
	defer timer.ObserveDuration()

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.AuthLoginTotal.WithLabelValues("failed").Inc()
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Extract device info from User-Agent header
	deviceInfo := c.GetHeader("User-Agent")
	if deviceInfo == "" {
		deviceInfo = "unknown"
	}

	// Extract IP address
	ipAddress := c.ClientIP()

	accessToken, refreshToken, err := h.service.Login(
		c.Request.Context(),
		req.Email,
		req.Password,
		deviceInfo,
		ipAddress,
		req.RememberMe,
	)
	if err != nil {
		if err == errors.ErrUserNotFound {
			metrics.AuthLoginTotal.WithLabelValues("failed").Inc()
			metrics.AuthFailedLoginAttempts.WithLabelValues("user_not_found").Inc()
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "invalid_credentials",
				Message: "Invalid email or password",
			})
			return
		}
		if err == errors.ErrInvalidPassword {
			metrics.AuthLoginTotal.WithLabelValues("failed").Inc()
			metrics.AuthFailedLoginAttempts.WithLabelValues("wrong_password").Inc()
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "invalid_credentials",
				Message: "Invalid email or password",
			})
			return
		}
		metrics.AuthLoginTotal.WithLabelValues("failed").Inc()
		h.logger.Error("Failed to login", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to login",
		})
		return
	}

	// Get user info for response
	user, err := h.service.ValidateToken(c.Request.Context(), accessToken)
	if err != nil {
		h.logger.Error("Failed to validate token after login", zap.Error(err))
	}

	metrics.AuthLoginTotal.WithLabelValues("success").Inc()
	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         toUserResponse(user),
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} RefreshTokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.AuthTokenRefreshTotal.WithLabelValues("failed").Inc()
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	accessToken, refreshToken, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		metrics.AuthTokenRefreshTotal.WithLabelValues("failed").Inc()
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "invalid_token",
			Message: "Invalid or expired refresh token",
		})
		return
	}

	metrics.AuthTokenRefreshTotal.WithLabelValues("success").Inc()
	c.JSON(http.StatusOK, RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// ValidateToken godoc
// @Summary Validate access token
// @Description Validate an access token
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} ValidateTokenResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/validate [get]
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		metrics.AuthTokenValidationTotal.WithLabelValues("invalid").Inc()
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "missing_token",
			Message: "Authorization header is required",
		})
		return
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	user, err := h.service.ValidateToken(c.Request.Context(), token)
	if err != nil {
		metrics.AuthTokenValidationTotal.WithLabelValues("invalid").Inc()
		c.JSON(http.StatusUnauthorized, ValidateTokenResponse{
			Valid: false,
		})
		return
	}

	metrics.AuthTokenValidationTotal.WithLabelValues("valid").Inc()
	c.JSON(http.StatusOK, ValidateTokenResponse{
		Valid: true,
		User:  toUserResponse(user),
	})
}

func toUserResponse(user *model.User) *UserResponse {
	if user == nil {
		return nil
	}
	return &UserResponse{
		ID:            user.ID().String(),
		Email:         user.Email(),
		Username:      user.Username(),
		Role:          string(user.Role()),
		EmailVerified: user.EmailVerified(),
		CreatedAt:     user.CreatedAt(),
	}
}
