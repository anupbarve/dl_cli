package core

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dl_cli/pkg/constants"
	"github.com/dl_cli/pkg/protocol"
	"github.com/dl_cli/pkg/protocol/dl_http"
	"github.com/dl_cli/pkg/utils"
)

// Download context for core package
type DownloadContextCore struct {
	dlPath string   //target path
	urls   []string //slice of url strings
}

// Creates destination directories on local filesystem
func createDstDir(path string) (string, error) {
	// Check if path exists
	info, err := os.Stat(path)
	if err == nil {
		// If path exists, check if it is a directory
		if !info.IsDir() {
			return "", fmt.Errorf("Target path is not a directory: %s", path)
		}
		// Check writable permissions
		if info.Mode().Perm()&(1<<(uint(7))) == 0 {
			return "", fmt.Errorf("Write permissions not set on target path: %s", path)
		}
	} else { // Target path does not exist. Creating the target directory

		err = os.MkdirAll(path, 0700)
		if err != nil {
			return "", err
		}
	}

	// Create a sub-dir based on current timestamp to avoid collisions
	t := time.Now()
	subDir := t.Format(constants.DateFormat)
	dstDir := filepath.Join(path, subDir)

	err = os.MkdirAll(dstDir, 0700)
	if err != nil {
		return "", err
	}

	return dstDir, nil
}

// Removes the directory from the path
func removeDstDir(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

// Splits comma separated string into a slice of strings
func parseURLs(path string) ([]string, error) {
	urlsSlice := strings.Split(path, ",")
	return urlsSlice, nil
}

// Core logic for download
func Download(dstPath string, urls string) error {

	// Creates base directory for destination path
	dstDir, err := createDstDir(dstPath)
	if err != nil {
		fmt.Printf("Target path validation failed. %v", err)
		return err
	}

	// Splits comma separated string into a slice of strings
	urlsSlice, err := parseURLs(urls)
	if err != nil {
		fmt.Printf("Parsing of input urls failed. %v", err)
		return err
	}

	numURLs := len(urlsSlice)
	dlURLs := make([]utils.DownloadURL, numURLs)

	// Populate slice of download URLs for each url
	for i, urlStr := range urlsSlice {
		var dlURL utils.DownloadURL
		u, err := url.Parse(urlStr)
		if err != nil {
			// Remove target directory in case of errors
			removeDstDir(dstDir)
			return err
		}
		if u.Scheme == "" || u.Host == "" {
			// Remove target directory in case of errors
			removeDstDir(dstDir)
			return fmt.Errorf("Malformed url string : %s", urlStr)
		}
		dlURL.Src = urlStr
		dlURL.SrcAbs = u.Path
		dlURL.Dst = filepath.Join(dstDir, u.Host)
		dlURL.Proto = u.Scheme
		dlURLs[i] = dlURL
	}

	// Create a waitgroup to track the downloads
	var wg sync.WaitGroup

	// There will be a separate thread per url
	wg.Add(numURLs)

	for i := 0; i < numURLs; i++ {
		go func(i int) {
			// Worker will be done after current thread is done
			// with it's work
			defer wg.Done()
			var protocol protocol.Protocol
			// For adding support to future protocols, add the case here.
			switch dlURLs[i].Proto {
			case "http", "https":
				// Absolute file path is the one used to create a file on local disk
				dlURLs[i].AbsFile = filepath.Join(dlURLs[i].Dst, dlURLs[i].SrcAbs)

				// Calculate absolute directory path, this covers any subdirectory structure
				// in the source URL
				dlURLs[i].AbsDir = path.Dir(dlURLs[i].AbsFile)

				// Create absolute directory path specific directory on local disk
				err := os.MkdirAll(dlURLs[i].AbsDir, 0700)
				if err != nil {
					dlURLs[i].Err = err
					return
				}

				protocol = dl_http.NewDlHttp(dlURLs[i])
				// Core logic for the download
				dlURLs[i].Err = protocol.Download()
			default:
				dlURLs[i].Err = fmt.Errorf("Error: Protocol Unsupported: %s, URL: %s", dlURLs[i].Proto, dlURLs[i].Src)
			}
		}(i)
	}

	// Block till all the downloads are done
	wg.Wait()

	// Print the results on console
	fmt.Printf("***************************\n")
	for i := 0; i < numURLs; i++ {
		fmt.Printf("Source URL   : %s\n", dlURLs[i].Src)
		fmt.Printf("Destination  : %s\n", filepath.Join(dlURLs[i].Dst, dlURLs[i].SrcAbs))
		if dlURLs[i].Err == nil {
			fmt.Printf("Status       : Success\n")
		} else {
			fmt.Printf("Status       : Failed, Reason : %s\n", dlURLs[i].Err)
		}
		fmt.Printf("***************************\n")
	}

	return nil
}
