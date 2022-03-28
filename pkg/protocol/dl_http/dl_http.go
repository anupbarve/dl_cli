package dl_http

import (
	"fmt"
	"github.com/dl_cli/pkg/utils"
)

type DlHttp struct {
	dlURL utils.DownloadURL
}

// NewDlHttp creates and returns a new instance of Dl_http
func NewDlHttp(dlURL utils.DownloadURL) *DlHttp {
	return &DlHttp{dlURL}
}

func (d *DlHttp) Download() error {
	fmt.Printf("Inside dl_http : %v\n", d.dlURL)
	return nil
}
