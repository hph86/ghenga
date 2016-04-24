package db

import "testing"

func TestNewFakePerson(t *testing.T) {
	p, err := NewFakePerson("de")
	if err != nil {
		t.Fatalf("NewFakePerson(): %v", err)
	}

	if err = p.Validate(); err != nil {
		t.Fatalf("NewFakePerson() not valid: %v", err)
	}
}
