package common

import (
	multierr "github.com/hashicorp/go-multierror"
	"github.com/juju/errors"
	"github.com/slosive/sloscribe/internal/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type (
	// Options is the list of options/flag available to the application,
	// plus the clients needed by the application to function.
	Options struct {
		// LogLevel used by the exporter's logger (none, debug, info, warn)
		LogLevel string
	}
)

// New creates a new instance of the application's options
func New() *Options {
	return new(Options)
}

// Prepare assigns the applications flag/options to the cobra cli
func (o *Options) Prepare(cmd *cobra.Command) *Options {
	o.addAppFlags(cmd.PersistentFlags())
	return o
}

// Complete initialises the components needed for the application to function given the options
func (o *Options) Complete() error {
	var err error
	if !logging.IsValidLevel(o.LogLevel) {
		// @aloe code invalid_log_level
		// @aloe title Invalid Log-Level Argument
		// @aloe summary The log level passed to the --log-level flag is not supported.
		// @aloe details The log level passed to the --log-level flag is not currently supported by the tool.
		// The following are supported: none, debug, info(default), warn.
		err = multierr.Append(err, errors.Errorf("invalid log-level %q was passed to --log-level flag", o.LogLevel))
	}
	return err
}

func (o *Options) addAppFlags(fs *pflag.FlagSet) {
	fs.StringVar(
		&o.LogLevel,
		"log-level",
		"info",
		"Only log messages with the given severity or above. One of: [none, debug, info, warn], errors will always be printed",
	)
}
