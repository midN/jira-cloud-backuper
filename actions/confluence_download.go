package actions

import (
	"fmt"
	"os"

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
		result, err := common.DownloadFile(c, downloadURL, out)
		if err != nil {
			return common.CliError(err)
		}

		fmt.Print(result)
		return nil
	}
}
