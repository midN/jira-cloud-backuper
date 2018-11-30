package main

import (
	"os"

	"github.com/midN/jira-cloud-backuper/commands"
	"github.com/midN/jira-cloud-backuper/flags"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "JIRA Cloud Backuper"
	app.Usage = "Backup your JIRA and Confluence Cloud"
	app.Version = "2.0"

	app.Flags = flags.Flags()
	app.Commands = commands.Commands()

	app.Run(os.Args)
}
