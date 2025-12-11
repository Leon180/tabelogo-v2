package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Leon180/tabelogo-v2/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides HTTP interface dependencies
var Module = fx.Module("spider.http",
	fx.Provide(
		NewSpiderHandler,
		NewSSEHandler,
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

	// Configure middleware
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

	return router
}

// RegisterRoutes registers all HTTP routes
func RegisterRoutes(
	lc fx.Lifecycle,
	router *gin.Engine,
	handler *SpiderHandler,
	sseHandler *SSEHandler,
	cfg *config.Config,
	logger *zap.Logger,
) {
	// Prometheus metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "spider-service",
		})
	})

	// API routes
	api := router.Group("/api/v1/spider")
	{
		// Handle OPTIONS for CORS preflight - must match each specific route
		api.OPTIONS("/scrape", func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		})
		api.OPTIONS("/jobs/:job_id", func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		})
		api.OPTIONS("/jobs/:job_id/stream", func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		})

		api.POST("/scrape", handler.Scrape)
		api.GET("/jobs/:job_id", handler.GetJobStatus)
		api.GET("/jobs/:job_id/stream", sseHandler.StreamJobStatus) // SSE endpoint
	}

	// Lifecycle hooks
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf(":%d", cfg.ServerPort)
			logger.Info("Starting Spider Service HTTP server",
				zap.String("address", addr),
				zap.String("environment", cfg.Environment),
			)

			go func() {
				if err := router.Run(addr); err != nil {
					logger.Fatal("Failed to start HTTP server", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping Spider Service HTTP server")
			return nil
		},
	})
}
