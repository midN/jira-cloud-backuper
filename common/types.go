package common

// BackupBody is a struct for a JSON sent to JIRA/Confluence backup endpoint
type BackupBody struct {
	CbAttachments string `json:"cbAttachments"`
	ExportToCloud string `json:"exportToCloud"`
}
