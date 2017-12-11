package actions

import (
	"fmt"

	"github.com/midN/jira-cloud-backuper/common"
	"gopkg.in/urfave/cli.v1"
)

// JiraDownload returns cli.Context related function
// which calls necessary JIRA APIs to download latest backup file
func JiraDownload() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		_, host, err := common.AuthUser(c)
		if err != nil {
			return common.CliError(err)
		}

		fmt.Println(host)
		return nil
	}
}
