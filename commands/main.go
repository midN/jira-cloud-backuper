package commands

import (
	"github.com/midN/jira-cloud-backuper/actions"
	"github.com/midN/jira-cloud-backuper/flags"
	"gopkg.in/urfave/cli.v1"
)

var (
	// SubCommands
	jiraBackupCommand = cli.Command{
		Name:   "jira",
		Usage:  "Backup JIRA Cloud",
		Action: actions.JiraBackup(),
	}

	jiraDownloadCommand = cli.Command{
		Name:   "jira",
		Usage:  "Download latest JIRA Cloud backup",
		Action: actions.JiraDownload(),
	}

	// Commands
	backupCommand = cli.Command{
		Name:    "backup",
		Aliases: []string{"bp"},
		Usage:   "backup ( jira or confluence )",
		Subcommands: []cli.Command{
			jiraBackupCommand,
		},
	}

	downloadCommand = cli.Command{
		Name:    "download",
		Aliases: []string{"dl"},
		Usage:   "download ( jira or confluence )",
		Flags:   flags.DlFlags(),
		Subcommands: []cli.Command{
			jiraDownloadCommand,
		},
	}
)

// Commands returns list of cli.Commands
func Commands() []cli.Command {
	return []cli.Command{
		backupCommand,
		downloadCommand,
	}
}
