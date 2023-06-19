package parser

import (
	"context"

	"github.com/juju/errors"
	"github.com/tfadeyi/slosive/internal/parser/options"
)

type (
	// Parser parses source files containing the sloth definitions
	Parser struct {
		// Opts contains the different options available to the parser.
		// These are applied by the parser constructor during initialization
		Opts *options.Options
	}
)

// New creates a new instance of the parser. See options.Option for more info on the available configuration.
func New(opts ...options.Option) (*Parser, error) {
	defaultOpts := new(options.Options)
	for _, opt := range opts {
		opt(defaultOpts)
	}
	for _, opt := range opts {
		opt(defaultOpts)
	}

	if defaultOpts.TargetSpecification == nil {
		return nil, errors.New("parsing strategy configuration was not set")
	}

	return &Parser{defaultOpts}, nil
}

// Parse parses the data source for the target annotations using the given parser configurations and returns a parsed specification.
func (p *Parser) Parse(ctx context.Context) (map[string]any, error) {
	return p.Opts.TargetSpecification.Parse(ctx)
}
