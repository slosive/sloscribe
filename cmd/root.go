/*
Copyright © 2023 Oluwole Fadeyi
*/
package cmd

import (
	"context"
	"os"

	commonoptions "github.com/slosive/xslosive/cmd/options/common"
	"github.com/slosive/xslosive/internal/logging"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd *cobra.Command

func cmd(opts *commonoptions.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "xslosive",
		Short: "Generate Sloth SLO/SLI definitions from code annotations.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.LoggerFromContext(cmd.Context())
			logger = logger.WithName("root")

			if err := opts.Complete(); err != nil {
				logger.Error(err, "flag argument error")
				return err
			}
			if opts.LogLevel != "" {
				logger = logger.SetLevel(opts.LogLevel)
			}
			cmd.SetContext(logging.ContextWithLogger(cmd.Context(), logger))
			return nil
		},
	}
	opts = opts.Prepare(cmd)
	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx context.Context) {
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}
func init() {
	opts := commonoptions.New()
	rootCmd = cmd(opts)
	rootCmd.AddCommand(specInitCmd(opts))
	rootCmd.AddCommand(versionCmd)
}
