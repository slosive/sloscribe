package parser

import (
	"context"

	"github.com/juju/errors"
	sloth "github.com/slok/sloth/pkg/prometheus/api/v1"
	"github.com/tfadeyi/sloth-simple-comments/internal/parser/options"
)

type (
	// Parser parses source files containing the sloth definitions
	Parser struct {
		Opts *options.Options
	}
)

// New creates a new instance of the parser, defaults to golang parsing strategy if non are passed
func New(opts ...options.Option) (*Parser, error) {
	defaultOpts := new(options.Options)
	for _, opt := range opts {
		opt(defaultOpts)
	}
	for _, opt := range opts {
		opt(defaultOpts)
	}

	if defaultOpts.Strategy == nil {
		return nil, errors.New("parsing strategy configuration was not set")
	}

	return &Parser{defaultOpts}, nil
}

func (p *Parser) Parse(ctx context.Context) (*sloth.Spec, error) {
	return p.Opts.Strategy.Parse(ctx)
}
