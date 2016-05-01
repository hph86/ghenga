package db

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
