package common

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type progressResponse struct {
	Result   string `json:"result"`
	Progress string `json:"progress"`
	Message  string `json:"message"`
}

// JiraCheckBackupProgress checks current status of a latest backup
func JiraCheckBackupProgress(client http.Client, id string, host string) (string, error) {
	var respJSON = new(progressResponse)
	url := host + "/rest/backup/1/export/getProgress?taskId=" + id
	resp, _ := client.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &respJSON)
	resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New(string(body))
	}
	return jiraDownloadURL(respJSON.Result, host), nil
}

func jiraDownloadURL(path string, host string) string {
	url := host + "/plugins/servlet/" + path

	return url
}
