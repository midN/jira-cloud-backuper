package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"

	cli "gopkg.in/urfave/cli.v1"
)

// Auth struct is used as JSON sent to JIRA for authentication
type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthUser authenticates user with provided username/pw in JIRA using basicauth
func AuthUser(username string, password string, host string) (http.Client, error) {
	jar, _ := cookiejar.New(nil)
	client := http.Client{Jar: jar}
	auth, _ := json.Marshal(Auth{
		Username: username,
		Password: password,
	})

	resp, _ := client.Post(
		host+"/rest/auth/1/session",
		"application/json",
		bytes.NewBuffer(auth),
	)
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != 200 {
		return client, errors.New(string(b))
	} else {
		return client, nil
	}
}

// DefineCredentials build host based on domain, validates if given url exists
// and return username/pw/host
func DefineCredentials(c *cli.Context) (string, string, string, error) {
	host := fmt.Sprintf("https://%s.atlassian.net", c.GlobalString("domain"))
	_, err := http.Get(host)
	if err != nil {
		return "", "", "", err
	}
	return c.GlobalString("username"), c.GlobalString("password"), host, nil
}
