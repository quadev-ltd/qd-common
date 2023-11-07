package log

import (
	"github.com/rs/zerolog"
)

// Loggerer is the interface of the logger
type Loggerer interface {
	Error(err error, message string)
	Info(message string)
	Warn(message string)
}

// Logger is the logger of the application
type Logger struct {
	log zerolog.Logger
}

// Error logs an error
func (logger *Logger) Error(err error, message string) {
	logger.log.Error().Err(err).Msg(message)
}

// Info logs an info
func (logger *Logger) Info(message string) {
	logger.log.Info().Msg(message)
}

// Warn logs a warning
func (logger *Logger) Warn(message string) {
	logger.log.Warn().Msg(message)
}
