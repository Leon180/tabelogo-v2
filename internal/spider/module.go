package spider

import (
	"github.com/Leon180/tabelogo-v2/internal/spider/application"
	"github.com/Leon180/tabelogo-v2/internal/spider/domain"
	"github.com/Leon180/tabelogo-v2/internal/spider/infrastructure"
	"github.com/Leon180/tabelogo-v2/internal/spider/interfaces/http"
	"github.com/Leon180/tabelogo-v2/pkg/config"
	"github.com/Leon180/tabelogo-v2/pkg/logger"
	"go.uber.org/fx"
)

// Module provides all spider service dependencies
var Module = fx.Module("spider",
	// Configuration and logging
	config.Module,
	logger.Module,

	// Domain layer
	domain.Module,

	// Infrastructure layer
	infrastructure.Module,

	// Application layer
	application.Module,

	// Interface layer
	http.Module,
)
