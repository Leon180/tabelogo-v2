package grpc

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// ConnectionConfig holds gRPC connection configuration
type ConnectionConfig struct {
	Address          string
	Timeout          time.Duration
	MaxRetries       int
	KeepAliveTime    time.Duration
	KeepAliveTimeout time.Duration
}

// NewMapServiceConnection creates a new gRPC connection to Map Service
func NewMapServiceConnection(cfg *ConnectionConfig, logger *zap.Logger) (*grpc.ClientConn, error) {
	logger.Info("Connecting to Map Service",
		zap.String("address", cfg.Address),
		zap.Duration("timeout", cfg.Timeout),
	)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		cfg.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                cfg.KeepAliveTime,
			Timeout:             cfg.KeepAliveTimeout,
			PermitWithoutStream: true,
		}),
		grpc.WithUnaryInterceptor(loggingInterceptor(logger)),
	)
	if err != nil {
		logger.Error("Failed to connect to Map Service",
			zap.Error(err),
			zap.String("address", cfg.Address),
		)
		return nil, err
	}

	logger.Info("Successfully connected to Map Service")
	return conn, nil
}

// loggingInterceptor logs gRPC calls
func loggingInterceptor(logger *zap.Logger) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()

		err := invoker(ctx, method, req, reply, cc, opts...)

		duration := time.Since(start)
		if err != nil {
			logger.Error("gRPC call failed",
				zap.String("method", method),
				zap.Duration("duration", duration),
				zap.Error(err),
			)
		} else {
			logger.Debug("gRPC call succeeded",
				zap.String("method", method),
				zap.Duration("duration", duration),
			)
		}

		return err
	}
}
