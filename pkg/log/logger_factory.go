package log

import (
	"context"
	"os"

	"github.com/gustavo-m-franco/qd-common/pkg/config"

	"github.com/rs/zerolog"
)

// CorrelationIDKey is the key of the correlation ID in the metadata
const CorrelationIDKey = "correlation_id"

// LogFactoryer	is the interface for creating a logger
type LogFactoryer interface {
	NewLogger() Loggerer
	NewLoggerWithCorrelationID(ctx context.Context) (Loggerer, error)
}

// LogFactory is the factory for creating a logger
type LogFactory struct {
	environment string
}

// NewLogFactory creates a new log factory
func NewLogFactory(environment string) LogFactoryer {
	return &LogFactory{
		environment: environment,
	}
}

func setUpLevel(log zerolog.Logger, environment string) zerolog.Logger {
	if environment == config.ProductionEnvironment {
		log = log.Level(zerolog.WarnLevel)
	} else {
		log = log.Level(zerolog.DebugLevel)
	}
	return log
}

// NewLogger creates a new logger for the given environment
func (logFactory *LogFactory) NewLogger() Loggerer {
	var log zerolog.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	log = setUpLevel(log, logFactory.environment)
	return &Logger{
		log: log,
	}
}

// NewLoggerWithCorrelationID creates a new logger with the correlation ID
func (logFactory *LogFactory) NewLoggerWithCorrelationID(ctx context.Context) (Loggerer, error) {
	var log zerolog.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	correlationID, error := GetCorrelationIDFromContext(ctx)
	if error != nil {
		return nil, error
	}
	log = setUpLevel(
		log,
		logFactory.environment,
	).With().Str(CorrelationIDKey, *correlationID).Logger()
	return &Logger{
		log: log,
	}, nil
}
