package log

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type loggerKey string

// LoggerKey is the key for the logger in the context
const LoggerKey loggerKey = "logger"

// TODO unit test
// LoggerInterceptor is the interceptor for the logger
func CreateLoggerInterceptor(
	logFactory LogFactoryer,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		logger, err := logFactory.NewLoggerWithCorrelationID(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Internal server error. Dubious request")
		}
		newCtx := context.WithValue(ctx, LoggerKey, logger)
		return handler(newCtx, req)
	}
}

// GetLoggerFromContext returns the logger from the context
func GetLoggerFromContext(ctx context.Context) Loggerer {
	if logger, ok := ctx.Value(LoggerKey).(Loggerer); ok {
		return logger
	}
	return nil
}

// TODO unit test
// AddCorrelationIDToContext adds the correlation ID to the context
func AddCorrelationIDToContext(ctx context.Context, correlationID string) context.Context {
	md := metadata.New(map[string]string{
		CorrelationIDKey: correlationID,
	})
	return metadata.NewOutgoingContext(ctx, md)
}
