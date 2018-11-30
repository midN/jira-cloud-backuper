package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gopkg.in/urfave/cli.v1"
)

type jiraProgressResponse struct {
	Result   string `json:"result"`
	Progress int    `json:"progress"`
	Message  string `json:"message"`
	Status   string `json:"Status"`
}

// JiraWaitForBackupReadyness check status of a backup
// and loops until it's ready
func JiraWaitForBackupReadyness(c *cli.Context, id string) (string, error) {
	downloadURL, progress := "", 0
	var err error

	for progress < 100 {
		downloadURL, progress, err = jiraCheckBackupProgress(c, id)
		if err != nil {
			return downloadURL, err
		}

		if progress != 100 {
			fmt.Println("Backup is still in progress, status:",
				fmt.Sprint(progress, "%."),
				"Retrying in 10s")
			time.Sleep(10 * time.Second)
		}
	}

	return downloadURL, nil
}

func jiraCheckBackupProgress(c *cli.Context, id string) (string, int, error) {
	respJSON := jiraProgressResponse{}

	// Do request and handle any errors.
	body, err := DoRequest(c, "GET", "/rest/backup/1/export/getProgress?taskId="+id, map[string]string{}, nil)
	if err != nil {
		return "", 0, err
	}

	json.Unmarshal(body, &respJSON)

	// Check to make sure that the job hasn't failed before returning URL.
	if respJSON.Status == "Failed" {
		return "", 0, errors.New(respJSON.Result)
	}

	// Return download path and progress percent.
	return jiraDownloadURL(respJSON.Result), respJSON.Progress, nil
}

func jiraDownloadURL(path string) string {
	url := "/plugins/servlet/" + path

	return url
}
