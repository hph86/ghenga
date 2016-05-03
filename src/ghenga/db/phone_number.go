package db

import "fmt"

// PhoneNumber is a phone number of a specified type.
type PhoneNumber struct {
	ID       int64
	Number   string
	Type     string
	PersonID int64
}

// PhoneNumberJSON is the JSON representation of a phone number.
type PhoneNumberJSON struct {
	Type   string `json:"type"`
	Number string `json:"number"`
}

func (p PhoneNumber) String() string {
	return fmt.Sprintf("<PhoneNumber [%v] %v>", p.Type, p.Number)
}

// PhoneNumbers is a collection of phone numbers.
type PhoneNumbers []PhoneNumber

// Equals returns true iff other contains exactly the same phone numbers and
// types.
func (p PhoneNumbers) Equals(other PhoneNumbers) bool {
	found := make(map[PhoneNumber]bool)

	for _, num := range p {
		found[PhoneNumber{Number: num.Number, Type: num.Type}] = false
	}

	for _, num := range other {
		id := PhoneNumber{Number: num.Number, Type: num.Type}
		if _, ok := found[id]; !ok {
			return false
		}
		found[id] = true
	}

	for _, v := range found {
		if !v {
			return false
		}
	}

	return true
}
