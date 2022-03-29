package dl_http

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

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
	// Absolute file path is the one used to create a file on local disk
	absFilePath := filepath.Join(d.dlURL.Dst, d.dlURL.SrcAbs)

	// Calculate absolute directory path, this covers any subdirectory structure
	// in the source URL
	absDirPath := path.Dir(absFilePath)

	// Create absolute directory path specific directory on local disk
	err := os.MkdirAll(absDirPath, 0700)
	if err != nil {
		return err
	}

	// Create the target file based on absolute file path
	dst, err := os.Create(absFilePath)
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
		err = os.Remove(absFilePath)
		if err != nil {
			return fmt.Errorf("Download failed for URL: %s, Cleanup of %s failed.", d.dlURL.Src, absFilePath)
		}
		err = os.Remove(absDirPath)
		if err != nil {
			return fmt.Errorf("Download failed for URL: %s, Cleanup of %s failed.", d.dlURL.Src, d.dlURL.Dst)
		}
		return fmt.Errorf("Download failed for URL: %s", d.dlURL.Src)
	}

	// io.Copy is going to prevent the memory overflow. Writable chunks from response
	// body to the target file will be in chunks of 32k
	_, err = io.Copy(dst, resp.Body)
	if err != nil {
		err = os.Remove(absFilePath)
		if err != nil {
			return fmt.Errorf("Download failed for URL: %s, Cleanup of %s failed.", d.dlURL.Src, absFilePath)
		}
		err = os.Remove(absDirPath)
		if err != nil {
			return fmt.Errorf("Download failed for URL: %s, Cleanup of %s failed.", d.dlURL.Src, d.dlURL.Dst)
		}
		return fmt.Errorf("Download failed, possible out of space: %s", d.dlURL.Src)
	}

	return nil
}
