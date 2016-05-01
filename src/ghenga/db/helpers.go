package db

import (
	"math/rand"
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

	// for _, d := range []struct{
	// 	probability float32
	// 	tpe string
	// 	gen func() string
	// }{
	// 	{0.5, "mobile", f.CellPhoneNumber},
	// 	{0.9, "work", f.PhoneNumber},
	// 	{0.1, "fax", f.PhoneNumber},
	// 	{0.3, "other", f.PhoneNumber},
	// } {
	// 	if rand.Float32() < d.probability {
	// 		p.PhoneNumbers = append(p.PhoneNumbers, PhoneNumber{
	// 			Type: d.tpe,
	// 			Number: d.gen(),
	// 		})
	// 	}
	// }

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
