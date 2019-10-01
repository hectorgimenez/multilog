package multilog

import (
	"fmt"

	"go.uber.org/zap"
)

func WithZap(config zap.Config) Options {
	return func(logger *Logger) error {
		z, err := config.Build()
		if err != nil {
			return fmt.Errorf("error initializing Zap logger: %w", err)
		}
		logger.zapLogger = z

		return nil
	}
}
