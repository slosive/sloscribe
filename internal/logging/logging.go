package logging

import (
	"context"
	"log"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
)

type (
	Logger struct {
		logr.Logger
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
	return Logger{l}
}

// NewStandardLogger wraps the creation of a Golang standard logger
func NewStandardLogger() Logger {
	return Logger{stdr.New(log.Default())}
}
