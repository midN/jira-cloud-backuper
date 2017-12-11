package flags

import "gopkg.in/urfave/cli.v1"

// Flags returns collection of cli.Flags
func Flags() []cli.Flag {
	usernameFlag := cli.StringFlag{
		Name:   "username, u",
		Usage:  "JIRA Username",
		EnvVar: "JIRA_USERNAME",
	}

	passwordFlag := cli.StringFlag{
		Name:   "password, p",
		Usage:  "JIRA Password",
		EnvVar: "JIRA_PASSWORD",
	}

	domainFlag := cli.StringFlag{
		Name:   "domain, d",
		Usage:  "JIRA Domain ( DOMAIN.atlassian.net )",
		EnvVar: "JIRA_DOMAIN",
	}

	return []cli.Flag{
		usernameFlag,
		passwordFlag,
		domainFlag,
	}
}
