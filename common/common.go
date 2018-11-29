package common

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/midN/jira-cloud-backuper/common"

	cli "gopkg.in/urfave/cli.v1"
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

// GetAtlassianHostParameters returns the username, password, and hostname required
// to send API requests to Atlassian cloud services.
func GetAtlassianHostParameters(c *cli.Context) (string, string, string, error) {
	host := fmt.Sprintf("https://%s.atlassian.net", c.GlobalString("domain"))
	_, err := http.Get(host)
	if err != nil {
		return "", "", "", err
	}
	return c.GlobalString("username"), c.GlobalString("password"), host, nil
}

func DoRequest(c *cli.Context, requestType string, path string, headers map[string]string, body []byte) ([]byte, error) {
	client := http.Client{}

	// Get username, token, and hostname for request.
	user, token, host, err := GetAtlassianHostParameters(c)
	if err != nil {
		return []byte, common.CliError(err)
	}

	// Create the new request.
	req, _ := http.NewRequest(requestType, host+path, body)

	// Add the headers from the headers mapping.
	req.Header.Add("X-Atlassian-Token", "no-check")
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	// Set basic authentication for request.
	req.SetBasicAuth(user, token)

	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return []byte, errors.New(string(body))
	}

	return body, nil
}
