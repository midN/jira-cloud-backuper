package daemon

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/midN/jira-cloud-backuper/actions"
	cli "gopkg.in/urfave/cli.v1"
)

func daemon(ch chan os.Signal,
	ticker *time.Ticker,
	c *cli.Context,
	backupFunc func(*cli.Context) error,
	downloadFunc func(*cli.Context) error) {

	var err error

	for {
		select {
		case <-ch:
			ticker.Stop()
			fmt.Printf("\nExiting...")
			return
		case <-ticker.C:
			// Get backup function and run it.
			err = backupFunc(c)
			if err != nil {
				// If a backup was initiated less than 24 hours ago then the above
				// function will return an error (due to non-200 status code).
				log.Fatal(err)
			}

			// Get download function and run it.
			err = downloadFunc(c)
			if err != nil {
				// Log errors with download process.
				log.Fatal(err)
			}
		}
	}
}

// StartDaemon starts the program in daemon mode.
func StartDaemon() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		var backupFunc func(*cli.Context) error
		var downloadFunc func(*cli.Context) error

		// Determine which actions to select to run in daemon.
		switch m := c.Command.Name; m {
		case "confluence", "cf":
			backupFunc = actions.ConfluenceBackup()
			downloadFunc = actions.ConfluenceDownload()
		case "jira":
			backupFunc = actions.JiraBackup()
			downloadFunc = actions.JiraDownload()
		default:
			return errors.New("No service (confluence or jira) provided")
		}

		// Create channel and set notification signals.
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Kill, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		// Create a new ticker that runs every 24 hours.
		ticker := time.NewTicker(24 * time.Hour)

		// Run daemon.
		daemon(ch, ticker, c, backupFunc, downloadFunc)

		return nil
	}
}
