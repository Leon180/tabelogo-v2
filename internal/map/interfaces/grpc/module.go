package grpc

import (
	"context"
	"fmt"
	"net"

	mapv1 "github.com/Leon180/tabelogo-v2/api/gen/map/v1"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Module provides gRPC server dependencies
var Module = fx.Module("grpc",
	fx.Provide(NewServer),
	fx.Invoke(registerHooks),
)

// Params holds dependencies for gRPC server
type Params struct {
	fx.In
	Server *Server
	Logger *zap.Logger
}

// registerHooks registers lifecycle hooks for gRPC server
func registerHooks(lc fx.Lifecycle, params Params) {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor(params.Logger)),
	)

	// Register the MapService server
	mapv1.RegisterMapServiceServer(grpcServer, params.Server)

	// Register reflection service (for grpcurl and debugging)
	reflection.Register(grpcServer)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Get gRPC port from environment (default: 19083)
			port := "19083"
			// TODO: Get from config

			listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
			if err != nil {
				return fmt.Errorf("failed to listen on port %s: %w", port, err)
			}

			params.Logger.Info("Starting gRPC server",
				zap.String("port", port),
			)

			// Start server in goroutine
			go func() {
				if err := grpcServer.Serve(listener); err != nil {
					params.Logger.Fatal("gRPC server failed", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			params.Logger.Info("Stopping gRPC server")
			grpcServer.GracefulStop()
			return nil
		},
	})
}

// loggingInterceptor logs gRPC requests
func loggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		logger.Debug("gRPC request",
			zap.String("method", info.FullMethod),
		)

		resp, err := handler(ctx, req)
		if err != nil {
			logger.Error("gRPC request failed",
				zap.String("method", info.FullMethod),
				zap.Error(err),
			)
		}

		return resp, err
	}
}
