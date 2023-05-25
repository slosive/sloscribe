package cmd

import (
	"github.com/spf13/cobra"
	specoptions "github.com/tfadeyi/sloth-simple-comments/cmd/options/spec"
	"github.com/tfadeyi/sloth-simple-comments/internal/generate"
	"github.com/tfadeyi/sloth-simple-comments/internal/logging"
	"github.com/tfadeyi/sloth-simple-comments/internal/parser"
	"github.com/tfadeyi/sloth-simple-comments/internal/parser/lang"
	"github.com/tfadeyi/sloth-simple-comments/internal/parser/options"
	"github.com/tfadeyi/sloth-simple-comments/internal/parser/strategy/golang"
	"github.com/tfadeyi/sloth-simple-comments/internal/parser/strategy/wasm"
	"io"
)

func specGenerateCmd() *cobra.Command {
	opts := specoptions.New()
	var inputReader io.ReadCloser
	var languageParser options.Option
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
			if err := opts.Complete(); err != nil {
				return err
			}

			switch opts.SrcLanguage {
			case lang.Wasm:
				logger.Info("The wasm parser has not been fully implemented and shouldn't be used! It will have unexpected behaviours.")
				languageParser = wasm.Parser()
			default:
				languageParser = golang.Parser()
			}

			if opts.Source == "-" {
				inputReader = io.NopCloser(cmd.InOrStdin())
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.LoggerFromContext(cmd.Context())

			logger.Info("Parsing source code for slo definitions",
				"directories", opts.IncludedDirs,
				"source", opts.Source,
			)

			parser, err := parser.New(
				languageParser,
				options.Logger(&logger),
				options.SourceFile(opts.Source),
				options.SourceContent(inputReader),
				options.Include(opts.IncludedDirs...))
			if err != nil {
				return err
			}
			service, err := parser.Parse(cmd.Context())
			if err != nil {
				return err
			}

			logger.Info("Source code was parsed successfully")

			return generate.WriteSpecification(service, true, "", opts.Formats...)
		},
	}
	opts = opts.Prepare(cmd)
	return cmd
}

func init() {
	rootCmd.AddCommand(specGenerateCmd())
}
