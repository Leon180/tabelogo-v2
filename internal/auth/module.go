package auth

import (
	"github.com/Leon180/tabelogo-v2/internal/auth/application"
	"github.com/Leon180/tabelogo-v2/internal/auth/infrastructure"
	authgrpc "github.com/Leon180/tabelogo-v2/internal/auth/interfaces/grpc"
	authhttp "github.com/Leon180/tabelogo-v2/internal/auth/interfaces/http"
	"github.com/Leon180/tabelogo-v2/pkg/config"
	"github.com/Leon180/tabelogo-v2/pkg/logger"
	"go.uber.org/fx"
)

// Module provides the complete auth service
var Module = fx.Module("auth",
	// Include base modules
	config.Module,
	logger.Module,

	// Include auth layers
	infrastructure.Module,
	application.Module,
	authgrpc.Module,
	authhttp.Module, // Add HTTP module
)
