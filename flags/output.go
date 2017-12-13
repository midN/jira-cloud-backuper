package flags

import cli "gopkg.in/urfave/cli.v1"

// DlFlags returns collection of cli.Flags
func DlFlags() []cli.Flag {
	outputFlag := cli.StringFlag{
		Name:  "output, o",
		Usage: "Output to path/file",
		Value: "jira.zip",
	}

	return []cli.Flag{
		outputFlag,
	}
}
