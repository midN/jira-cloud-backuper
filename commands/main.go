package commands

import (
	"github.com/midN/jira-cloud-backuper/actions"
	"github.com/midN/jira-cloud-backuper/daemon"
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

	conflueceBackupCommand = cli.Command{
		Name:    "confluence",
		Aliases: []string{"cf"},
		Usage:   "Backup Confluence Cloud",
		Action:  actions.ConfluenceBackup(),
	}

	jiraDownloadCommand = cli.Command{
		Name:   "jira",
		Usage:  "Download latest JIRA Cloud backup",
		Action: actions.JiraDownload(),
	}

	confluenceDownloadCommand = cli.Command{
		Name:    "confluence",
		Aliases: []string{"cf"},
		Usage:   "Download latest Confluence Cloud backup",
		Action:  actions.ConfluenceDownload(),
	}

	confluenceDaemonCommand = cli.Command{
		Name:    "confluence",
		Aliases: []string{"cf"},
		Usage:   "Backup & download from Confluence Cloud automatically",
		Action:  daemon.StartDaemon(),
	}

	jiraDaemonCommand = cli.Command{
		Name:   "jira",
		Usage:  "Backup & download of JIRA automatically",
		Action: daemon.StartDaemon(),
	}

	// Commands
	backupCommand = cli.Command{
		Name:    "backup",
		Aliases: []string{"bp"},
		Usage:   "backup ( jira or confluence )",
		Subcommands: []cli.Command{
			jiraBackupCommand,
			conflueceBackupCommand,
		},
	}

	downloadCommand = cli.Command{
		Name:    "download",
		Aliases: []string{"dl"},
		Usage:   "download ( jira or confluence )",
		Flags:   flags.DlFlags(),
		Subcommands: []cli.Command{
			jiraDownloadCommand,
			confluenceDownloadCommand,
		},
	}

	daemonCommand = cli.Command{
		Name:    "daemon",
		Aliases: []string{"d"},
		Usage:   "Run daemon to download backups from jira or confluence",
		Subcommands: []cli.Command{
			confluenceDaemonCommand,
			jiraDaemonCommand,
		},
	}
)

// Commands returns list of cli.Commands
func Commands() []cli.Command {
	return []cli.Command{
		backupCommand,
		downloadCommand,
		daemonCommand,
	}
}
