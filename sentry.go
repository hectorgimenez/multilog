package multilog

import (
	"fmt"
	"reflect"
	"time"

	"github.com/getsentry/sentry-go"
)

type sentryHandler struct {
	client *sentry.Client
	level  LogLevel
}

func WithSentry(logLevel LogLevel, options sentry.ClientOptions) Options {
	return func(logger *Logger) error {
		client, err := sentry.NewClient(options)
		if err != nil {
			return fmt.Errorf("error initializing Sentry: %w", err)
		}
		logger.handlers = append(logger.handlers, sentryHandler{
			client: client,
			level:  logLevel,
		})

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
	if len(trace.Frames) > 3 {
		trace.Frames = trace.Frames[:len(trace.Frames)-3]
	}

	return trace
}

func (s sentryHandler) Level() LogLevel {
	return s.level
}

func (s sentryHandler) Flush() error {
	s.client.Flush(time.Second * 5)
	return nil
}

func (s sentryHandler) Debug(msg string) {
	event := sentryEventWithMessageAndLevel(msg, sentry.LevelDebug)
	event.Threads = []sentry.Thread{{
		Stacktrace: getSentryStackTrace(nil),
		Current:    true,
	}}
	s.client.CaptureEvent(event, nil, sentry.CurrentHub().Scope())
}

func (s sentryHandler) Info(msg string) {
	event := sentryEventWithMessageAndLevel(msg, sentry.LevelInfo)
	event.Threads = []sentry.Thread{{
		Stacktrace: getSentryStackTrace(nil),
		Current:    true,
	}}
	s.client.CaptureEvent(event, nil, sentry.CurrentHub().Scope())
}

func (s sentryHandler) Error(err error) {
	event := sentryEventWithMessageAndLevel(err.Error(), sentry.LevelError)
	event.Exception = []sentry.Exception{{
		Value:      err.Error(),
		Type:       reflect.TypeOf(err).String(),
		Stacktrace: getSentryStackTrace(err),
	}}
	s.client.CaptureEvent(event, nil, sentry.CurrentHub().Scope())
}
