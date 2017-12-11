package common

import (
	"fmt"

	"github.com/fatih/color"
	cli "gopkg.in/urfave/cli.v1"
)

// CliError formats provided error into Red colored string
// and returns cli.NewExitError
func CliError(err error) *cli.ExitError {
	redError := color.RedString(
		fmt.Sprintln("Request failed:", err),
	)
	return cli.NewExitError(redError, 1)
}
