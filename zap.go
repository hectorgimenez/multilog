package multilog

import (
	"fmt"

	"go.uber.org/zap"
)

type zapHandler struct {
	logger *zap.Logger
	level  LogLevel
}

func WithZap(logLevel LogLevel, config zap.Config) Options {
	return func(logger *Logger) error {
		z, err := config.Build()
		if err != nil {
			return fmt.Errorf("error initializing Zap logger: %w", err)
		}
		logger.handlers = append(logger.handlers, zapHandler{
			logger: z,
			level:  logLevel,
		})

		return nil
	}
}

func (z zapHandler) Flush() error {
	return z.logger.Sync()
}

func (z zapHandler) Debug(msg string) {
	z.logger.Debug(msg)
}

func (z zapHandler) Info(msg string) {
	z.logger.Info(msg)
}

func (z zapHandler) Error(err error) {
	z.logger.Error(err.Error(), zap.Error(err))
}

func (z zapHandler) Level() LogLevel {
	return z.level
}
