package log

import (
	"context"
	"os"

	"github.com/rs/zerolog"

	"github.com/gustavo-m-franco/qd-common/pkg/config"
)

// CorrelationIDKey is the key of the correlation ID in the metadata
const CorrelationIDKey = "correlation_id"

// Factoryer is the interface for creating a log factory to create a logger
type Factoryer interface {
	NewLogger() Loggerer
	NewLoggerWithCorrelationID(ctx context.Context) (Loggerer, error)
}

// Factory is the factory for creating a logger
type Factory struct {
	environment string
}

// NewLogFactory creates a new log factory
func NewLogFactory(environment string) Factoryer {
	return &Factory{
		environment: environment,
	}
}

func setUpLog(log zerolog.Logger, environment string) zerolog.Logger {
	if environment == config.ProductionEnvironment {
		log = log.Level(zerolog.WarnLevel)
	} else {
		log = log.Level(zerolog.DebugLevel)
	}
	return log
}

// NewLogger creates a new logger for the given environment
func (logFactory *Factory) NewLogger() Loggerer {
	var log zerolog.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	log = setUpLog(log, logFactory.environment)
	return &Logger{
		log: log,
	}
}

// NewLoggerWithCorrelationID creates a new logger with the correlation ID
func (logFactory *Factory) NewLoggerWithCorrelationID(ctx context.Context) (Loggerer, error) {
	var log zerolog.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	correlationID, error := GetCorrelationIDFromContext(ctx)
	if error != nil {
		return nil, error
	}
	log = setUpLog(
		log,
		logFactory.environment,
	).With().Str(CorrelationIDKey, *correlationID).Logger()
	return &Logger{
		log: log,
	}, nil
}
