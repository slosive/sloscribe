package sloth

import (
	"context"

	"github.com/tfadeyi/slotalk/internal/parser/lang"
	"github.com/tfadeyi/slotalk/internal/parser/specification/sloth/language"
	"github.com/tfadeyi/slotalk/internal/parser/specification/sloth/language/golang"
)

// Parser struct, stores the language parser used to parse the data source
type parser struct {
	languageParser language.Language
}

// Options is a struct contains all the configurations available for the sloth parser
type Options struct {
	Language   lang.Target
	GolangOpts golang.Options
}

// newParser client parser performs all checks at initialization time
func newParser(opts Options) *parser {
	var selectedLanguageParser language.Language
	switch opts.Language {
	case lang.Go:
		selectedLanguageParser = golang.NewParser(opts.GolangOpts)
	case lang.Rust:
	}
	return &parser{
		languageParser: selectedLanguageParser,
	}
}

// Parse parses the sloth specification using the target language parser
func (p parser) Parse(ctx context.Context) (any, error) {
	return p.languageParser.Parse(ctx)
}
