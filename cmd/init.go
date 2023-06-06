package cmd

import (
	"io"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
	commonoptions "github.com/tfadeyi/slotalk/cmd/options/common"
	initoptions "github.com/tfadeyi/slotalk/cmd/options/init"
	"github.com/tfadeyi/slotalk/internal/generate"
	"github.com/tfadeyi/slotalk/internal/logging"
	"github.com/tfadeyi/slotalk/internal/parser"
	"github.com/tfadeyi/slotalk/internal/parser/lang"
	"github.com/tfadeyi/slotalk/internal/parser/options"
	"github.com/tfadeyi/slotalk/internal/parser/specification/sloth"
)

func specInitCmd(common *commonoptions.Options) *cobra.Command {
	opts := initoptions.New(common)
	var inputReader io.ReadCloser
	var targetLanguage options.Option
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Init generates the Sloth definition specification from source code comments.",
		Long: `The init command parses files in the target directory for comments using the @sloth tags,
i.e: 
	// @sloth.slo name chatgpt
	// @sloth.slo objective 95.0


These are then used to generate Sloth definition specifications. 
i.e:
	version: prometheus/v1
	service: "chatgpt"
	slos:
		- name: chat-gpt-availability
		  objective: 95
`,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.LoggerFromContext(cmd.Context())
			logger = logger.WithName("init")

			if err := opts.Complete(); err != nil {
				logger.Error(err, "flag argument error")
				return err
			}

			switch opts.SrcLanguage {
			case lang.Rust:
				err := errors.New("The rust parser has not been fully implemented and shouldn't be used! It will have unexpected behaviours.")
				logger.Error(err, "")
				return err
				targetLanguage = options.Language(lang.Rust)
			default:
				targetLanguage = options.Language(lang.Go)
			}

			if opts.Source == "-" {
				inputReader = io.NopCloser(cmd.InOrStdin())
			}

			cmd.SetContext(logging.ContextWithLogger(cmd.Context(), logger))
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.LoggerFromContext(cmd.Context())
			logger.Info("Parsing source code for slo definitions",
				"directories", opts.IncludedDirs,
				"source", opts.Source,
			)

			parser, err := parser.New(
				targetLanguage,
				sloth.Parser(),
				options.Logger(&logger),
				options.SourceFile(opts.Source),
				options.SourceContent(inputReader),
				options.Include(opts.IncludedDirs...))
			if err != nil {
				logger.Error(err, "parser initialization")
				return err
			}

			service, err := parser.Parse(cmd.Context())
			if err != nil {
				logger.Error(err, "parser parsing error")
				return err
			}

			logger.Info("Source code was parsed!")
			logger.Info("Printing result specification to stdout.")
			return generate.WriteSpecification(service, true, "", opts.Formats...)
		},
	}
	opts = opts.Prepare(cmd)
	return cmd
}
