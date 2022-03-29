package dl_http

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dl_cli/pkg/utils"
)

// In-memory structure for DlHttp, wrapper on top of DownloadURL structure
type DlHttp struct {
	dlURL utils.DownloadURL
}

// NewDlHttp creates and returns a new instance of Dl_http
func NewDlHttp(dlURL utils.DownloadURL) *DlHttp {
	return &DlHttp{dlURL}
}

// Http protocol specific download function
func (d *DlHttp) Download() error {

	// Create the target file based on absolute file path
	dst, err := os.Create(d.dlURL.AbsFile)
	if err != nil {
		return err
	}

	// Use http library to download actual source file
	resp, err := http.Get(d.dlURL.Src)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Rollback in case of failure
	if resp.StatusCode != http.StatusOK {
		// okay to ignore the removal errors for now
		os.Remove(d.dlURL.AbsFile)
		os.Remove(d.dlURL.AbsDir)
		return fmt.Errorf("Download failed for URL: %s", d.dlURL.Src)
	}

	// io.Copy is going to prevent the memory overflow. Writable chunks from response
	// body to the target file will be in chunks of 32k
	_, err = io.Copy(dst, resp.Body)
	if err != nil {
		// okay to ignore the removal errors for now
		os.Remove(d.dlURL.AbsFile)
		os.Remove(d.dlURL.AbsDir)
		return fmt.Errorf("Download failed, possible out of space: %s", d.dlURL.Src)
	}

	return nil
}
