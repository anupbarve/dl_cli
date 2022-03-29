package cmd

import (
	"fmt"

	"github.com/dl_cli/pkg/core"
	"github.com/spf13/cobra"
)

// usage example
var downloadExample = `
  # Download the files represented by comma separated urls at a specified path.
  dl_cli download -p <path to download> -u <comma separated list of urls>

  # Sample command to download files,
  dl_cli download -p /tmp/test -u http://my.file.com/file,ftp://other.file.com/other,sftp://and.also.this/ending 
 `

// dlCmdDownload -- To download files from list of urls.
var (
	dlCmdDownload = &cobra.Command{
		Use:     "download",
		Short:   "download files from list of urls",
		Example: downloadExample,
		Long:    `download files from list of urls`,
		Run:     dlCmdDownloadRun,
	}
)

type DownloadContext struct {
	dlPath string
	urls   string
}

var dlCtx DownloadContext

func init() {
	rootCmd.AddCommand(dlCmdDownload)
	dlCmdDownload.Flags().StringVarP(&dlCtx.dlPath, "path", "p", "", "Path to download the files")
	dlCmdDownload.Flags().StringVarP(&dlCtx.urls, "urls", "u", "", "List of comma separated urls")
}

// To download files from list of urls
func dlCmdDownloadRun(cmd *cobra.Command, args []string) {
	err := core.Download(dlCtx.dlPath, dlCtx.urls)
	if err != nil {
		fmt.Printf("Download operation failed. %v\n", err)
	} else {
		fmt.Printf("Download operation complete, all URLs processed. \n")
	}
}
