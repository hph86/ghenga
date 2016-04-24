package db

import (
	"testing"

	"github.com/jmoiron/modl"
	"github.com/manveru/faker"
)

// NewFakePerson returns a Person struct filled with fake data.
func NewFakePerson(lang string) (*Person, error) {
	f, err := faker.New(lang)
	if err != nil {
		return nil, err
	}

	p := NewPerson(f.FirstName() + " " + f.LastName())

	p.EmailAddress = f.Email()
	p.PhoneMobile = f.CellPhoneNumber()
	p.PhoneWork = f.PhoneNumber()

	p.Comment = "fake profile"

	return p, nil
}

// InsertFakeData will populate the db with fake (but realistic) data.
func InsertFakeData(dbm *modl.DbMap, people int) error {
	for i := 0; i < people; i++ {
		p, err := NewFakePerson("de")
		if err != nil {
			return err
		}

		err = dbm.Insert(p)
		if err != nil {
			return err
		}
	}

	return nil
}

// TestDBFilled returns an in-memory database filled with fake data.
func TestDBFilled(t *testing.T, people int) (*modl.DbMap, func()) {
	db, cleanup := TestDB(t)

	err := InsertFakeData(db, people)
	if err != nil {
		t.Fatalf("TestFillDB(): %v", err)
	}

	return db, cleanup
}
