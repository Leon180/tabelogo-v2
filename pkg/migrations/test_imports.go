// +build ignore

// This file is used to test imports
package migrations

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Test that all imports work
func testImports() {
	var _ *sql.DB
	var _ *migrate.Migrate
	var _ *postgres.Config
	var _ fx.Option
	var _ *zap.Logger
}
