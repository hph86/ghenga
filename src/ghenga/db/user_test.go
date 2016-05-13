package db

import "testing"

func TestUserAdd(t *testing.T) {
	db, cleanup := TestDB(t)
	defer cleanup()

	u, err := NewUser("foo", "bar")
	if err != nil {
		t.Fatal(err)
	}

	err = db.Insert(u)
	if err != nil {
		t.Fatal(err)
	}

	u2, err := FindUser(db, "foo")
	if err != nil {
		t.Fatal(err)
	}

	if !u2.CheckPassword("bar") {
		t.Fatalf("password for test user does not match hash")
	}

	if u2.CheckPassword("xxx") {
		t.Fatalf("wrong password for test user was accepted")
	}
}
