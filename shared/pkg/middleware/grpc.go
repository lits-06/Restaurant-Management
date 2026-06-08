package middleware

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"go.uber.org/zap"
)

// LoggingInterceptor logs gRPC requests and responses
func LoggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Call the handler
		resp, err := handler(ctx, req)

		// Log the request
		duration := time.Since(start)
		statusCode := codes.OK
		if err != nil {
			statusCode = status.Code(err)
		}

		logger.Info("gRPC request",
			zap.String("method", info.FullMethod),
			zap.Duration("duration", duration),
			zap.String("status", statusCode.String()),
		)

		if err != nil {
			logger.Error("gRPC error",
				zap.String("method", info.FullMethod),
				zap.Error(err),
			)
		}

		return resp, err
	}
}

// RecoveryInterceptor recovers from panics and returns internal error
func RecoveryInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered",
					zap.String("method", info.FullMethod),
					zap.Any("panic", r),
				)
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}

// ChainUnaryServer chains multiple unary interceptors
func ChainUnaryServer(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		chain := func(currentInter grpc.UnaryServerInterceptor, currentHandler grpc.UnaryHandler) grpc.UnaryHandler {
			return func(ctx context.Context, req interface{}) (interface{}, error) {
				return currentInter(ctx, req, info, currentHandler)
			}
		}

		chainedHandler := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			chainedHandler = chain(interceptors[i], chainedHandler)
		}

		return chainedHandler(ctx, req)
	}
}
