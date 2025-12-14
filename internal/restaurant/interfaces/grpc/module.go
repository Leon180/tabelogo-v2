package grpc

import (
	"context"
	"fmt"
	"net"

	restaurantv1 "github.com/Leon180/tabelogo-v2/api/gen/restaurant/v1"
	"github.com/Leon180/tabelogo-v2/pkg/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Module provides gRPC server dependencies
var Module = fx.Module("restaurant.grpc",
	fx.Provide(NewRestaurantServer),
	fx.Invoke(RegisterGRPCServer),
)

// RegisterGRPCServer registers and starts the gRPC server
func RegisterGRPCServer(
	lc fx.Lifecycle,
	server *RestaurantServer,
	cfg *config.Config,
	logger *zap.Logger,
) {
	grpcServer := grpc.NewServer()
	restaurantv1.RegisterRestaurantServiceServer(grpcServer, server)

	// Enable reflection for grpcurl
	reflection.Register(grpcServer)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf(":%d", cfg.GRPCPort)
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				return fmt.Errorf("failed to listen on %s: %w", addr, err)
			}

			go func() {
				logger.Info("gRPC server starting", zap.String("addr", addr))
				if err := grpcServer.Serve(lis); err != nil {
					logger.Error("gRPC server failed", zap.Error(err))
				}
			}()

			logger.Info("gRPC server started successfully", zap.String("addr", addr))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping gRPC server")
			grpcServer.GracefulStop()
			logger.Info("gRPC server stopped")
			return nil
		},
	})
}
