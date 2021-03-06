package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:  "dl_cli",
	Long: `CLI to download files.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf(err.Error())
	}
}

func init() {
	// To tell Cobra not to provide the default completion command.
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
