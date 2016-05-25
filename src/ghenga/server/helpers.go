package server

import (
	"ghenga/db"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"golang.org/x/net/context"
)

// cleanupErr runs fn and sets err to the returned error if err is nil.
func cleanupErr(err *error, fn func() error) {
	e := fn()
	if *err == nil {
		*err = e
	}
}

const (
	fakePersonProfiles = 200
	fakeUserProfiles   = 5
)

// TestEnv returns a test environment running on an in-memory database filled
// with test data.
func TestEnv(t *testing.T) (env *Env, cleanup func()) {
	db, dbcleanup := db.TestDBFilled(t, fakePersonProfiles, fakeUserProfiles)

	env = &Env{
		DbMap: db,
		Cfg: Config{
			SessionDuration: 600 * time.Second,
		},
	}

	if os.Getenv("SERVERTRACE") != "" {
		env.Logger.Debug = log.New(os.Stdout, "", log.LstdFlags)
		env.Logger.Error = log.New(os.Stdout, "", log.LstdFlags)
	}

	return env, func() { dbcleanup() }
}

// TestSrv bundles a test server with a test environment.
type TestSrv struct {
	*httptest.Server
	*Env
}

// TestServer returns an *httptest.Server running the ghenga API on an
// in-memory DB filled with fake data.
func TestServer(t *testing.T) (srv *TestSrv, cleanup func()) {
	env, envcleanup := TestEnv(t)

	ctx, cancel := context.WithCancel(context.TODO())

	srv = &TestSrv{
		Server: httptest.NewServer(NewRouter(ctx, env)),
		Env:    env,
	}

	return srv, func() {
		srv.Close()
		envcleanup()
		cancel()
	}
}
