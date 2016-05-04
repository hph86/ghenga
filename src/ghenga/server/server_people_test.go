package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"
)

func marshal(t *testing.T, item interface{}) string {
	buf, err := json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}

	return string(buf)
}

func postJSON(t *testing.T, url string, data []byte) []byte {
	res, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatalf("POST request to %v failed: %v", url, err)
	}

	t.Logf("POST %v -> %v (%v)", url, res.StatusCode, res.Status)

	if res.StatusCode != 201 {
		t.Fatalf("invalid status code, want 201, got %v, body:\n  %s", res.StatusCode, data)
	}

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("cannot read response body: %v", err)
	}

	err = res.Body.Close()
	if err != nil {
		t.Fatalf("close body: %v", err)
	}

	return responseBody
}

func readFixture(t *testing.T, filename string) []byte {
	p, err := ioutil.ReadFile(filepath.Join("test-fixtures", "sample_person.json"))
	if err != nil {
		t.Fatalf("unable to read test fixture: %v", err)
	}

	return p
}

func TestAddPerson(t *testing.T) {
	srv, cleanup := TestServer(t)
	defer cleanup()

	p := readFixture(t, "sample_person.json")

	postJSON(t, srv.URL+"/api/person", p)
}
