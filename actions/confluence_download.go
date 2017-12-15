package actions

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/midN/jira-cloud-backuper/common"
	"gopkg.in/urfave/cli.v1"
)

// ConfluenceDownload returns cli.Context related function
// which calls necessary JIRA APIs to download latest backup file
func ConfluenceDownload() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		filename := c.GlobalString("output")
		if filename == "" {
			filename = "confluence.zip"
		}
		out, err := os.Create(filename)
		if err != nil {
			return common.CliError(err)
		}
		defer out.Close()

		client, host, err := common.AuthUser(c)
		if err != nil {
			return common.CliError(err)
		}

		downloadURL, err := common.ConfluenceWaitForBackupReadyness(client, host)
		if err != nil {
			return common.CliError(err)
		}

		fmt.Println("Downloading to", filename)
		result, err := downloadLatestConfluence(client, downloadURL, out)
		if err != nil {
			return common.CliError(err)
		}

		fmt.Print(result)
		return nil
	}
}

func downloadLatestConfluence(client http.Client, url string, out *os.File) (string, error) {
	resp, _ := client.Get(url)
	if resp.StatusCode == 404 {
		return "", errors.New("File not found at " + url)
	}
	defer resp.Body.Close()

	readerpt := &common.PassThru{Reader: resp.Body, Length: resp.ContentLength}
	count, err := io.Copy(out, readerpt)
	if err != nil {
		return "", err
	}

	return color.GreenString(fmt.Sprintln(
		"Download finished, file size:", count, "bytes.", "File:", out.Name())), nil
}
