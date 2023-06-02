package golang

import (
	"github.com/tfadeyi/slotalk/internal/parser/options"
)

func Parser() options.Option {
	return func(e *options.Options) {
		e.Strategy = newParser(e.Logger, e.SourceFile, e.SourceContent, e.IncludedDirs...)
	}
}
