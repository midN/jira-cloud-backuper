# jira-cloud-backuper
This application can be used to:

    * Backup JIRA Cloud and Confluence Cloud
    * Download latest JIRA Cloud and Confluence Cloud backup

You will need to create an Atlassian Cloud token in order to use this
application. Instructions for that can be found [here](https://confluence.atlassian.com/cloud/api-tokens-938839638.html).
Use this token in place of your account password when using this
application (i.e. the `--password` flag).

# Usage
Call `jira-cloud-backuper -h` to get list of usable commands or read below.

**Main commands and flags:**

![main](images/main.png)

**Backup commands and flags:**

![backup](images/backup.png)

**Download commands and flags:**

![download](images/download.png)

# Installation
Download latest Release or `go install`
