package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type jiraProgressResponse struct {
	Result   string `json:"result"`
	Progress int    `json:"progress"`
	Message  string `json:"message"`
	Status   string `json:"Status"`
}

// JiraWaitForBackupReadyness check status of a backup
// and loops until it's ready
func JiraWaitForBackupReadyness(client http.Client, id string, host string) (string, error) {
	downloadURL, progress := "", 0
	var err error

	for progress < 100 {
		downloadURL, progress, err = jiraCheckBackupProgress(client, id, host)
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

func jiraCheckBackupProgress(client http.Client, id string, host string) (string, int, error) {
	var respJSON = new(jiraProgressResponse)
	url := host + "/rest/backup/1/export/getProgress?taskId=" + id
	resp, _ := client.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &respJSON)
	resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", 0, errors.New(string(body))
	}

	if respJSON.Status == "Failed" {
		return "", 0, errors.New(respJSON.Result)
	}

	return jiraDownloadURL(respJSON.Result, host), respJSON.Progress, nil
}

func jiraDownloadURL(path string, host string) string {
	url := host + "/plugins/servlet/" + path

	return url
}
