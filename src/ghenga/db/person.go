package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Person is a person in the database.
type Person struct {
	Name         string
	Title        string
	Department   string
	EmailAddress string
	PhoneWork    string
	PhoneMobile  string
	PhoneFax     string
	PhoneOther   string

	// Address
	Street     string
	PostalCode string
	State      string
	City       string
	Country    string

	Comment string

	// the following attributes are managed by the server
	ID        int64
	ChangedAt time.Time
	CreatedAt time.Time
	Version   int64
}

// PersonJSON is the JSON representation of a Person as returned or consumed by
// the API.
type PersonJSON struct {
	ID           *int64             `json:"id,omitempty"`
	Name         *string            `json:"name,omitempty"`
	Title        *string            `json:"title,omitempty"`
	Department   *string            `json:"department,omitempty"`
	EmailAddress *string            `json:"email_address,omitempty"`
	PhoneNumbers *[]PhoneNumberJSON `json:"phone_numbers,omitempty"`

	Address *AddressJSON `json:"address,omitempty"`

	Comment *string `json:"comment,omitempty"`

	ChangedAt string `json:"changed_at,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`

	Version int64 `json:"version"`
}

// PhoneNumberJSON is the JSON representation of a phone number.
type PhoneNumberJSON struct {
	Type   string `json:"type"`
	Number string `json:"number"`
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
		Version:   1,
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

	numbers := []PhoneNumberJSON{}
	for _, num := range []struct{ t, n string }{
		{"work", p.PhoneWork}, {"mobile", p.PhoneMobile}, {"other", p.PhoneOther},
	} {
		if num.n == "" {
			continue
		}

		numbers = append(numbers, PhoneNumberJSON{
			Type:   num.t,
			Number: num.n,
		})
	}
	jp.PhoneNumbers = &numbers

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

	if otherPerson.PhoneNumbers != nil {
		for _, pn := range *otherPerson.PhoneNumbers {
			switch pn.Type {
			case "work":
				p.PhoneWork = pn.Number
			case "mobile":
				p.PhoneMobile = pn.Number
			case "other":
				p.PhoneOther = pn.Number
			default:
				return fmt.Errorf("unknown phone number type %v", pn.Type)
			}
		}
	}

	if otherPerson.Comment != nil {
		p.Comment = *otherPerson.Comment
	}

	return nil
}

func (p Person) String() string {
	return fmt.Sprintf("<Person %q (%v)>", p.Name, p.ID)
}
