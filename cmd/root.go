/*
Copyright © 2023 Oluwole Fadeyi
*/
package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/tfadeyi/sloth-simple-comments/internal/logging"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sli-app",
	Short: "Generate sloth definitions from comments in the application's code.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx context.Context) {
	log := logging.LoggerFromContext(ctx)
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}
}
