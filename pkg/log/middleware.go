package log

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type loggerKey string

// LoggerKey is the key for the logger in the context
const LoggerKey loggerKey = "logger"

// TODO unit test

// AddNewCorrelationIDToContext can ber used as a middleware that adds a new correlation ID to the context
func AddNewCorrelationIDToContext(context *gin.Context) {
	// Generate a random UUID
	correlationID := uuid.New().String()

	ctxWithCorrelationID := AddCorrelationIDToOutgoingContext(context.Request.Context(), correlationID)
	context.Request = context.Request.WithContext(ctxWithCorrelationID)

	// Continue processing the request
	context.Next()
}

// CreateLoggerInterceptor is the interceptor that intercepts the gRPC
// calls and adds a logger with a correlation ID to the context
func CreateLoggerInterceptor(
	logFactory Factoryer,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		logger, err := logFactory.NewLoggerWithCorrelationID(ctx)
		if err != nil {
			return nil, err
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

// AddCorrelationIDToOutgoingContext adds the correlation ID to the context
func AddCorrelationIDToOutgoingContext(ctx context.Context, correlationID string) context.Context {
	md := metadata.New(map[string]string{
		CorrelationIDKey: correlationID,
	})
	return metadata.NewOutgoingContext(ctx, md)
}

// AddCorrelationIDToIncomingContext adds the correlation ID to the context
func AddCorrelationIDToIncomingContext(ctx context.Context, correlationID string) context.Context {
	md := metadata.New(map[string]string{
		CorrelationIDKey: correlationID,
	})
	return metadata.NewIncomingContext(ctx, md)
}

// GetCorrelationIDFromContext returns the correlation ID obtained from the context
func GetCorrelationIDFromContext(ctx context.Context) (*string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("Metadata not found in context")
	}
	correlationIDs, exists := md[CorrelationIDKey]
	if !exists || len(correlationIDs) != 1 {
		return nil, errors.New("Correlation ID not found in metadata")
	}
	return &correlationIDs[0], nil
}
