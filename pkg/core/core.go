package core

import (
	"os"
	"fmt"
	"strings"
	"sync"
	"time"
	"path/filepath"
	"net/url"

	"github.com/dl_cli/pkg/protocol"
	"github.com/dl_cli/pkg/utils"
	"github.com/dl_cli/pkg/protocol/dl_http"
)

type DownloadContextCore struct {
	dlPath string
	urls   []string
}

func createDstDir(path string) (string, error) {

	// TODO, check permissions for base dir

	err := os.MkdirAll(path, 0700)
	if err != nil {
		return "", err
	}

	t := time.Now()
	subDir := t.Format("2006-01-02-15-04-05")

	dstDir := filepath.Join(path, subDir)

        err = os.MkdirAll(dstDir, 0700)
        if err != nil {
                return "", err
        }

	return dstDir, nil
}

func parseURLs(path string) ([]string, error) {
	urlsSlice := strings.Split(path, ",")
	// TODO validate the URL
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
			return err
		}
		dlURL.Src = urlStr
		dlURL.Dst = filepath.Join(dstDir, u.Host)
		dlURL.IsDir = false
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
				dlURLs[i].Err = protocol.Download()
			default:
				dlURLs[i].Err = fmt.Errorf("Error: Protocol Unsupported %s, URL %s", dlURLs[i].Proto, dlURLs[i].Src)
			}
		}(i)
	}

	// Block till all the downloads are done.
	wg.Wait()
	for i := 0; i < numURLs; i++ {
		fmt.Println(dlURLs[i])
	}

	return nil
}
