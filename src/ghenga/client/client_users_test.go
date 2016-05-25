package client

import (
	"ghenga/db"
	"ghenga/server"
	"testing"
)

func TestClientUsersCRUD(t *testing.T) {
	srv, cleanup := server.TestServer(t)
	defer cleanup()

	client := TestClient(t, srv.URL, "admin", "geheim")

	users, err := client.ListUsers()
	if err != nil {
		t.Errorf("error listing users: %v", err)
	}

	t.Logf("found %v users", len(users))

	user := db.User{
		Login:    "x",
		Password: "x",
		Admin:    true,
	}

	u, err := client.CreateUser(user)
	if err != nil {
		t.Errorf("creating user failed: %v", err)
	}

	t.Logf("created user %v", u)
}
