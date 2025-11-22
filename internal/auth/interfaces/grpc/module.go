package grpc

import (
	"context"
	"fmt"
	"net"

	authv1 "github.com/Leon180/tabelogo-v2/api/gen/auth/v1"
	"github.com/Leon180/tabelogo-v2/pkg/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Module provides gRPC interface layer dependencies
var Module = fx.Module("auth.grpc",
	fx.Provide(
		NewGRPCServer,
		NewAuthServer,
	),
	fx.Invoke(RegisterServer),
)

// NewGRPCServer creates a new gRPC server
func NewGRPCServer() *grpc.Server {
	return grpc.NewServer()
}

// RegisterServer registers the auth service with the gRPC server and manages its lifecycle
func RegisterServer(
	lc fx.Lifecycle,
	server *grpc.Server,
	authServer *AuthServer,
	cfg *config.Config,
	logger *zap.Logger,
) {
	// Register services
	authv1.RegisterAuthServiceServer(server, authServer)
	reflection.Register(server)

	// Manage lifecycle
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
			if err != nil {
				return fmt.Errorf("failed to listen on port %d: %w", cfg.GRPCPort, err)
			}

			go func() {
				logger.Info("Starting gRPC server",
					zap.Int("port", cfg.GRPCPort),
					zap.String("environment", cfg.Environment),
				)
				if err := server.Serve(listener); err != nil {
					logger.Fatal("Failed to serve gRPC", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping gRPC server")
			server.GracefulStop()
			logger.Info("gRPC server stopped")
			return nil
		},
	})
}
