package server

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func login(t *testing.T, srv *TestSrv, username, password string) (token string, valid time.Duration) {
	req, err := http.NewRequest(http.MethodGet, srv.URL+"/login/token", nil)
	if err != nil {
		t.Fatalf("unable to create login request: %v", err)
	}

	req.SetBasicAuth(username, password)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("login attempt failed: %v", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Fatalf("error closing body: %v", err)
		}
	}()

	var response struct {
		Token    string `json:"token"`
		ValidFor uint   `json:"valid_for"`
	}

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("reading body: %v", err)
	}

	unmarshal(t, buf, &response)

	if response.Token == "" || response.ValidFor == 0 {
		t.Fatalf("invalid response from login endpoint: %v", buf)
	}

	return response.Token, time.Second * time.Duration(response.ValidFor)
}

func TestLogin(t *testing.T) {
	srv, cleanup := TestServer(t)
	defer cleanup()

	token, valid := login(t, srv, "admin", "geheim")
	t.Logf("token: %v, valid for: %v", token, valid)

}
