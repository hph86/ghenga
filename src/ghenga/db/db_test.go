package db

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var TestEntries = []Person{
	{Name: "foo", EmailAddress: "1"},
	{Name: "bar", EmailAddress: "2"},
	{Name: "baz", EmailAddress: "3"},
}

func TestInit(t *testing.T) {
	d, err := ioutil.TempDir("", "ghenga-db-test-")
	if err != nil {
		t.Fatalf("error creating tempdir: %v", err)
	}

	dbfile := filepath.Join(d, "test.db")

	db, err := Init(dbfile)
	if err != nil {
		t.Fatalf("unable to init database: %v", err)
	}

	// add a few entries to the db
	for i, entry := range TestEntries {
		err = db.Insert(&entry)
		if err != nil {
			t.Fatalf("entry %v: insert database: %v", i, err)
		}
	}

	// close the db and open it again
	if err = db.Db.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}

	// reopen db
	db, err = Init(dbfile)
	if err != nil {
		t.Fatalf("unable to init database: %v", err)
	}

	var entries []Person

	// read all entries
	err = db.Select(&entries, "select * from people order by id")
	if err != nil {
		t.Fatalf("select: %v", err)
	}

	if len(entries) != len(TestEntries) {
		t.Fatalf("wrong number of test entries returned: want %d, got %d",
			len(TestEntries), len(entries))
	}

	for i, entry := range entries {
		if entry.Name != TestEntries[i].Name {
			t.Errorf("entry %d wrong name: want %q, got %q", i,
				entry.Name, TestEntries[i].Name)
		}
		if entry.EmailAddress != TestEntries[i].EmailAddress {
			t.Errorf("entry %d wrong email address: want %q, got %q", i,
				entry.EmailAddress, TestEntries[i].EmailAddress)
		}
	}

	if err = db.Db.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}

	if err = os.RemoveAll(d); err != nil {
		t.Fatalf("error removing tempdir: %v", err)
	}
}
