package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Person is a person in the database.
type Person struct {
	ID           int64
	Name         string
	EmailAddress string
	PhoneWork    string
	PhoneMobile  string
	PhoneOther   string

	Comment string

	ChangedAt time.Time
	CreatedAt time.Time
}

// PersonJSON is the JSON representation of a Person as returned or consumed by
// the API.
type PersonJSON struct {
	ID           *int64             `json:"id,omitempty"`
	Name         *string            `json:"name,omitempty"`
	EmailAddress *string            `json:"email_address,omitempty"`
	PhoneNumbers *[]PhoneNumberJSON `json:"phone_numbers,omitempty"`

	Comment *string `json:"comment,omitempty"`

	ChangedAt string `json:"changed_at,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

// PhoneNumberJSON is the JSON representation of a phone number.
type PhoneNumberJSON struct {
	Type   string `json:"type"`
	Number string `json:"number"`
}

// NewPerson returns a new person record with the timestamps set to the current
// time.
func NewPerson(name string) *Person {
	ts := time.Now()

	return &Person{
		Name:      name,
		CreatedAt: ts,
		ChangedAt: ts,
	}
}

const timeLayout = "2006-01-02T15:04:05-07:00"

// MarshalJSON returns the JSON representation of p.
func (p *Person) MarshalJSON() ([]byte, error) {
	jp := PersonJSON{
		ID:           &p.ID,
		Name:         &p.Name,
		EmailAddress: &p.EmailAddress,

		Comment:   &p.Comment,
		ChangedAt: p.ChangedAt.Format(timeLayout),
		CreatedAt: p.CreatedAt.Format(timeLayout),
	}

	var numbers []PhoneNumberJSON

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

	return json.Marshal(jp)
}

// Validate checks if p is valid and returns an error if not.
func (p *Person) Validate() error {
	if p.Name == "" {
		return errors.New("name is empty")
	}

	return nil
}

// Update changes the fields present in other.
func (p *Person) Update(other PersonJSON) error {
	if other.Name != nil {
		p.Name = *other.Name
	}

	if other.EmailAddress != nil {
		p.EmailAddress = *other.EmailAddress
	}

	if other.PhoneNumbers != nil {
		for _, pn := range *other.PhoneNumbers {
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

	if other.Comment != nil {
		p.Comment = *other.Comment
	}

	return nil
}

func (p Person) String() string {
	return fmt.Sprintf("<Person %q (%v)>", p.Name, p.ID)
}
