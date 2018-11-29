package actions

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/midN/jira-cloud-backuper/common"
	"gopkg.in/urfave/cli.v1"
)

// ConfluenceBackup returns cli.Context related function
// which calls necessary JIRA APIs to initialize a backup action
func ConfluenceBackup() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		err := initiateConfluenceBackup(c)
		if err != nil {
			return common.CliError(err)
		}
		fmt.Println("Backup initiated")

		// Cannot check percentage here since Confluence backup API
		// returns fake percentage which goes over 100 lol.
		// Can it reach 9000+?, that is the question.
		downloadURL, err := common.ConfluenceWaitForBackupReadiness(c)
		if err != nil {
			return common.CliError(err)
		}

		fmt.Println(color.GreenString(fmt.Sprintln(
			"Done, please use same app to download file or direct link:",
			downloadURL)))
		return nil
	}
}

func initiateConfluenceBackup(c *cli.Context) error {
	headers := map[string]string{"Content-Type": "application/json"}
	jsonBody, _ := json.Marshal(common.BackupBody{
		CbAttachments: "true",
		ExportToCloud: "true",
	})

	body, err := common.DoRequest(c, "POST", "/wiki/rest/obm/1.0/runbackup", headers, bytes.NewBuffer(jsonBody))
	if err != nil {
		return common.CliError(err)
	}

	return nil
}
