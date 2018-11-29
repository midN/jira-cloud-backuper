package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gopkg.in/urfave/cli.v1"
)

type confluenceProgressResponse struct {
	Result   string `json:"fileName"`
	Progress string `json:"alternativePercentage"`
	Message  string `json:"currentStatus"`
}

// ConfluenceWaitForBackupReadiness check status of a backup
// and loops until it's ready.
func ConfluenceWaitForBackupReadiness(c *cli.Context) (string, error) {
	downloadURL, fileName, status, progress := "", "", "", ""
	var err error

	for fileName == "" {
		downloadURL, fileName, status, progress, err = confluenceCheckBackupProgress(c)
		if err != nil {
			return "", err
		}

		if fileName == "" {
			fmt.Println("Backup is still in progress, status:",
				status,
				"Progress:",
				progress,
				"Retrying in 10s")
			time.Sleep(10 * time.Second)
		}
	}

	return downloadURL, nil
}

func confluenceCheckBackupProgress(c *cli.Context) (string, string, string, string, error) {
	respJSON := confluenceProgressResponse{}

	// Do request and handle any errors.
	body, err := DoRequest(c, "GET", "/wiki/rest/obm/1.0/getprogress", map[string]string{}, nil)
	if err != nil {
		return "", "", "", "", errors.New(string(body))
	}

	json.Unmarshal(body, &respJSON)

	// Return the download URL, filename, current message, and percentage complete.
	return confluenceDownloadURL(c, respJSON.Result), respJSON.Result, respJSON.Message, respJSON.Progress, nil
}

func confluenceDownloadURL(c *cli.Context, path string) string {
	_, _, host, _ := getAtlassianHostParameters(c)
	url := host + "/wiki/download/" + path

	return url
}
