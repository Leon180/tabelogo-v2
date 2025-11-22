package application

import (
	"github.com/Leon180/tabelogo-v2/pkg/config"
	"github.com/Leon180/tabelogo-v2/pkg/jwt"
	"go.uber.org/fx"
)

// Module provides application layer dependencies
var Module = fx.Module("auth.application",
	fx.Provide(
		NewJWTMaker,
		NewAuthService,
	),
)

// NewJWTMaker creates a new JWT maker
func NewJWTMaker(cfg *config.Config) (jwt.Maker, error) {
	return jwt.NewJWTMaker(cfg.JWT.Secret)
}
