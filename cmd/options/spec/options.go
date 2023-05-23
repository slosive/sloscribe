package spec

import (
	multierr "github.com/hashicorp/go-multierror"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	goaloe "github.com/tfadeyi/go-aloe"
	"github.com/tfadeyi/sloth-simple-comments/internal/generate"
	"github.com/tfadeyi/sloth-simple-comments/internal/parser/lang"
	"os"
)

type (
	// Options is the list of options/flag available to the application,
	// plus the clients needed by the application to function.
	Options struct {
		StdOutput    bool
		Formats      []string
		IncludedDirs []string
		SrcLanguage  lang.SourceLanguage
	}
)

// New creates a new instance of the application's options
func New() *Options {
	return new(Options)
}

// Prepare assigns the applications flag/options to the cobra cli
func (o *Options) Prepare(cmd *cobra.Command) *Options {
	o.addAppFlags(cmd.Flags())
	return o
}

// Complete initialises the components needed for the application to function given the options
func (o *Options) Complete() error {
	var err error
	// @aloe code unsupported_output_format
	// @aloe title Unsupported Output Format Error
	// @aloe summary The format passed to the --format flag is not supported.
	// @aloe details The format passed to the --lang flag is not currently supported by the tool.
	// The following are the supported languages: yaml(default), json.
	for _, format := range o.Formats {
		if ok := generate.IsValidOutputFormat(format); !ok {
			err = goaloe.Default().Error(
				multierr.Append(err, errors.Errorf("invalid format %q was passed to --format flag", format)),
				"unsupported_output_format")
		}
	}

	// @aloe code unsupported_language
	// @aloe title Unsupported Language Error
	// @aloe summary The language passed to the --lang flag is not supported.
	// @aloe details The source language passed to the --lang flag is not currently supported by the tool.
	// The following are the supported languages: go, wasm(experimental).
	if ok := lang.IsSupportedLanguage(o.SrcLanguage); !ok {
		err = goaloe.Default().Error(
			multierr.Append(err, errors.Errorf("unsupported language %q was passed to --lang flag", o.SrcLanguage)),
			"unsupported_language")
	}

	return err
}

func GetWorkingDirOrDie() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

func (o *Options) addAppFlags(fs *pflag.FlagSet) {
	fs.StringSliceVar(
		&o.IncludedDirs,
		"dirs",
		[]string{GetWorkingDirOrDie()},
		"Comma separated list of directories to be parses by the tool",
	)
	fs.StringSliceVar(
		&o.Formats,
		"format",
		[]string{"yaml"},
		"Output format (yaml,json)",
	)
	fs.BoolVar(
		&o.StdOutput,
		"stdout",
		false,
		"Print output to standard output.",
	)
	fs.StringVar(
		(*string)(&o.SrcLanguage),
		"lang",
		"go",
		"Language of the source files. (go, wasm)",
	)
}
