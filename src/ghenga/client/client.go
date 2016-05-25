// Package client implements a client for the gheng API.
//
// This package implements a client for the ghenga API. It is mainly used in
// tests against the API implementation itself.
package client

import (
	"encoding/json"
	"io"
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

	// http request/response tracing
	trace io.Writer
}

// New returns a new Client. In the parameter `url` it expects the base URL for
// the ghenga server to use, e.g. `http://server:8080`.
func New(url string) *Client {
	return &Client{
		BaseURL: url,
		C:       http.DefaultClient,
	}
}

// LoginResponse is returned by the /api/login/token and /api/login/info
// endpoints after successful login.
type LoginResponse struct {
	Token    string `json:"token"`
	ValidFor uint   `json:"valid_for"`
}

// TraceOn enables printing all HTTP requests and responses to the given writer.
func (c *Client) TraceOn(wr io.Writer) {
	c.trace = wr
}

// TraceOff disables debug output for HTTP requests.
func (c *Client) TraceOff() {
	c.trace = nil
}

// do executes the http request req.
func (c *Client) do(req *http.Request) (*http.Response, error) {
	res, err := c.C.Do(req)
	dumpHTTP(c.trace, req, res)
	return res, err
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

	res, httpErr := c.do(req)
	if err != nil {
		return "", probe.Trace(httpErr)
	}

	defer func() {
		e := res.Body.Close()
		if err == nil {
			err = probe.Trace(e)
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

	c.Token = lr.Token
	c.TokenValidUntil = time.Now().Add(time.Duration(lr.ValidFor) * time.Second)

	return lr.Token, nil
}

// newRequest returns a new HTTP request with the token set in the header
// `X-Auth-Token`, if available.
func (c *Client) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, probe.Trace(err, method, url, body)
	}

	if c.Token != "" {
		req.Header.Set("X-Auth-Token", c.Token)
	}

	return req, nil
}

// get executes an HTTP get request, with the authentication header set if a
// token is available.
func (c *Client) get(url string) (*http.Response, error) {
	req, err := c.newRequest("GET", url, nil)
	if err != nil {
		return nil, probe.Trace(err, url)
	}

	res, err := c.do(req)
	if err != nil {
		return nil, probe.Trace(err)
	}

	return res, nil
}

// getJSON executes an HTTP get request using get() and tries to unmarshal the
// response into data. It expects an HTTP status code of 200 (OK).
func (c *Client) getJSON(url string, data interface{}) error {
	res, err := c.get(url)
	if err != nil {
		return probe.Trace(err)
	}

	if res.StatusCode != http.StatusOK {
		return probe.Trace(ParseError(res))
	}

	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(data); err != nil {
		return probe.Trace(err)
	}

	return probe.Trace(res.Body.Close())
}

// Check queries the API server whether the token is still valid.
func (c *Client) Check() error {
	var lr LoginResponse
	err := c.getJSON(c.BaseURL+"/api/login/info", &lr)
	if err != nil {
		return probe.Trace(err)
	}

	return nil
}

// Logout invalidates the session token.
func (c *Client) Logout() error {
	if c.Token == "" {
		return nil
	}

	res, err := c.get(c.BaseURL + "/api/login/invalidate")
	if err != nil {
		return probe.Trace(err)
	}

	if res.StatusCode != http.StatusOK {
		return probe.Trace(ParseError(res))
	}

	return res.Body.Close()
}
