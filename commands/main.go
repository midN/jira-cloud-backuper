package commands

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/midN/jira-cloud-backuper/actions"
)

var (
	// SubCommands
	jiraCommand = cli.Command{
		Name:    "jira",
		Usage:   "Backup JIRA Cloud",
		Action:  actions.JiraAction(),
	}

	// Command
	bkupCommand = cli.Command{
		Name:    "backup",
		Usage:   "backup ( jira or confluence )",
		Subcommands: []cli.Command{
			jiraCommand,
		},
	}
)

func Commands() []cli.Command {
	return []cli.Command{
		bkupCommand,
	}
}
