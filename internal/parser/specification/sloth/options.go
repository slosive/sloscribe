package sloth

import (
	"github.com/slosive/sloscribe/internal/parser/options"
	"github.com/slosive/sloscribe/internal/parser/specification/sloth/language/golang"
)

// Parser returns the options.Option to run the parser targeting sloth as a specification
func Parser(kubernetes bool) options.Option {
	return func(opts *options.Options) {
		opts.TargetSpecification = newParser(Options{
			Language: opts.TargetLanguage,
			GolangOpts: golang.Options{
				Logger:           opts.Logger,
				SourceFile:       opts.SourceFile,
				SourceContent:    opts.SourceContent,
				InputDirectories: opts.IncludedDirs,
				Kubernetes:       kubernetes,
			},
		})
	}
}
