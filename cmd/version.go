package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tfadeyi/slotalk/internal/logging"
	"github.com/tfadeyi/slotalk/internal/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Returns the binary build information.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		log := logging.LoggerFromContext(ctx)
		log = log.WithName("version")
		log.Debug(version.BuildInfo())
		log.Info(version.Info())
	},
}
