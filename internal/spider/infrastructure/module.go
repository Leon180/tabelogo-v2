package infrastructure

import (
	"github.com/Leon180/tabelogo-v2/internal/spider/domain/repositories"
	"github.com/Leon180/tabelogo-v2/internal/spider/infrastructure/persistence"
	"github.com/Leon180/tabelogo-v2/internal/spider/infrastructure/scraper"
	"go.uber.org/fx"
)

// Module provides infrastructure layer dependencies
var Module = fx.Module("spider.infrastructure",
	fx.Provide(
		// Repositories
		fx.Annotate(
			persistence.NewInMemoryJobRepository,
			fx.As(new(repositories.JobRepository)),
		),
		// Scraper
		scraper.NewScraper,
	),
)
