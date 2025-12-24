package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Leon180/tabelogo-v2/internal/restaurant/docs"
	"github.com/Leon180/tabelogo-v2/pkg/config"
	"github.com/Leon180/tabelogo-v2/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides HTTP interface dependencies
var Module = fx.Module("restaurant.http",
	fx.Provide(
		NewRestaurantHandler,
		NewHTTPServer,
		NewAuthMiddleware,
	),
	fx.Invoke(RegisterRoutes),
)

// NewAuthMiddleware creates auth middleware for Restaurant Service
func NewAuthMiddleware(cfg *config.Config, redis *redis.Client, logger *zap.Logger) *middleware.AuthMiddleware {
	return middleware.NewAuthMiddleware(cfg.JWT.Secret, redis, logger)
}

// NewHTTPServer creates a new HTTP server
func NewHTTPServer(cfg *config.Config) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Configure middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(CORSMiddleware()) // Enable CORS for frontend
	router.Use(MetricsMiddleware())

	return router
}

// RegisterRoutes registers all HTTP routes
func RegisterRoutes(
	lc fx.Lifecycle,
	router *gin.Engine,
	handler *RestaurantHandler,
	authMW *middleware.AuthMiddleware,
	cfg *config.Config,
	logger *zap.Logger,
) {
	// Swagger documentation endpoints
	// Set base path for Swagger docs
	docs.SwaggerInfo.BasePath = "/api/v1"

	// Serve Swagger UI at /swagger/*
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Also serve at service-specific path for consistency
	router.GET("/restaurant-service/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Prometheus metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public restaurant routes (read-only, optional auth for tracking)
		publicRestaurants := v1.Group("/restaurants")
		publicRestaurants.Use(authMW.Optional()) // Track authenticated users but don't require auth
		{
			publicRestaurants.GET("/:id", handler.GetRestaurant)
			publicRestaurants.GET("/search", handler.SearchRestaurants)
			publicRestaurants.GET("/quick-search/:place_id", handler.QuickSearchByPlaceID)
		}

		// Protected restaurant routes (write operations, admin only)
		protectedRestaurants := v1.Group("/restaurants")
		protectedRestaurants.Use(authMW.RequireAuth())
		{
			// Admin-only operations
			protectedRestaurants.POST("", authMW.RequireRole("admin"), handler.CreateRestaurant)
			// Allow authenticated users to update restaurant details (e.g., Japanese name)
			protectedRestaurants.PATCH("/:id", handler.UpdateRestaurant)
		}

		// Protected favorite routes (require authentication)
		favorites := v1.Group("/favorites")
		favorites.Use(authMW.RequireAuth())
		{
			favorites.POST("", handler.AddToFavorites)
		}

		// Protected user favorites routes (require authentication)
		userFavorites := v1.Group("/users/:userId/favorites")
		userFavorites.Use(authMW.RequireAuth())
		{
			userFavorites.GET("", handler.GetUserFavorites)
		}
	}

	// Lifecycle hooks
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf(":%d", cfg.ServerPort)
			logger.Info("Starting HTTP server",
				zap.String("addr", addr),
				zap.String("swagger", fmt.Sprintf("http://localhost:%d/swagger", cfg.ServerPort)),
				zap.String("metrics", fmt.Sprintf("http://localhost:%d/metrics", cfg.ServerPort)),
			)

			go func() {
				if err := router.Run(addr); err != nil && err != http.ErrServerClosed {
					logger.Fatal("Failed to start HTTP server", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping HTTP server")
			return nil
		},
	})
}
