// Package db contains the models for the ghenga database.
package db

import (
	"database/sql"
	"testing"

	// import the sqlite driver
	_ "github.com/mattn/go-sqlite3"

	"github.com/jmoiron/modl"
)

// configDBMap creates a new mapping on the given database and creates the
// tables (if necessary).
func configDBMap(db *sql.DB) (*modl.DbMap, error) {
	dbmap := modl.NewDbMap(db, modl.SqliteDialect{})
	dbmap.AddTable(Person{}, "people").SetKeys(true, "ID")

	return dbmap, dbmap.CreateTablesIfNotExists()
}

// Init opens the database. If the database does not exist yet, it is created.
func Init(dbfile string) (*modl.DbMap, error) {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}

	return configDBMap(db)
}

// TestDB returns an in-memory database suitable for testing.
func TestDB(t *testing.T) (*modl.DbMap, func()) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("sql.Open(): %v", err)
	}

	dbmap, err := configDBMap(db)
	if err != nil {
		t.Fatalf("configDBMap(): %v", err)
	}

	return dbmap, func() {
		err := dbmap.Db.Close()
		if err != nil {
			t.Fatalf("db.Close(): %v", err)
		}
	}
}
