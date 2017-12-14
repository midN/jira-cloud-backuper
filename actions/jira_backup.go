package actions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"
	"github.com/midN/jira-cloud-backuper/common"
	"gopkg.in/urfave/cli.v1"
)

type backupBody struct {
	CbAttachments string `json:"cbAttachments"`
	ExportToCloud string `json:"exportToCloud"`
}

type backupResponse struct {
	TaskID string `json:"taskId"`
}

// JiraBackup returns cli.Context related function
// which calls necessary JIRA APIs to initalize a backup action
func JiraBackup() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		client, host, err := common.AuthUser(c)
		if err != nil {
			return common.CliError(err)
		}

		backupID, err := initiateBackup(client, host)
		if err != nil {
			return common.CliError(err)
		}

		downloadURL, err := common.JiraWaitForBackupReadyness(client, backupID, host)
		if err != nil {
			return common.CliError(err)
		}

		fmt.Println(color.GreenString(fmt.Sprintln(
			"Done, please use same app to download file or direct link:",
			downloadURL)))
		return nil
	}
}

func initiateBackup(client http.Client, host string) (string, error) {
	var respJSON = new(backupResponse)
	jsonBody, _ := json.Marshal(backupBody{
		"true",
		"true",
	})

	resp, _ := client.Post(
		host+"/rest/backup/1/export/runbackup",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &respJSON)
	resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New(string(body))
	}
	return respJSON.TaskID, nil
}
