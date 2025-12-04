package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Leon180/tabelogo-v2/internal/restaurant/docs"
	"github.com/Leon180/tabelogo-v2/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides HTTP interface dependencies
var Module = fx.Module("restaurant.http",
	fx.Provide(
		NewRestaurantHandler,
		NewHTTPServer,
	),
	fx.Invoke(RegisterRoutes),
)

// NewHTTPServer creates a new HTTP server
func NewHTTPServer(cfg *config.Config) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add recovery middleware
	router.Use(gin.Recovery())

	// Add logger middleware
	router.Use(gin.Logger())

	// Add metrics middleware
	router.Use(MetricsMiddleware())

	return router
}

// RegisterRoutes registers all HTTP routes
func RegisterRoutes(
	lc fx.Lifecycle,
	router *gin.Engine,
	handler *RestaurantHandler,
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
		// Restaurant routes
		restaurants := v1.Group("/restaurants")
		{
			restaurants.POST("", handler.CreateRestaurant)
			restaurants.GET("/:id", handler.GetRestaurant)
			restaurants.GET("/search", handler.SearchRestaurants)
		}

		// Favorite routes
		favorites := v1.Group("/favorites")
		{
			favorites.POST("", handler.AddToFavorites)
		}

		// User favorites routes
		v1.GET("/users/:userId/favorites", handler.GetUserFavorites)
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
