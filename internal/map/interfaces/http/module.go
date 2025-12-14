package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	_ "github.com/Leon180/tabelogo-v2/internal/map/docs"
	"github.com/Leon180/tabelogo-v2/pkg/config"
	"github.com/Leon180/tabelogo-v2/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/ulule/limiter/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides HTTP interface layer dependencies
var Module = fx.Module("map.http",
	fx.Provide(
		NewMapHandler,
		NewHTTPServer,
	),
	fx.Invoke(RegisterRoutes),
)

// NewHTTPServer creates a new Gin HTTP server
func NewHTTPServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// CORS middleware - must be before routes
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Metrics middleware
	router.Use(MetricsMiddleware())

	return router
}

// RegisterRoutes registers HTTP routes and manages server lifecycle
func RegisterRoutes(
	lc fx.Lifecycle,
	router *gin.Engine,
	handler *MapHandler,
	cfg *config.Config,
	logger *zap.Logger,
	redisClient *redis.Client,
) {
	// Health check
	router.GET("/health", handler.HealthCheck)

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/map-service/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create rate limiters
	quickSearchLimiter := middleware.NewRateLimiter(redisClient, middleware.RateLimiterConfig{
		Rate:   limiter.Rate{Period: time.Minute, Limit: 60},
		Logger: logger,
	})

	advanceSearchLimiter := middleware.NewRateLimiter(redisClient, middleware.RateLimiterConfig{
		Rate:   limiter.Rate{Period: time.Minute, Limit: 30},
		Logger: logger,
	})

	// Map API routes
	mapGroup := router.Group("/api/v1/map")
	{
		// Handle OPTIONS for CORS preflight
		mapGroup.OPTIONS("/quick_search", func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		})
		mapGroup.OPTIONS("/advance_search", func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		})

		// Apply rate limiters to endpoints
		mapGroup.POST("/quick_search", quickSearchLimiter, handler.QuickSearch)
		mapGroup.POST("/advance_search", advanceSearchLimiter, handler.AdvanceSearch)
	}

	// Manage lifecycle
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info("Starting HTTP server",
					zap.Int("port", cfg.ServerPort),
					zap.String("environment", cfg.Environment),
				)
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal("Failed to serve HTTP", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping HTTP server")
			return server.Shutdown(ctx)
		},
	})
}
