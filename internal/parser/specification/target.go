package specification

import (
	"context"
)

type (
	// Target is the specification target interface, it defines the specification target contract that
	// all new targets should adhere to.
	Target interface {
		// Parse returns a specification struct given a data source, returns error if parsing fails
		Parse(ctx context.Context) (map[string]any, error)
	}
)
