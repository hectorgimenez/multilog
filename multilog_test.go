package multilog

import (
	"errors"
	"testing"
)

type fakeHandler struct {
	flushCalled bool
	debugCalled bool
	infoCalled  bool
	errorCalled bool
	customLevel LogLevel
}

func (f *fakeHandler) Flush() error {
	f.flushCalled = true
	return nil
}

func (f *fakeHandler) Debug(msg string) {
	f.debugCalled = true
}

func (f *fakeHandler) Info(msg string) {
	f.infoCalled = true
}

func (f *fakeHandler) Error(err error) {
	f.errorCalled = true
}

func (f *fakeHandler) Level() LogLevel {
	return f.customLevel
}

func TestNewLogger(t *testing.T) {
	t.Run("Create new Logger without any handler", func(t *testing.T) {
		logger, err := NewLogger()
		if err != nil {
			t.Errorf("Error not expected: %v", err)
		}
		if len(logger.handlers) != 0 {
			t.Errorf("Logger should not contain any handler, found: %d", len(logger.handlers))
		}
	})
}

func TestLogger_Log(t *testing.T) {
	tests := []struct {
		name        string
		handler     *fakeHandler
		debugCalled bool
		infoCalled  bool
		errorCalled bool
	}{
		{
			name: "Given a handler with Error level should not log anything less than error",
			handler: &fakeHandler{
				customLevel: Error,
			},
			debugCalled: false,
			infoCalled:  false,
			errorCalled: true,
		},
		{
			name: "Given a handler with Info level should log info and error",
			handler: &fakeHandler{
				customLevel: Info,
			},
			debugCalled: false,
			infoCalled:  true,
			errorCalled: true,
		},
		{
			name: "Given a handler with Info level should log everything",
			handler: &fakeHandler{
				customLevel: Debug,
			},
			debugCalled: true,
			infoCalled:  true,
			errorCalled: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				handlers: []Handler{tt.handler},
			}
			l.Debug("Debug")
			l.Info("Info")
			l.Error(errors.New("error"))

			if tt.handler.debugCalled != tt.debugCalled {
				t.Errorf("Debug is %v expected %v", tt.handler.debugCalled, tt.debugCalled)
			}
			if tt.handler.infoCalled != tt.infoCalled {
				t.Errorf("Info is %v expected %v", tt.handler.infoCalled, tt.infoCalled)
			}
			if tt.handler.errorCalled != tt.errorCalled {
				t.Errorf("Error is %v expected %v", tt.handler.errorCalled, tt.errorCalled)
			}
		})
	}
}
