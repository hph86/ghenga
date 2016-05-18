package db

import (
	"log"
	"math/rand"
	"os"
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
	if rand.Float32() <= 0.2 {
		p.Title = "CEO"
	}
	p.Department = "Testers"
	p.EmailAddress = f.Email()

	for _, d := range []struct {
		probability float32
		tpe         string
		gen         func() string
	}{
		{0.5, "mobile", f.CellPhoneNumber},
		{0.9, "work", f.PhoneNumber},
		{0.1, "fax", f.PhoneNumber},
		{0.3, "other", f.PhoneNumber},
	} {
		if rand.Float32() < d.probability {
			p.PhoneNumbers = append(p.PhoneNumbers, PhoneNumber{
				Type:   d.tpe,
				Number: d.gen(),
			})
		}
	}

	p.Comment = "fake profile"
	p.ID = rand.Int63n(20)

	if rand.Float32() <= 0.6 {
		p.Street = f.StreetAddress()
		p.PostalCode = f.PostCode()
		if rand.Float32() < 0.4 {
			p.State = "CA"
		}
		p.City = f.City()
		p.Country = f.Country()
	}

	return p, nil
}

// NewFakeUser returns a User struct filled with fake data. The password is
// always set to "geheim".
func NewFakeUser(lang string) (*User, error) {
	f, err := faker.New(lang)
	if err != nil {
		return nil, err
	}

	return NewUser(f.UserName(), "geheim")
}

// InsertFakeData will populate the db with fake (but realistic) data. Among
// others, users named "admin" and "user" with the password "geheim" are
// created.
func InsertFakeData(dbm *modl.DbMap, people, user int) error {
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

	for _, s := range []struct {
		name  string
		admin bool
	}{{"admin", true}, {"user", false}} {
		u, err := NewUser(s.name, "geheim")
		if err != nil {
			return err
		}

		u.Admin = s.admin
		if err := dbm.Insert(u); err != nil {
			return err
		}
	}

	for i := 0; i < user; i++ {
		u, err := NewFakeUser("de")
		if err != nil {
			return err
		}

		err = dbm.Insert(u)
		if err != nil {
			// ignore errors for fake data
			continue
		}
	}

	return nil
}

// TestDB returns an in-memory database suitable for testing. If the
// environment variable GHENGA_TEST_DB is set to a file name, this is used
// instead.
func TestDB(t *testing.T) (*modl.DbMap, func()) {
	filename := os.Getenv("GHENGA_TEST_DB")
	if filename == "" {
		filename = ":memory:"
	}

	dbmap, err := Init(filename)
	if err != nil {
		t.Fatalf("unable to initialize db: %v", err)
	}

	if os.Getenv("DBTRACE") != "" {
		dbmap.TraceOn("DB: ", log.New(os.Stderr, "", log.LstdFlags))
	}

	return dbmap, func() {
		err := dbmap.Db.Close()
		if err != nil {
			t.Fatalf("db.Close(): %v", err)
		}
	}
}

// TestDBFilled returns an in-memory database filled with fake data.
func TestDBFilled(t *testing.T, people, user int) (*modl.DbMap, func()) {
	db, cleanup := TestDB(t)

	err := InsertFakeData(db, people, user)
	if err != nil {
		t.Fatalf("TestDBFilled(): %v", err)
	}

	return db, cleanup
}
