package db

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"
)

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

var testUsers = []struct {
	name string
	u    User
}{
	{
		name: "testuser1",
		u: User{
			Login:        "foobar",
			Admin:        false,
			PasswordHash: "foobarbaz",
			ChangedAt:    parseTime("2016-04-24T10:30:07+02:00"),
			CreatedAt:    parseTime("2016-04-24T10:30:07+02:00"),
			Version:      23,
		},
	},
	{
		name: "testuser2",
		u: User{
			Login:        "x",
			Admin:        true,
			PasswordHash: "xxy",
			ChangedAt:    parseTime("2016-03-24T10:30:07+02:00"),
			CreatedAt:    parseTime("2016-01-24T10:30:07+02:00"),
			Version:      5,
		},
	},
}

func TestUserVersion(t *testing.T) {
	db, cleanup := TestDBFilled(t, 20, 3)
	defer cleanup()

	var u User
	err := db.SelectOne(&u, "SELECT * FROM users WHERE id = 2")
	if err != nil {
		t.Fatal(err)
	}

	u.Version = 25
	_, err = db.Update(&u)
	if err == nil {
		t.Fatalf("expected error due to outdated version not found")
	}
}

func TestUserMarshal(t *testing.T) {
	for i, test := range testUsers {
		buf := marshal(t, test.u)

		golden := filepath.Join("testdata", "TestUserMarshal_"+test.name+".golden")
		if *update {
			err := ioutil.WriteFile(golden, buf, 0644)
			if err != nil {
				t.Fatalf("test %d: update golden file %v failed: %v", i, golden, err)
			}
		}

		expected, err := ioutil.ReadFile(golden)
		if err != nil {
			t.Errorf("test %d: unable to read golden file %v", i, golden)
			continue
		}
		if !bytes.Equal(buf, expected) {
			t.Errorf("test %d (%v) wrong JSON returned:\nwant:\n%s\ngot:\n%s", i, test.name, expected, buf)
		}
	}
}

func TestUserUnmarshal(t *testing.T) {
	for i, test := range testUsers {
		golden := filepath.Join("testdata", "TestUserMarshal_"+test.name+".golden")
		buf, err := ioutil.ReadFile(golden)
		if err != nil {
			t.Errorf("test %d: unable to read golden file %v", i, golden)
			continue
		}

		var u User
		unmarshal(t, buf, &u)

		buf2 := marshal(t, u)

		if !bytes.Equal(buf, buf2) {
			t.Errorf("test %d (%v) wrong JSON returned:\nwant:\n%s\ngot:\n%s", i, test.name, buf, buf2)
		}
	}
}

var testUserValidate = []struct {
	name  string
	valid bool
	u     User
}{
	{
		name:  "invalid1",
		valid: false,
		u: User{
			Login: "",
		},
	},
}

func TestUserValidate(t *testing.T) {
	for i, test := range testUsers {
		if err := test.u.Validate(); err != nil {
			t.Errorf("test %v (%v) failed: test User is invalid: %v", test.name, i, err)
		}
	}

	for i, test := range testUserValidate {
		err := test.u.Validate()
		if test.valid && err != nil {
			t.Errorf("test %v (%v) failed: test User should be valid but is invalid: %v", test.name, i, err)
		}

		if !test.valid && err == nil {
			t.Errorf("test %v (%v) failed: test User should be invalid but is valid", test.name, i)
		}
	}
}

func TestUserUpdate(t *testing.T) {
	db, cleanup := TestDBFilled(t, 20, 3)
	defer cleanup()

	u, err := FindUser(db, "user")
	if err != nil {
		t.Fatalf("unable to load user %q: %v", "user", err)
	}

	u.Login = "foo bar"
	if _, err = db.Update(u); err != nil {
		t.Fatalf("unable to update user: %v", err)
	}

	u.Admin = !u.Admin
	u.Version = 1
	if _, err = db.Update(u); err == nil {
		t.Fatalf("update did not fail despite wrong version field")
	}
}

func TestUserUpdatePassword(t *testing.T) {
	db, cleanup := TestDBFilled(t, 20, 3)
	defer cleanup()

	u, err := FindUser(db, "user")
	if err != nil {
		t.Fatalf("unable to load user %q: %v", "user", err)
	}

	if !u.CheckPassword("geheim") {
		t.Fatalf("password for account `user` is not `geheim`")
	}

	u.Password = "foobar2"
	if _, err = db.Update(u); err != nil {
		t.Errorf("unable to update user: %v", err)
	}

	if u.CheckPassword("geheim") {
		t.Errorf("password for account `user` is still `geheim`")
	}

	if !u.CheckPassword("foobar2") {
		t.Errorf("changed password for account `user` is not `foobar2`")
	}

	u2, err := FindUser(db, "user")
	if err != nil {
		t.Fatalf("unable to load user %q: %v", "user", err)
	}

	if !u2.CheckPassword("foobar2") {
		t.Errorf("changed password for account `user` in the db is not `foobar2`")
	}
}
