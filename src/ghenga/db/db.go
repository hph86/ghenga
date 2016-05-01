// Package db contains the models for the ghenga database.
package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	// import the sqlite driver
	_ "github.com/mattn/go-sqlite3"

	"github.com/jmoiron/modl"
	"github.com/jmoiron/sqlx"
	"github.com/rubenv/modl-migrate"
)

const dialect = "sqlite3"

func init() {
	// use our own name mapper
	sqlx.NameMapper = ToSnakeCase
}

func findMigrationsDir() (dir string, err error) {
	subdir := filepath.Join("migrations", dialect)

	dir, err = os.Getwd()
	if err != nil {
		return "", err
	}

	for dir != "" {
		d := filepath.Join(dir, subdir)
		if fi, err := os.Stat(d); err == nil && fi.Mode().IsDir() {
			return d, nil
		}

		dir = filepath.Dir(dir)
	}

	return "", fmt.Errorf("directory %q not found", subdir)
}

func migrateUp(db *sql.DB) error {
	dir, err := findMigrationsDir()
	if err != nil {
		return err
	}

	src := &migrate.FileMigrationSource{Dir: dir}

	_, err = migrate.Exec(db, dialect, src, migrate.Up)
	return err
}

// configDBMap creates a new mapping on the given database and creates the
// tables (if necessary).
func configDBMap(db *sql.DB) (*modl.DbMap, error) {
	dbmap := modl.NewDbMap(db, modl.SqliteDialect{})
	dbmap.AddTableWithName(Person{}, "people").SetKeys(true, "id")
	dbmap.AddTableWithName(PhoneNumber{}, "phone_numbers").SetKeys(true, "id")

	return dbmap, migrateUp(db)
}

// Init opens the database. If the database does not exist yet, it is created.
func Init(dbfile string) (*modl.DbMap, error) {
	db, err := sql.Open(dialect, dbfile)
	if err != nil {
		return nil, err
	}

	return configDBMap(db)
}

// TestDB returns an in-memory database suitable for testing.
func TestDB(t *testing.T) (*modl.DbMap, func()) {
	db, err := sql.Open(dialect, ":memory:")
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
