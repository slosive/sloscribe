package golang

import (
	"github.com/tfadeyi/sloth-simple-comments/internal/parser/options"
)

func Parser() options.Option {
	return func(e *options.Options) {
		e.Strategy = newParser(e.Logger, e.IncludedDirs...)
	}
}
