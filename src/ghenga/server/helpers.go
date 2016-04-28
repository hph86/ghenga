package server

import (
	"ghenga/db"
	"testing"
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
