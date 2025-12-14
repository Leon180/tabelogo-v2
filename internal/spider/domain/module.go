package domain

import (
	"github.com/Leon180/tabelogo-v2/internal/spider/domain/models"
	"go.uber.org/fx"
)

// Module provides domain layer dependencies
var Module = fx.Module("spider.domain",
	fx.Provide(
		models.NewScraperConfig,
	),
)
