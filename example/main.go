package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"

	"github.com/hectorgimenez/multilog"
)

func main() {
	logger, err := multilog.NewLogger(
		multilog.WithZap(zap.NewDevelopmentConfig()),
		multilog.WithSentry(sentry.ClientOptions{
			Dsn: "", // Enter Sentry DSN here
		}),
	)
	if err != nil {
		log.Panic("error creating logger", err)
	}
	defer func() {
		if err := logger.Flush(); err != nil {
			fmt.Printf("error flushing logger: %v", err)
		}
	}()

	testLogger(logger)
}

func testLogger(logger *multilog.Logger) {
	logger.Info("This a simple INFO")
	logger.Debug("We are logging at Debug level")
	logger.Error(errors.New("we can log also an error"))
}
