package log

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type loggerKey string

// LoggerKey is the key for the logger in the context
const LoggerKey loggerKey = "logger"

// TODO: unit test

// CreateGinLoggerMiddleware is the middleware that adds a logger with a correlation ID to the context
func CreateGinLoggerMiddleware(logFactory Factoryer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a new context with the logger
		logger, err := logFactory.NewLoggerWithCorrelationID(c.Request.Context())
		if err != nil {
			// Handle error if needed
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		newCtx := context.WithValue(c.Request.Context(), LoggerKey, logger)

		// Set the new context in the Gin request
		c.Request = c.Request.WithContext(newCtx)

		// Continue processing the request
		c.Next()
	}
}

// AddNewCorrelationIDToContext can ber used as a middleware that adds a new correlation ID to the context
func AddNewCorrelationIDToContext(context *gin.Context) {
	// Generate a random UUID
	correlationID := uuid.New().String()

	ctxWithCorrelationID := AddCorrelationIDToIncomingContext(context.Request.Context(), correlationID)
	context.Request = context.Request.WithContext(ctxWithCorrelationID)

	contextWithCorrelationID, err := TransferCorrelationIDToOutgoingContext(context.Request.Context())
	if err != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.Request = context.Request.WithContext(contextWithCorrelationID)

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
func GetLoggerFromContext(ctx context.Context) (Loggerer, error) {
	if logger, ok := ctx.Value(LoggerKey).(Loggerer); ok {
		return logger, nil
	}
	return nil, errors.New("Logger not found in context")
}

// TransferCorrelationIDToOutgoingContext transfers the correlation ID from the incoming context to the outgoing context
func TransferCorrelationIDToOutgoingContext(ctx context.Context) (context.Context, error) {
	correlationID, err := GetCorrelationIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error getting correlation ID from context: %v", err)
	}
	newOutgoingCtx := AddCorrelationIDToOutgoingContext(ctx, *correlationID)
	return newOutgoingCtx, nil
}

// AddCorrelationIDToOutgoingContext adds the correlation ID to the context
func AddCorrelationIDToOutgoingContext(ctx context.Context, correlationID string) context.Context {
	existingMD, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		existingMD = metadata.New(map[string]string{})
	}
	newMD := metadata.New(map[string]string{
		CorrelationIDKey: correlationID,
	})
	mergedMD := metadata.Join(existingMD, newMD)
	return metadata.NewOutgoingContext(ctx, mergedMD)
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
