package log

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"

	"github.com/gustavo-m-franco/qd-common/pkg/config"
)

func TestNewLogger(t *testing.T) {
	environments := []struct {
		name          string
		environment   string
		expectedLevel zerolog.Level
	}{
		{"Production Environment", config.ProductionEnvironment, zerolog.WarnLevel},
		{"Non-Production Environment", "development", zerolog.DebugLevel},
	}

	for _, env := range environments {
		t.Run(env.name, func(t *testing.T) {
			factory := NewLogFactory(env.environment)
			logger := factory.NewLogger().(*Logger)
			assert.Equal(t, env.expectedLevel, logger.log.GetLevel(), "Log level should match the environment setting")
		})
	}
}

func TestNewLoggerWithCorrelationID(t *testing.T) {
	t.Run("With_Correlation_ID", func(t *testing.T) {
		correlationID := "test-correlation-id"
		md := metadata.New(map[string]string{
			CorrelationIDKey: correlationID,
		})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		factory := NewLogFactory("development")
		_, err := factory.NewLoggerWithCorrelationID(ctx)
		assert.NoError(t, err)
	})

	t.Run("Without_Correlation_ID", func(t *testing.T) {
		ctx := context.Background()
		factory := NewLogFactory("development")
		_, err := factory.NewLoggerWithCorrelationID(ctx)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "Metadata not found in context")
	})
}
