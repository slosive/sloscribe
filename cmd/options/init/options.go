package init

import (
	"os"

	multierr "github.com/hashicorp/go-multierror"
	"github.com/juju/errors"
	"github.com/slosive/sloscribe/cmd/options/common"
	"github.com/slosive/sloscribe/internal/generate"
	"github.com/slosive/sloscribe/internal/parser/lang"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type (
	// Options is the list of options/flag available to the application,
	// plus the clients needed by the application to function.
	Options struct {
		Formats        []string
		IncludedDirs   []string
		Source         string
		SourceLanguage lang.Target
		Specification  string
		ToFile         bool
		Services       []string
		Target         string
		*common.Options
	}
)

// New creates a new instance of the application's options
func New(c *common.Options) *Options {
	opts := new(Options)
	opts.Options = c
	return opts
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
			err = multierr.Append(err, errors.Errorf("invalid format %q was passed to --format flag", format))
		}
	}

	// @aloe code unsupported_language
	// @aloe title Unsupported TargetLanguage Error
	// @aloe summary The language passed to the --lang flag is not supported.
	// @aloe details The source language passed to the --lang flag is not currently supported by the tool.
	// The following are the supported languages: go, wasm(experimental).
	if ok := lang.IsSupportedLanguage(o.SourceLanguage); !ok {
		err = multierr.Append(err, errors.Errorf("unsupported language %q was passed to --lang flag", o.SourceLanguage))
	}
	return err
}

func getWorkingDirOrDie() string {
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
		[]string{getWorkingDirOrDie()},
		"Comma separated list of directories to be recursively parsed by the tool",
	)
	fs.StringSliceVar(
		&o.Formats,
		"format",
		[]string{"yaml"},
		"Format of the output returned by the tool. Available: yaml, json.",
	)
	fs.StringVar(
		(*string)(&o.SourceLanguage),
		"lang",
		"go",
		"Target source code language. Available: go.",
	)
	fs.StringVarP(
		&o.Source,
		"file",
		"f",
		"",
		"Source code file to parse for annotations. Example: ./metrics.go",
	)
	fs.BoolVar(
		&o.ToFile,
		"to-file",
		false,
		"Tells the tool to save the generated specifications to file, under ./slo_definitions.",
	)
	fs.StringSliceVar(
		&o.Services,
		"service-selector",
		[]string{},
		"Comma separated list of service specification names. These will select the output service specifications returned by the tool. Example: --service-selector app1,app3 ",
	)
	fs.StringVar(
		&o.Target,
		"specification",
		"sloth",
		"The SLO specification the tool should parse the source file for. Available: sloth, sloth-k8s.",
	)
}
