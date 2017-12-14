package actions

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/midN/jira-cloud-backuper/common"
	"gopkg.in/urfave/cli.v1"
)

// PassThru wraps an existing io.Reader.
//
// It simply forwards the Read() call, while displaying
// the results from individual calls to it
type PassThru struct {
	io.Reader
	total    int64 // Total # of bytes transferred
	length   int64 // Expected length
	progress float64
}

// Read 'overrides' the underlying io.Reader's Read method.
// This is the one that will be called by io.Copy(). We simply
// use it to keep track of byte counts and then forward the call.
func (pt *PassThru) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	if n > 0 {
		pt.total += int64(n)
		percentage := float64(pt.total) / float64(pt.length) * float64(100)
		is := fmt.Sprintf("%6.2f", percentage)
		if percentage-pt.progress > 5 {
			fmt.Print(is + "%\n")
			pt.progress = percentage
		}
	}

	return n, err
}

// JiraDownload returns cli.Context related function
// which calls necessary JIRA APIs to download latest backup file
func JiraDownload() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		filename := c.GlobalString("output")
		out, err := os.Create(filename)
		if err != nil {
			return common.CliError(err)
		}
		defer out.Close()

		client, host, err := common.AuthUser(c)
		if err != nil {
			return common.CliError(err)
		}

		latestID, err := latestTaskID(client, host)
		if err != nil {
			return common.CliError(err)
		}

		downloadURL, err := common.JiraWaitForBackupReadyness(client, latestID, host)
		if err != nil {
			return common.CliError(err)
		}

		fmt.Println("Downloading to", filename)
		result, err := downloadLatest(client, downloadURL, out)
		if err != nil {
			return common.CliError(err)
		}

		fmt.Print(result)
		return nil
	}
}

func latestTaskID(client http.Client, host string) (string, error) {
	url := host + "/rest/backup/1/export/lastTaskId"
	resp, _ := client.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New(string(body))
	}
	return string(body), nil
}

func downloadLatest(client http.Client, url string, out *os.File) (string, error) {
	resp, _ := client.Get(url)
	if resp.StatusCode == 404 {
		return "", errors.New("File not found at " + url)
	}
	defer resp.Body.Close()

	readerpt := &PassThru{Reader: resp.Body, length: resp.ContentLength}
	count, err := io.Copy(out, readerpt)
	if err != nil {
		return "", err
	}

	return color.GreenString(fmt.Sprintln(
		"Download finished, file size:", count, "bytes.", "File:", out.Name())), nil
}
