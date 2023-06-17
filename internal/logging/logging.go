package logging

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
)

type (
	Logger struct {
		logr.Logger
		*sync.Mutex
	}
)

// ContextWithLogger wraps the logr NewContext function
func ContextWithLogger(ctx context.Context, l Logger) context.Context {
	return logr.NewContext(ctx, l.Logger)
}

// LoggerFromContext wraps the LoggerFromContext or creates a Zap production logger
func LoggerFromContext(ctx context.Context) Logger {
	l, err := logr.FromContext(ctx)
	if err != nil {
		return NewStandardLogger()
	}
	return Logger{l, new(sync.Mutex)}
}

// NewStandardLogger wraps the creation of a Golang standard logger
func NewStandardLogger() Logger {
	return Logger{stdr.New(log.Default()), new(sync.Mutex)}
}

func (l *Logger) WithName(name string) Logger {
	return Logger{l.Logger.WithName(name), new(sync.Mutex)}
}

// SetLevel sets the global level against which all info logs will be compared.
// If this is greater than or equal to the "V" of the logger, the message will be logged.
// A higher value here means more logs will be written
func (l *Logger) SetLevel(lvl string) Logger {
	l.Lock()
	defer l.Unlock()
	stdr.SetVerbosity(int(findLogLevel(lvl)))
	return *l
}

// Info logs the message using the info log level (2)
func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.V(int(stdr.Info)).Info(msg, keysAndValues...)
}

// Warn logs the error using the warning log level (3)
func (l *Logger) Warn(err error, keysAndValues ...interface{}) {
	l.V(int(stdr.Error)).Info(err.Error(), keysAndValues...)
}

// Debug logs the message using the info debug level (1)
func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	l.V(5).Info(msg, keysAndValues...)
}

func findLogLevel(lvl string) stdr.MessageClass {
	switch strings.TrimSpace(strings.ToLower(lvl)) {
	case "info":
		return stdr.Info
	case "debug":
		return stdr.MessageClass(5)
	case "warn":
		return stdr.Error
	case "none":
		return stdr.None
	}
	return stdr.Info
}

// IsValidLevel returns true if the input level is one of the valid values ("info", "debug", "warn", "none")
func IsValidLevel(lvl string) bool {
	switch strings.TrimSpace(strings.ToLower(lvl)) {
	case "info", "debug", "warn", "none":
		return true
	}
	return false
}
