package grpc

import (
	"context"
	"fmt"
	"net"

	spiderv1 "github.com/Leon180/tabelogo-v2/api/gen/spider/v1"
	"github.com/Leon180/tabelogo-v2/pkg/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Module provides gRPC server dependencies
var Module = fx.Module("grpc",
	fx.Provide(
		NewSpiderServer,
		NewGRPCServer,
	),
	fx.Invoke(RegisterServer),
)

// NewGRPCServer creates a new gRPC server
func NewGRPCServer(logger *zap.Logger) *grpc.Server {
	server := grpc.NewServer()

	// Enable reflection for grpcurl and other tools
	reflection.Register(server)

	logger.Info("gRPC server created")
	return server
}

// RegisterServer registers the Spider service with the gRPC server
func RegisterServer(
	grpcServer *grpc.Server,
	spiderServer *SpiderServer,
	cfg *config.Config,
	logger *zap.Logger,
	lc fx.Lifecycle,
) {
	// Register service
	spiderv1.RegisterSpiderServiceServer(grpcServer, spiderServer)
	logger.Info("Spider service registered with gRPC server")

	// Lifecycle hooks
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf(":%d", cfg.GRPCPort)
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				return fmt.Errorf("failed to listen on %s: %w", addr, err)
			}

			logger.Info("Starting Spider Service gRPC server",
				zap.String("address", addr),
			)

			go func() {
				if err := grpcServer.Serve(lis); err != nil {
					logger.Fatal("Failed to serve gRPC", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping Spider Service gRPC server")
			grpcServer.GracefulStop()
			return nil
		},
	})
}
