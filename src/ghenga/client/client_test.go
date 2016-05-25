package client

import (
	"ghenga/server"
	"os"
	"testing"
)

func TestClientLogin(t *testing.T) {
	srv, cleanup := server.TestServer(t)
	defer cleanup()

	c := New(srv.URL)
	if os.Getenv("HTTPTRACE") != "" {
		c.TraceOn(os.Stderr)
	}

	err := c.Check()
	if err == nil {
		t.Errorf("info endpoint returned no error without a token")
	}

	token, err := c.Login("x", "x")
	if err == nil {
		t.Errorf("login with invalid credentials suceeded, token %v", token)
	}

	token, err = c.Login("user", "geheim")
	if err != nil {
		t.Errorf("login with valid credentials failed: %v, token %v", err, token)
	}

	if err = c.Check(); err != nil {
		t.Errorf("check with valid credentials failed: %v", err)
	}

	if err = c.Logout(); err != nil {
		t.Errorf("error invalidating token: %v", err)
	}

	if err = c.Check(); err == nil {
		t.Errorf("check succeeded after logout")
	}
}
