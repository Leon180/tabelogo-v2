package main

import (
	"github.com/Leon180/tabelogo-v2/internal/restaurant"
	"go.uber.org/fx"
)

// @title Restaurant Service API
// @version 1.0
// @description Restaurant management service for Tabelogo platform
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@tabelogo.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:18082
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	fx.New(
		restaurant.Module,
	).Run()
}
