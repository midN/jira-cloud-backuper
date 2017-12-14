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

// ConfluenceBackup returns cli.Context related function
// which calls necessary JIRA APIs to initalize a backup action
func ConfluenceBackup() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		client, host, err := common.AuthUser(c)
		if err != nil {
			return common.CliError(err)
		}

		err = initiateConfluenceBackup(client, host)
		if err != nil {
			return common.CliError(err)
		}
		fmt.Print("Ok good")

		// Cannot check percentage here since Confluence backup API
		// returns fake percentage which goes over 100 lol.
		// Can it reach 9000+?, that is the question.
		downloadURL, err := common.ConfluenceWaitForBackupReadyness(client, host)
		if err != nil {
			return common.CliError(err)
		}

		fmt.Println(color.GreenString(fmt.Sprintln(
			"Done, please use same app to download file or direct link:",
			downloadURL)))
		return nil
	}
}

func initiateConfluenceBackup(client http.Client, host string) error {
	jsonBody, _ := json.Marshal(common.BackupBody{
		"true",
		"true",
	})

	resp, _ := client.Post(
		host+"/wiki/rest/obm/1.0/runbackup",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(string(body))
	}
	return nil
}
