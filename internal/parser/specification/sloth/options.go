package sloth

import (
	"github.com/tfadeyi/slotalk/internal/parser/options"
	"github.com/tfadeyi/slotalk/internal/parser/specification/sloth/language/golang"
)

// Parser returns the options.Options struct containing
func Parser() options.Option {
	return func(opts *options.Options) {
		opts.Specification = newParser(Options{
			Language: opts.Language,
			GolangOpts: golang.Options{
				Logger:           opts.Logger,
				SourceFile:       opts.SourceFile,
				SourceContent:    opts.SourceContent,
				InputDirectories: opts.IncludedDirs,
			},
		})
	}
}
