package client

import (
	"ghenga/server"
	"testing"
)

func TestClientLogin(t *testing.T) {
	srv, cleanup := server.TestServer(t)
	defer cleanup()

	c := New(srv.URL)

	token, err := c.Login("x", "x")
	if err == nil {
		t.Errorf("login with invalid credentials suceeded, token %v", token)
	}

	t.Logf("error for invalid credentials is %v", err)

	token, err = c.Login("user", "geheim")
	if err != nil {
		t.Errorf("login with valid credentials failed: %v, token %v", err, token)
	}
}
