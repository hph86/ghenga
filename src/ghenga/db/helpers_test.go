package db

import (
	"testing"

	"github.com/jmoiron/modl"
)

func TestNewFakePerson(t *testing.T) {
	p, err := NewFakePerson("de")
	if err != nil {
		t.Fatalf("NewFakePerson(): %v", err)
	}

	if err = p.Validate(); err != nil {
		t.Fatalf("NewFakePerson() not valid: %v", err)
	}
}

func countPeople(t *testing.T, db *modl.DbMap) int {
	count := 0
	if err := db.Dbx.Get(&count, "select count(*) from people"); err != nil {
		t.Fatalf("Dbx.Get(): %v", err)
	}

	return count
}

func TestInsertFakeData(t *testing.T) {
	db, cleanup := TestDB(t)
	defer cleanup()

	before := countPeople(t, db)
	if err := InsertFakeData(db, 234); err != nil {
		t.Fatalf("InsertFakeData(): %v", err)
	}

	after := countPeople(t, db)
	if before >= after {
		t.Fatal("no additional people inserted by InsertFakeData()")
	}
}

func TestTestDBFilled(t *testing.T) {
	db, cleanup := TestDBFilled(t, 234)
	defer cleanup()

	people := countPeople(t, db)
	if people < 234 {
		t.Fatalf("expcted at least %v people in the db, got %v", 234, people)
	}
}
