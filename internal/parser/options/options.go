package options

import (
	"github.com/tfadeyi/sloth-simple-comments/internal/logging"
	"github.com/tfadeyi/sloth-simple-comments/internal/parser/strategy"
	"io"
)

type (
	Options struct {
		Strategy      strategy.ParsingStrategy
		IncludedDirs  []string
		Logger        *logging.Logger
		SourceFile    string
		SourceContent io.ReadCloser
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

func SourceFile(file string) Option {
	return func(e *Options) {
		e.SourceFile = file
	}
}

func SourceContent(content io.ReadCloser) Option {
	return func(e *Options) {
		e.SourceContent = content
	}
}
