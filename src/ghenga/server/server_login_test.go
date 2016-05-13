package server

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func loginRequest(t *testing.T, srv *TestSrv, username, password string) (status int, body []byte) {
	req, err := http.NewRequest(http.MethodGet, srv.URL+"/login/token", nil)
	if err != nil {
		t.Fatalf("unable to create login request: %v", err)
	}

	if username != "" {
		req.SetBasicAuth(username, password)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("login attempt failed: %v", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Fatalf("error closing body: %v", err)
		}
	}()

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("reading body: %v", err)
	}

	return res.StatusCode, buf
}

func login(t *testing.T, srv *TestSrv, username, password string) (token string) {
	status, body := loginRequest(t, srv, username, password)

	if status != http.StatusOK {
		t.Fatalf("unexpected status %v, expected 200", status)
	}

	var response struct {
		Token    string `json:"token"`
		ValidFor uint   `json:"valid_for"`
	}

	unmarshal(t, body, &response)

	if response.Token == "" || response.ValidFor == 0 {
		t.Fatalf("invalid response from login endpoint: %v", body)
	}

	return response.Token
}

var invalidUsernamePasswords = []struct {
	u string
	p string
}{
	{"admin", "geheimX"},
	{"admin", ""},
	{"", "geheimX"},
	{"", ""},
}

func TestLogin(t *testing.T) {
	srv, cleanup := TestServer(t)
	defer cleanup()

	token := login(t, srv, "admin", "geheim")
	if token == "" {
		t.Fatalf("invalid response for valid login request: token %v", token)
	}

	for _, test := range invalidUsernamePasswords {
		status, body := loginRequest(t, srv, test.u, test.p)
		if status != http.StatusUnauthorized {
			t.Errorf("invalid response for invalid login request (%v/%v): status %v, body:\n%s",
				test.u, test.p, status, body)
		}
	}
}
