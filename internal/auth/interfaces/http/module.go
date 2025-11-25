package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Leon180/tabelogo-v2/internal/auth/docs"
	"github.com/Leon180/tabelogo-v2/pkg/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides HTTP interface layer dependencies
var Module = fx.Module("auth.http",
	fx.Provide(
		NewAuthHandler,
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

	return router
}

// RegisterRoutes registers HTTP routes and manages server lifecycle
func RegisterRoutes(
	lc fx.Lifecycle,
	router *gin.Engine,
	handler *AuthHandler,
	cfg *config.Config,
	logger *zap.Logger,
) {
	// Register auth routes
	authGroup := router.Group("/api/v1/auth")
	{
		// Handle OPTIONS for CORS preflight - must match each specific route
		authGroup.OPTIONS("/register", func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		})
		authGroup.OPTIONS("/login", func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		})
		authGroup.OPTIONS("/refresh", func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		})
		authGroup.OPTIONS("/validate", func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		})

		authGroup.POST("/register", handler.Register)
		authGroup.POST("/login", handler.Login)
		authGroup.POST("/refresh", handler.RefreshToken)
		authGroup.GET("/validate", handler.ValidateToken)
	}

	// Swagger documentation endpoints
	// Use service-specific path (/auth-service/swagger/) to avoid conflicts with other services
	router.GET("/auth-service/swagger/doc.json", func(c *gin.Context) {
		c.String(http.StatusOK, docs.SwaggerInfo.ReadDoc())
	})

	router.GET("/auth-service/swagger/index.html", func(c *gin.Context) {
		// Read file directly to avoid http.ServeFile's automatic redirects
		absPath, err := filepath.Abs("./internal/auth/docs/index.html")
		if err != nil {
			logger.Error("Failed to resolve Swagger UI path", zap.Error(err))
			c.String(http.StatusInternalServerError, "Internal server error")
			return
		}

		content, err := os.ReadFile(absPath)
		if err != nil {
			logger.Error("Failed to read Swagger UI file", zap.Error(err))
			c.String(http.StatusNotFound, "Swagger UI not found")
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", content)
	})

	// Redirect shortcuts for convenience
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/auth-service/swagger/index.html")
	})
	router.GET("/auth-service/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/auth-service/swagger/index.html")
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

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
