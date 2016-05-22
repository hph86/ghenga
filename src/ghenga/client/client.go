// Package client implements a client for the gheng API.
//
// This package implements a client for the ghenga API. It is mainly used in
// tests against the API implementation itself.
package client

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fd0/probe"
)

// Client is used to communicate with the ghenga API.
type Client struct {
	// the ghenga API base URL
	BaseURL string

	// the authentication token
	Token string

	// end of token validity period
	TokenValidUntil time.Time

	// used by the methods on client to connect to the ghenga API
	C *http.Client
}

// New returns a new Client. In the parameter `url` it expects the base URL for
// the ghenga server to use, e.g. `http://server:8080`.
func New(url string) *Client {
	return &Client{
		BaseURL: url,
		C:       http.DefaultClient,
	}
}

// LoginResponse is returned by the /api/login/token endpoint after successful
// login.
type LoginResponse struct {
	Token    string `json:"token"`
	ValidFor uint   `json:"valid_for"`
}

// Login tries to log into the API with the given credentials. On success, the
// authentication token is returned and stored within the Client struct for
// further use.
func (c *Client) Login(username, password string) (token string, err error) {
	url := c.BaseURL + "/api/login/token"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", probe.Trace(err, username)
	}

	req.SetBasicAuth(username, password)

	res, httpErr := c.C.Do(req)
	if err != nil {
		return "", probe.Trace(httpErr)
	}

	defer func() {
		e := res.Body.Close()
		if err == nil {
			err = e
		}
	}()

	if res.StatusCode != http.StatusOK {
		return "", probe.Trace(ParseError(res))
	}

	var lr LoginResponse
	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&lr); err != nil {
		return "", probe.Trace(err)
	}

	return lr.Token, nil
}
