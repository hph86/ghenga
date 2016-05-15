package server

import "testing"

func search(t *testing.T, query string) []Person {
	return nil
}

func TestSearchPerson(t *testing.T) {
	srv, cleanup := TestServer(t)
	defer cleanup()

	p := readFixture(t, "sample_person.json")

	token := login(t, srv, "admin", "geheim")

	status, _ := request(t, token, "POST", srv.URL+"/api/person", p)
	if status != 201 {
		t.Fatalf("invalid status code, want 201, got %v, body:\n  %s", status, string(p))
	}
}
