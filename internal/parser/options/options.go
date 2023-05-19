package options

import (
	"github.com/tfadeyi/sloth-simple-comments/internal/logging"
	"github.com/tfadeyi/sloth-simple-comments/internal/parser/strategy"
)

type (
	Options struct {
		Strategy     strategy.ParsingStrategy
		IncludedDirs []string
		Logger       *logging.Logger
	}
	Option func(p *Options)
)


func Include(dirs ...string) Option {
	return func(e *Options) {
		e.IncludedDirs = dirs
	}
}

func Logger(logger *logging.Logger) Option {
	return func(e *Options) {
		e.Logger = logger
	}
}

func ParserStrategy(p strategy.ParsingStrategy) Option {
	return func(e *Options) {
		e.Strategy = p
	}
}