package common

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/fatih/color"
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

// getAtlassianHostParameters returns the username, password, and hostname required
// to send API requests to Atlassian cloud services.
func getAtlassianHostParameters(c *cli.Context) (string, string, string, error) {
	host := fmt.Sprintf("https://%s.atlassian.net", c.GlobalString("domain"))
	_, err := http.Get(host)
	if err != nil {
		return "", "", "", err
	}
	return c.GlobalString("username"), c.GlobalString("password"), host, nil
}

// DoRequest function provides common interface for sending API requests to Atlassian cloud.
func DoRequest(c *cli.Context, requestType string, path string, headers map[string]string, body io.Reader) ([]byte, error) {
	client := http.Client{}

	// Get username, token, and hostname for request.
	user, token, host, err := getAtlassianHostParameters(c)
	if err != nil {
		return nil, CliError(err)
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

	// Do request and read the body into byte-array.
	resp, _ := client.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(string(respBody))
	}

	return respBody, nil
}

// DownloadFile downloads zipped backup from Confluence or JIRA.
func DownloadFile(c *cli.Context, path string, out *os.File) (string, error) {
	body, err := DoRequest(c, "GET", path, map[string]string{}, nil)
	if err != nil {
		return "", err
	}

	// Initialize PassThru reader and copy file contents to disk.
	contentReader := bytes.NewReader(body)
	readerpt := &PassThru{Reader: contentReader, Length: contentReader.Size()}
	count, err := io.Copy(out, readerpt)
	if err != nil {
		return "", err
	}

	return color.GreenString(fmt.Sprintln(
		"Download finished, file size:", count, "bytes.", "File:", out.Name())), nil
}
