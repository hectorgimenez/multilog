package multilog

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
)

func WithSentry(options sentry.ClientOptions) Options {
	return func(logger *Logger) error {
		client, err := sentry.NewClient(options)
		if err != nil {
			return fmt.Errorf("error initializing Sentry: %w", err)
		}
		logger.sentryLogger = client

		return nil
	}
}

func sentryEventWithMessageAndLevel(msg string, level sentry.Level) *sentry.Event {
	event := sentry.NewEvent()
	event.Message = msg
	event.Level = level
	event.Timestamp = time.Now().Unix()

	return event
}

func getSentryStackTrace(err error) *sentry.Stacktrace {
	if err != nil {
		if trace := sentry.ExtractStacktrace(err); trace != nil {
			return trace
		}
	}

	trace := sentry.NewStacktrace()

	// Remove last two frames in order to find real stacktrace position
	if len(trace.Frames) > 0 {
		trace.Frames = trace.Frames[:len(trace.Frames)-2]
	}

	return trace
}
