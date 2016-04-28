package server

import (
	"ghenga/db"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// cleanupErr runs fn and sets err to the returned error if err is nil.
func cleanupErr(err *error, fn func() error) {
	e := fn()
	if *err == nil {
		*err = e
	}
}

const fakePersonProfiles = 200

// TestEnv returns a test environment running on an in-memory database filled
// with test data.
func TestEnv(t *testing.T) (env *Env, cleanup func()) {
	db, dbcleanup := db.TestDBFilled(t, fakePersonProfiles)

	env = &Env{
		DbMap: db,
	}

	return env, func() { dbcleanup() }
}

// TestServer returns an *httptest.Server running the ghenga API on an
// in-memory DB filled with fake data.
func TestServer(t *testing.T) (srv *httptest.Server, cleanup func()) {
	env, envcleanup := TestEnv(t)

	r := mux.NewRouter()
	PeopleHandler(env, r)
	srv = httptest.NewServer(r)

	return srv, func() {
		srv.Close()
		envcleanup()
	}
}
