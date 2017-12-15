package common

import (
	"fmt"
	"io"
)

// BackupBody is a struct for a JSON sent to JIRA/Confluence backup endpoint
type BackupBody struct {
	CbAttachments string `json:"cbAttachments"`
	ExportToCloud string `json:"exportToCloud"`
}

// PassThru wraps an existing io.Reader.
//
// It simply forwards the Read() call, while displaying
// the results from individual calls to it
type PassThru struct {
	io.Reader
	Total    int64 // Total # of bytes transferred
	Length   int64 // Expected length
	Progress float64
}

// Read 'overrides' the underlying io.Reader's Read method.
// This is the one that will be called by io.Copy(). We simply
// use it to keep track of byte counts and then forward the call.
func (pt *PassThru) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	if n > 0 {
		pt.Total += int64(n)
		percentage := float64(pt.Total) / float64(pt.Length) * float64(100)
		is := fmt.Sprintf("%6.2f", percentage)
		if percentage-pt.Progress > 5 {
			fmt.Print(is + "%\n")
			pt.Progress = percentage
		}
	}

	return n, err
}
