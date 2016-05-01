package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/modl"
)

// Person is a person in the database.
type Person struct {
	ID           int64
	Name         string
	Title        string
	Department   string
	EmailAddress string
	PhoneNumbers []PhoneNumber `db:"-"`

	// Address
	Street     string
	PostalCode string
	State      string
	City       string
	Country    string

	Comment string

	ChangedAt time.Time
	CreatedAt time.Time
	Version   int64
}

// PersonJSON is the JSON representation of a Person as returned or consumed by
// the API.
type PersonJSON struct {
	ID           *int64            `json:"id,omitempty"`
	Name         *string           `json:"name,omitempty"`
	Title        *string           `json:"title,omitempty"`
	Department   *string           `json:"department,omitempty"`
	EmailAddress *string           `json:"email_address,omitempty"`
	PhoneNumbers []PhoneNumberJSON `json:"phone_numbers"`

	Address *AddressJSON `json:"address,omitempty"`

	Comment *string `json:"comment,omitempty"`

	ChangedAt string `json:"changed_at,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`

	Version int64 `json:"version"`
}

// AddressJSON is the JSON representation of an address.
type AddressJSON struct {
	Street     string `json:"street,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	State      string `json:"state,omitempty"`
	City       string `json:"city,omitempty"`
	Country    string `json:"country,omitempty"`
}

// NewPerson returns a new person record.
func NewPerson(name string) *Person {
	ts := time.Now()
	return &Person{
		Name:      name,
		CreatedAt: ts,
		ChangedAt: ts,
	}
}

// timeLayout is the format used for the timestamps in the JSON representation
// of a Person.
const timeLayout = "2006-01-02T15:04:05-07:00"

// MarshalJSON returns the JSON representation of p.
func (p Person) MarshalJSON() ([]byte, error) {
	jp := PersonJSON{
		ID:   &p.ID,
		Name: &p.Name,

		ChangedAt: p.ChangedAt.Format(timeLayout),
		CreatedAt: p.CreatedAt.Format(timeLayout),
		Version:   p.Version,
	}

	if p.Title != "" {
		jp.Title = &p.Title
	}

	if p.Department != "" {
		jp.Department = &p.Department
	}

	if p.EmailAddress != "" {
		jp.EmailAddress = &p.EmailAddress
	}

	jp.PhoneNumbers = []PhoneNumberJSON{}
	for _, pn := range p.PhoneNumbers {
		jp.PhoneNumbers = append(jp.PhoneNumbers, PhoneNumberJSON{
			Type:   pn.Type,
			Number: pn.Number,
		})
	}

	if p.Street != "" || p.PostalCode != "" || p.State != "" || p.City != "" || p.Country != "" {
		jp.Address = &AddressJSON{
			Street:     p.Street,
			PostalCode: p.PostalCode,
			State:      p.State,
			City:       p.City,
			Country:    p.Country,
		}
	}

	if p.Comment != "" {
		jp.Comment = &p.Comment
	}

	return json.Marshal(jp)
}

// Validate checks if p is valid and returns an error if not.
func (p *Person) Validate() error {
	if p.Name == "" {
		return errors.New("name is empty")
	}

	return nil
}

// Update changes the fields present in otherPerson.
func (p *Person) Update(otherPerson PersonJSON) error {
	if otherPerson.Name != nil {
		p.Name = *otherPerson.Name
	}

	if otherPerson.EmailAddress != nil {
		p.EmailAddress = *otherPerson.EmailAddress
	}

	if otherPerson.Comment != nil {
		p.Comment = *otherPerson.Comment
	}

	return nil
}

// PostInsert is run after a person is saved into the database. It is
// used to handle phone numbers associated with a person.
func (p *Person) PostInsert(db modl.SqlExecutor) error {
	if len(p.PhoneNumbers) == 0 {
		return nil
	}

	for _, num := range p.PhoneNumbers {
		num.PersonID = p.ID
		err := db.Insert(&num)
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadPhoneNumbers loads the phone numbers associated with the person.
func (p *Person) LoadPhoneNumbers(db *modl.DbMap) error {
	return db.Select(&p.PhoneNumbers, "SELECT * FROM phone_numbers WHERE person_id = ?", p.ID)
}

func (p Person) String() string {
	numbers := ""
	if len(p.PhoneNumbers) > 0 {
		for _, num := range p.PhoneNumbers {
			numbers += fmt.Sprintf(", %v [%v]", num.Number, num.Type)
		}
	}

	return fmt.Sprintf("<Person (%v)%s>", p.Name, numbers)
}
