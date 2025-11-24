package main

import (
	"github.com/Leon180/tabelogo-v2/internal/auth"
	"go.uber.org/fx"
)

// @title Tabelogo Auth Service API
// @version 1.0
// @description Authentication and user management service for Tabelogo platform
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@tabelogo.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8081
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	fx.New(
		// Load the complete auth module
		auth.Module,
	).Run()
}
