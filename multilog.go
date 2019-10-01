package multilog

import (
	"log"
	"reflect"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

const (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
	Error LogLevel = "error"
	Fatal LogLevel = "fatal"
	Panic LogLevel = "panic"
)

type Logger struct {
	zapLogger    *zap.Logger
	sentryLogger *sentry.Client
}

type Options func(*Logger) error
type LogLevel string

func NewLogger(options ...Options) (*Logger, error) {
	logger := &Logger{}
	for _, opt := range options {
		if err := opt(logger); err != nil {
			return nil, err
		}
	}

	return logger, nil
}

func (l *Logger) Flush() error {
	l.sentryLogger.Flush(time.Second * 5)

	if l.zapLogger != nil {
		return l.zapLogger.Sync()
	}

	return nil
}

func (l *Logger) Debug(msg string) {
	if l.sentryLogger != nil {
		event := sentryEventWithMessageAndLevel(msg, sentry.LevelDebug)
		event.Threads = []sentry.Thread{{
			Stacktrace: getSentryStackTrace(nil),
			Current:    true,
		}}
		l.sentryLogger.CaptureEvent(event, nil, sentry.CurrentHub().Scope())
	}

	if l.zapLogger != nil {
		l.zapLogger.Debug(msg)
	} else {
		log.Printf(msg)
	}
}

func (l *Logger) Info(msg string) {
	if l.sentryLogger != nil {
		event := sentryEventWithMessageAndLevel(msg, sentry.LevelInfo)
		event.Threads = []sentry.Thread{{
			Stacktrace: getSentryStackTrace(nil),
			Current:    true,
		}}
		l.sentryLogger.CaptureEvent(event, nil, sentry.CurrentHub().Scope())
	}

	if l.zapLogger != nil {
		l.zapLogger.Info(msg)
	} else {
		log.Printf(msg)
	}
}

func (l *Logger) Error(err error) {
	if l.sentryLogger != nil {
		event := sentryEventWithMessageAndLevel(err.Error(), sentry.LevelError)
		event.Exception = []sentry.Exception{{
			Value:      err.Error(),
			Type:       reflect.TypeOf(err).String(),
			Stacktrace: getSentryStackTrace(err),
		}}
		l.sentryLogger.CaptureEvent(event, nil, sentry.CurrentHub().Scope())
	}

	if l.zapLogger != nil {
		l.zapLogger.Error(err.Error(), zap.Error(err))
	} else {
		log.Printf(err.Error())
	}
}
