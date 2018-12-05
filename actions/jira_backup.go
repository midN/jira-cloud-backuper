package actions

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/midN/jira-cloud-backuper/common"
	"gopkg.in/urfave/cli.v1"
)

type backupResponse struct {
	TaskID string `json:"taskId"`
}

// JiraBackup returns cli.Context related function
// which calls necessary JIRA APIs to initalize a backup action
func JiraBackup() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		backupID, err := initiateJiraBackup(c)
		if err != nil {
			return common.CliError(err)
		}
		fmt.Println("Backup initiated")

		downloadURL, err := common.JiraWaitForBackupReadyness(c, backupID)
		if err != nil {
			return common.CliError(err)
		}

		fmt.Println(color.GreenString(fmt.Sprintln(
			"Done, please use same app to download file or direct link:",
			downloadURL)))
		return nil
	}
}

func initiateJiraBackup(c *cli.Context) (string, error) {
	headers := map[string]string{"Content-Type": "application/json"}
	jsonBody, _ := json.Marshal(common.BackupBody{
		CbAttachments: "true",
		ExportToCloud: "true",
	})

	body, err := common.DoRequest(c, "POST", "/rest/backup/1/export/runbackup", headers, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	respJSON := backupResponse{}
	json.Unmarshal(body, &respJSON)

	return respJSON.TaskID, nil
}
