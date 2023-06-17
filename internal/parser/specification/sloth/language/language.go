package language

import (
	"context"
)

type (
	// Language is the parsing strategy used by the Parser to parse comments in the different source files
	Language interface {
		Parse(ctx context.Context) (map[string]any, error)
	}
)
