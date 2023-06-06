package language

import (
	"context"
	sloth "github.com/slok/sloth/pkg/prometheus/api/v1"
)

type (
	// Language is the parsing strategy used by the Parser to parse comments in the different source files
	Language interface {
		Parse(ctx context.Context) (*sloth.Spec, error)
	}
)
