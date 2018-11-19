package common

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"

	"gopkg.in/urfave/cli.v1"
)

// AuthUser authenticates user with provided username/token against Atlassian REST API using basicauth.
func AuthUser(c *cli.Context) (http.Client, string, error) {
	jar, _ := cookiejar.New(nil)
	client := http.Client{Jar: jar}

	user, pw, host, err := defineCredentials(c)
	if err != nil {
		return client, host, err
	}

	// Create a new request and add headers.
	req, err := http.NewRequest("POST", host, nil)
	req.Header.Add("X-Atlassian-Token", "no-check")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	// Set basic authentication with username and token.
	req.SetBasicAuth(user, pw)

	// Set request.
	resp, err := client.Do(req)

	b, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	// Check for valid status to return correct error code.
	var respError error
	switch c := resp.StatusCode; c {
	case 200:
		respError = nil
	case 302:
		// This case sometimes happens when you try to authenticate against a cloud instance running a single product
		// i.e. Confluence or JIRA (but not both).
		respError = nil
	default:
		respError = errors.New(string(b))
	}

	return client, host, respError
}

func defineCredentials(c *cli.Context) (string, string, string, error) {
	host := fmt.Sprintf("https://%s.atlassian.net", c.GlobalString("domain"))
	_, err := http.Get(host)
	if err != nil {
		return "", "", "", err
	}
	return c.GlobalString("username"), c.GlobalString("password"), host, nil
}
