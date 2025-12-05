package mapservice

import (
	"github.com/Leon180/tabelogo-v2/internal/map/application"
	"github.com/Leon180/tabelogo-v2/internal/map/infrastructure"
	mapgrpc "github.com/Leon180/tabelogo-v2/internal/map/interfaces/grpc"
	maphttp "github.com/Leon180/tabelogo-v2/internal/map/interfaces/http"
	"github.com/Leon180/tabelogo-v2/pkg/config"
	"github.com/Leon180/tabelogo-v2/pkg/logger"
	"go.uber.org/fx"
)

// Module provides the complete map service
var Module = fx.Module("map",
	// Include base modules
	config.Module,
	logger.Module,

	// Include map layers
	infrastructure.Module,
	application.Module,
	maphttp.Module,
	mapgrpc.Module,
)
