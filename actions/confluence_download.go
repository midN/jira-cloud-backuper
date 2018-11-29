package actions

import (
	"bytes"
	"fmt"
	"io"
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

		downloadURL, err := common.ConfluenceWaitForBackupReadiness(c)
		if err != nil {
			return common.CliError(err)
		}

		fmt.Println("Downloading to", filename)
		result, err := downloadLatestConfluence(c, downloadURL, out)
		if err != nil {
			return common.CliError(err)
		}

		fmt.Print(result)
		return nil
	}
}

func downloadLatestConfluence(c *cli.Context, path string, out *os.File) (string, error) {
	body, err := common.DoRequest(c, "GET", path, map[string]string{}, nil)
	if err != nil {
		return "", err
	}

	// Initialize PassThru reader and copy file contents to disk.
	contentReader := bytes.NewReader(body)
	readerpt := &common.PassThru{Reader: contentReader, Length: contentReader.Size()}
	count, err := io.Copy(out, readerpt)
	if err != nil {
		return "", err
	}

	return color.GreenString(fmt.Sprintln(
		"Download finished, file size:", count, "bytes.", "File:", out.Name())), nil
}
