package options

import (
	"io"

	"github.com/tfadeyi/slotalk/internal/logging"
	"github.com/tfadeyi/slotalk/internal/parser/lang"
	"github.com/tfadeyi/slotalk/internal/parser/specification"
)

type (
	Options struct {
		Specification specification.Target
		IncludedDirs  []string
		Logger        *logging.Logger
		SourceFile    string
		SourceContent io.ReadCloser
		Language      lang.Target
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
		log := logger.WithName("parser")
		e.Logger = &log
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

func Language(lang lang.Target) Option {
	return func(o *Options) {
		o.Language = lang
	}
}

func Specification(target specification.Target) Option {
	return func(o *Options) {
		o.Specification = target
	}
}
