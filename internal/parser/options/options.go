package options

import (
	"io"

	"github.com/slosive/sloscribe/internal/logging"
	"github.com/slosive/sloscribe/internal/parser/lang"
	"github.com/slosive/sloscribe/internal/parser/specification"
)

type (
	// Options is a struct contains all the configurations available for the parser
	Options struct {
		// TargetSpecification is the specification targeted by the parser, i.e: sloth.
		// the parser will parse the source code for the target annotations.
		// Option: func Specification(target specification.Target) Option
		TargetSpecification specification.Target

		// TargetLanguage is the language targeted by the parser, i.e: go.
		// Option: func Language(lang lang.Target) Option
		TargetLanguage lang.Target

		// IncludedDirs is the array containing all the directories that will be parsed by the parser.
		// SourceFile and SourceContent will override this, if present.
		// Option: func Include(dirs ...string) Option
		IncludedDirs []string

		// Logger is the parser's logger
		// Option: func Logger(logger *logging.Logger) Option
		Logger *logging.Logger

		// SourceFile is the file the parser will parse. Shouldn't be used together with SourceContent
		// Option: func SourceFile(file string) Option
		SourceFile string

		// SourceContent is the io.Reader the parser will parse. Shouldn't be used together with SourceFile
		// Option: func SourceContent(content io.ReadCloser) Option
		SourceContent io.ReadCloser
	}
	// Option is a more atomic to configure the different Options rather than passing the entire Options struct.
	Option func(p *Options)
)

// Include configure the parser to parse the given included directories
// SourceFile and SourceContent will override this, if present.
func Include(dirs ...string) Option {
	return func(e *Options) {
		e.IncludedDirs = dirs
	}
}

// Logger configure the parser's logger
func Logger(logger *logging.Logger) Option {
	return func(e *Options) {
		log := logger.WithName("parser")
		e.Logger = &log
	}
}

// SourceFile configure the parser to parse a specific file
// Shouldn't be used together with SourceContent
func SourceFile(file string) Option {
	return func(e *Options) {
		e.SourceFile = file
	}
}

// SourceContent configure the parser to parse a specific io.Reader
// Shouldn't be used together with SourceFile
func SourceContent(content io.ReadCloser) Option {
	return func(e *Options) {
		e.SourceContent = content
	}
}

// Language configure the parser to parse using a specific target language
func Language(lang lang.Target) Option {
	return func(o *Options) {
		o.TargetLanguage = lang
	}
}

// Specification configure the parser to parse for a specific target specification
func Specification(target specification.Target) Option {
	return func(o *Options) {
		o.TargetSpecification = target
	}
}
