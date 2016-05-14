package server

import (
	"net/http"
	"testing"
)

var testURLs = []struct {
	path   string
	method string
	body   string
}{
	{"/api/login/token", "GET", ""},
	{"/api/login/check", "GET", ""},
	{"/api/person", "GET", ""},
	{"/api/person", "POST", `{"name": "foo"}`},
	{"/api/person/23", "GET", ""},
	{"/api/person/23", "PUT", `{"name": "foo"}`},
}

func TestServerAuthentication(t *testing.T) {
	srv, cleanup := TestServer(t)
	defer cleanup()

	for _, test := range testURLs {
		var body []byte
		if test.body != "" {
			body = []byte(test.body)
		}

		code, _ := request(t, "", test.method, srv.URL+test.path, body)
		if code != http.StatusUnauthorized {
			t.Errorf("%v %v succeeded without authentication", test.method, test.path)
		}
	}

	token := login(t, srv, "admin", "geheim")
	if token == "" {
		t.Fatalf("invalid response for valid login request: token %v", token)
	}

	for _, test := range testURLs {
		var body []byte
		if test.body != "" {
			body = []byte(test.body)
		}

		code, _ := request(t, "", test.method, srv.URL+test.path, body)
		if code == http.StatusUnauthorized {
			t.Errorf("%v %v failed with authentication", test.method, test.path)
		}
	}
}
