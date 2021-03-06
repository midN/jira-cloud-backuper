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

type auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthUser authenticates user with provided username/pw in JIRA using basicauth
func AuthUser(c *cli.Context) (http.Client, string, error) {
	jar, _ := cookiejar.New(nil)
	client := http.Client{Jar: jar}

	user, pw, host, err := defineCredentials(c)
	if err != nil {
		return client, host, err
	}

	auth, _ := json.Marshal(auth{
		Username: user,
		Password: pw,
	})

	resp, _ := client.Post(
		host+"/rest/auth/1/session",
		"application/json",
		bytes.NewBuffer(auth),
	)
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != 200 {
		return client, host, errors.New(string(b))
	}
	return client, host, nil
}

func defineCredentials(c *cli.Context) (string, string, string, error) {
	host := fmt.Sprintf("https://%s.atlassian.net", c.GlobalString("domain"))
	_, err := http.Get(host)
	if err != nil {
		return "", "", "", err
	}
	return c.GlobalString("username"), c.GlobalString("password"), host, nil
}
