package cmd

import (
	"fmt"

	"github.com/dl_cli/pkg/constants"
	"github.com/spf13/cobra"
)

// usage example
var versionExample = `
  # Check the current version of downloader CLI being used.
  dl_cli version
 `

// versionCmd represents "Version of file downloader being used.".
var (
	versionCmd = &cobra.Command{
		Use:     "version",
		Short:   "Current version of file downloader CLI being used",
		Example: versionExample,
		Long:    `Current version of file downloader CLI being used`,
		Run: func(cmd *cobra.Command, args []string) {
			//Prints the current version of file downloader CLI being used.
			fmt.Println(constants.CLIVersion)
		},
	}
)

func init() {
	// Add version command to the root command
	rootCmd.AddCommand(versionCmd)
}
