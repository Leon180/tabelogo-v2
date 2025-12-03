package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Leon180/tabelogo-v2/pkg/config"
	"github.com/gin-gonic/gin"
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

	router := gin.Default()
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

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// Lifecycle hooks
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf(":%d", cfg.ServerPort)
			logger.Info("Starting HTTP server", zap.String("addr", addr))

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
