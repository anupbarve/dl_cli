package utils

type DownloadURL struct {
	Src     string // Source URL
	SrcAbs  string // Absolute Source path
	AbsFile string // Absolute file path
	AbsDir  string // Absolute dir path
	Dst     string // Target path where the file will be downloaded
	Proto   string // Source protocol
	Err     error  // Error context
}
