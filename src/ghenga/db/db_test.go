package db

import (
	"os"
	"testing"

	"github.com/jmoiron/modl"
)

var testDB *modl.DbMap

func TestMain(m *testing.M) {
	var cleanup func()
	testDB, cleanup = TestDB(20, 5)
	res := m.Run()
	cleanup()
	os.Exit(res)
}
