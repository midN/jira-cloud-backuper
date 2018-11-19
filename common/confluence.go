package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type confluenceProgressResponse struct {
	Result   string `json:"fileName"`
	Progress string `json:"alternativePercentage"`
	Message  string `json:"currentStatus"`
}

// ConfluenceWaitForBackupReadiness check status of a backup
// and loops until it's ready
func ConfluenceWaitForBackupReadiness(client http.Client, host string) (string, error) {
	downloadURL, fileName, status, progress := "", "", "", ""
	var err error

	for fileName == "" {
		downloadURL, fileName, status, progress, err = confluenceCheckBackupProgress(client, host)
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

func confluenceCheckBackupProgress(client http.Client, host string) (string, string, string, string, error) {
	var respJSON = new(confluenceProgressResponse)
	url := host + "/wiki/rest/obm/1.0/getprogress"
	resp, _ := client.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &respJSON)
	resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", "", "", "", errors.New(string(body))
	}
	return confluenceDownloadURL(respJSON.Result, host), respJSON.Result, respJSON.Message, respJSON.Progress, nil
}

func confluenceDownloadURL(path string, host string) string {
	url := host + "/wiki/download/" + path

	return url
}
