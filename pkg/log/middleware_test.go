package log

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"

	"github.com/quadev-ltd/qd-common/pkg/log/mock"
)

func TestGetLoggerFromContext(t *testing.T) {
	t.Run("Context_With_Logger", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		mockLogger := mock.NewMockLoggerer(controller)
		ctx := context.WithValue(context.Background(), LoggerKey, mockLogger)

		logger, err := GetLoggerFromContext(ctx)
		assert.NoError(t, err)
		assert.Equal(t, mockLogger, logger, "Expected to get the mock logger from context")
	})

	t.Run("Context_Without_Logger", func(t *testing.T) {
		ctx := context.Background()

		logger, err := GetLoggerFromContext(ctx)
		assert.Error(t, err)
		assert.Equal(t, "Logger not found in context", err.Error())
		assert.Nil(t, logger, "Expected to get nil when the logger is not present in context")
	})

	t.Run("Context_With_Non_Loggerer_Value_For_LoggerKey", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), LoggerKey, "not-a-logger")

		logger, err := GetLoggerFromContext(ctx)
		assert.Error(t, err)
		assert.Equal(t, "Logger not found in context", err.Error())
		assert.Nil(t, logger, "Expected to get nil when the logger is not present in context")
	})
}

func TestAddCorrelationIDToContext(t *testing.T) {
	t.Run("AddCorrelationIDToContext_Add_Correlation_ID", func(t *testing.T) {
		correlationID := "test-correlation-id"
		ctx := context.Background()

		// Call the function to add the correlation ID to the context
		ctxWithCorrelationID := AddCorrelationIDToOutgoingContext(ctx, correlationID)

		// Retrieve the metadata from the context
		md, ok := metadata.FromOutgoingContext(ctxWithCorrelationID)
		if !ok {
			t.Fatal("Expected metadata to be present in context")
		}

		// Assert that the correlation ID is correctly set in the metadata
		assert.Equal(t, correlationID, md[CorrelationIDKey][0], "Expected correlation ID to be set correctly in metadata")
	})
}

func TestGetCorrelationIDFromContext(t *testing.T) {
	t.Run("Context_With_Correlation_ID", func(t *testing.T) {
		correlationID := "test-correlation-id"
		md := metadata.New(map[string]string{
			CorrelationIDKey: correlationID,
		})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		retrievedID, err := GetCorrelationIDFromContext(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedID)
		assert.Equal(t, correlationID, *retrievedID, "Expected to retrieve the correct correlation ID")
	})

	t.Run("Context_Without_Metadata", func(t *testing.T) {
		ctx := context.Background()

		_, err := GetCorrelationIDFromContext(ctx)
		assert.Error(t, err)
		assert.Equal(t, "Metadata not found in context", err.Error(), "Expected an error when metadata is missing")
	})

	t.Run("Context_With_Metadata_But_Without_Correlation_ID", func(t *testing.T) {
		md := metadata.New(map[string]string{})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		_, err := GetCorrelationIDFromContext(ctx)
		assert.Error(t, err)
		assert.Equal(t, "Correlation ID not found in metadata", err.Error(), "Expected an error when correlation ID is missing")
	})

	t.Run("Context_With_Multiple_Correlation_IDs", func(t *testing.T) {
		md := metadata.Pairs(
			CorrelationIDKey, "id1",
			CorrelationIDKey, "id2",
		)
		ctx := metadata.NewIncomingContext(context.Background(), md)

		_, err := GetCorrelationIDFromContext(ctx)
		assert.Error(t, err)
		assert.Equal(t, "Correlation ID not found in metadata", err.Error(), "Expected an error when multiple correlation IDs are present")
	})
}
