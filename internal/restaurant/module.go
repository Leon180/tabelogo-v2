package restaurant

import (
	"github.com/Leon180/tabelogo-v2/internal/restaurant/application"
	"github.com/Leon180/tabelogo-v2/internal/restaurant/infrastructure"
	restaurantgrpc "github.com/Leon180/tabelogo-v2/internal/restaurant/interfaces/grpc"
	restauranthttp "github.com/Leon180/tabelogo-v2/internal/restaurant/interfaces/http"
	"github.com/Leon180/tabelogo-v2/pkg/config"
	"github.com/Leon180/tabelogo-v2/pkg/logger"
	"go.uber.org/fx"
)

// Module provides all restaurant service dependencies
var Module = fx.Module("restaurant",
	config.Module,
	logger.Module,
	infrastructure.Module,
	application.Module,
	restaurantgrpc.Module,
	restauranthttp.Module,
)
