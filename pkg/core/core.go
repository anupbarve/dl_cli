package core

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dl_cli/pkg/protocol"
	"github.com/dl_cli/pkg/protocol/dl_http"
	"github.com/dl_cli/pkg/utils"
)

type DownloadContextCore struct {
	dlPath string
	urls   []string
}

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
	subDir := t.Format("2006-01-02-15-04-05")

	dstDir := filepath.Join(path, subDir)

	err = os.MkdirAll(dstDir, 0700)
	if err != nil {
		return "", err
	}

	return dstDir, nil
}

func removeDstDir(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func parseURLs(path string) ([]string, error) {
	urlsSlice := strings.Split(path, ",")
	return urlsSlice, nil
}

func Download(path string, urls string) error {

	dstDir, err := createDstDir(path)
	if err != nil {
		fmt.Printf("Target path validation failed. %v", err)
		return err
	}

	urlsSlice, err := parseURLs(urls)
	if err != nil {
		fmt.Printf("Parsing of input urls failed. %v", err)
		return err
	}

	numURLs := len(urlsSlice)
	dlURLs := make([]utils.DownloadURL, numURLs)

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

	// There will be a separate thread per url.
	wg.Add(numURLs)

	for i := 0; i < numURLs; i++ {
		go func(i int) {
			// Worker will be done after current thread is done
			// with it's work.
			defer wg.Done()
			var protocol protocol.Protocol
			// For adding support to future protocols, add the case here.
			switch dlURLs[i].Proto {
			case "http":
				protocol = dl_http.NewDlHttp(dlURLs[i])
				// Core logic for the download
				dlURLs[i].Err = protocol.Download()
			case "https":
				protocol = dl_http.NewDlHttp(dlURLs[i])
				// Core logic for the download
				dlURLs[i].Err = protocol.Download()
			default:
				dlURLs[i].Err = fmt.Errorf("Error: Protocol Unsupported: %s, URL: %s", dlURLs[i].Proto, dlURLs[i].Src)
			}
		}(i)
	}

	// Block till all the downloads are done.
	wg.Wait()
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
