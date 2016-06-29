package server

import "testing"

type User struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Admin     bool   `json:"admin"`
	CreatedAt string `json:"created_at"`
	ChangedAt string `json:"changed_at"`
	Version   int    `json:"version"`
}

func createUser(t *testing.T, token, url string, data []byte) []byte {
	status, body := request(t, token, "POST", url, data)
	if status != 201 {
		t.Fatalf("invalid status code, want 201, got %v, body:\n  %s", status, string(data))
	}

	return body
}

func TestUserCRUD(t *testing.T) {
	srv, cleanup := TestServer(t)
	defer cleanup()

	fu := readFixture(t, "sample_user.json")

	token := login(t, srv, "admin", "geheim")

	body := createUser(t, token, srv.URL+"/api/user", fu)

	t.Logf("body: %v", string(body))

	// person := verifyPerson(t, "Nicolai Person", body)

	// status, body = request(t, token, "GET", fmt.Sprintf("%s/api/person/%d", srv.URL, person.ID), nil)
	// if status != 200 {
	// 	t.Fatalf("reading person again yielded unexpected status %d", status)
	// }

	// t.Logf("person: %v", person)

	// person = verifyPerson(t, person.Name, body)
	// person.Name = "Robert Niemand"

	// t.Logf("person: %v", person)

	// status, body = request(t, token, "PUT", fmt.Sprintf("%s/api/person/%d", srv.URL, person.ID), marshal(t, person))
	// if status != 200 {
	// 	t.Fatalf("updating person, invalid status %d", status)
	// }

	// verifyPerson(t, person.Name, body)

	// status, body = request(t, token, "GET", fmt.Sprintf("%s/api/person/%d", srv.URL, person.ID), nil)
	// if status != 200 {
	// 	t.Fatalf("reading person again yielded unexpected status %d", status)
	// }

	// verifyPerson(t, person.Name, body)

	// deletePerson(t, token, srv.URL, person.ID)
}

// func TestUserList(t *testing.T) {
// 	srv, cleanup := TestServer(t)
// 	defer cleanup()

// 	token := login(t, srv, "admin", "geheim")

// 	status, body := request(t, token, "GET", srv.URL+"/api/person", nil)
// 	if status != 200 {
// 		t.Fatalf("reading list of persons failed with invalid status: want 200, got %d", status)
// 	}

// 	var list []Person
// 	unmarshal(t, body, &list)
// 	if len(list) == 0 {
// 		t.Fatalf("got no persons from test server")
// 	}

// 	t.Logf("loaded %d person records", len(list))
// }

// var invalidUserTests = []string{
// 	`{}`,
// 	`{"id": 23}`,
// 	`{"email_address": "foo@example.com"}`,
// }

// func TestInvalidUser(t *testing.T) {
// 	srv, cleanup := TestServer(t)
// 	defer cleanup()

// 	token := login(t, srv, "admin", "geheim")

// 	for _, test := range invalidPersonTests {
// 		status, body := request(t, token, "POST", srv.URL+"/api/person", []byte(test))
// 		if status != 400 {
// 			t.Fatalf("status code for invalid person not found, want 400, got %v, body:\n  %s", status, body)
// 		}
// 	}
// }
