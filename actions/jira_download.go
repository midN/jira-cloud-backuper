package actions

import (
	"fmt"
	"os"
	"time"

	"github.com/midN/jira-cloud-backuper/common"
	"gopkg.in/urfave/cli.v1"
)

// JiraDownload returns cli.Context related function
// which calls necessary JIRA APIs to download latest backup file
func JiraDownload() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		filename := c.GlobalString("output")
		if filename == "" {
			timeString := time.Now().Format("2006-01-02")
			filename = fmt.Sprintf("%s-jira.zip", timeString)
		}
		out, err := os.Create(filename)
		if err != nil {
			return common.CliError(err)
		}
		defer out.Close()

		latestID, err := latestJiraTaskID(c)
		if err != nil {
			return common.CliError(err)
		}

		downloadURL, err := common.JiraWaitForBackupReadyness(c, latestID)
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

func latestJiraTaskID(c *cli.Context) (string, error) {
	body, err := common.DoRequest(c, "GET", "/rest/backup/1/export/lastTaskId", map[string]string{}, nil)
	if err != nil {
		return "", nil
	}

	return string(body), nil
}
