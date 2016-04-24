package db

import (
	"github.com/jmoiron/modl"
	"github.com/manveru/faker"
)

// TestNewFakePerson returns a Person struct filled with fake data.
func TestNewFakePerson(lang string) (*Person, error) {
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

// TestFillDB will populate the db with fake (but realistic) data.
func TestFillDB(dbm *modl.DbMap, people int) error {
	for i := 0; i < people; i++ {
		p, err := TestNewFakePerson("de")
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
