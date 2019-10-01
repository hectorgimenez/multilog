package multilog

import (
	"errors"
	"fmt"
	"log"
)

const (
	Debug LogLevel = 1
	Info  LogLevel = 2
	Error LogLevel = 3
	Fatal LogLevel = 4
	Panic LogLevel = 5
)

type Handler interface {
	Flush() error
	Debug(msg string)
	Info(msg string)
	Error(err error)
	Level() LogLevel
}

type Logger struct {
	handlers []Handler
}

type Options func(*Logger) error
type LogLevel int8

func NewLogger(options ...Options) (*Logger, error) {
	logger := &Logger{}
	for _, opt := range options {
		if err := opt(logger); err != nil {
			return nil, err
		}
	}

	return logger, nil
}

func WithHandler(handler Handler) Options {
	return func(logger *Logger) error {
		if handler == nil {
			return errors.New("handler must be specified")
		}
		logger.handlers = append(logger.handlers, handler)

		return nil
	}
}

func (l *Logger) Flush() (err error) {
	for _, handler := range l.handlers {
		e := handler.Flush()
		if e != nil {
			err = fmt.Errorf("error flushing: %w", e)
		}
	}

	return err
}

func (l *Logger) Debug(msg string) {
	if len(l.handlers) == 0 {
		log.Print(msg)
	}

	for _, handler := range l.handlers {
		if handler.Level() > Debug {
			continue
		}
		handler.Debug(msg)
	}
}

func (l *Logger) Info(msg string) {
	if len(l.handlers) == 0 {
		log.Print(msg)
	}

	for _, handler := range l.handlers {
		if handler.Level() > Info {
			continue
		}
		handler.Info(msg)
	}
}

func (l *Logger) Error(err error) {
	if len(l.handlers) == 0 {
		log.Print(err)
	}

	for _, handler := range l.handlers {
		if handler.Level() > Error {
			continue
		}
		handler.Error(err)
	}
}
