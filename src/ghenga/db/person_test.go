package db

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"testing"
	"time"
)

var update = flag.Bool("update", false, "update golden files")

func parseTime(s string) time.Time {
	t, err := time.Parse(timeLayout, s)
	if err != nil {
		panic(err)
	}

	return t
}

var testPersons = []struct {
	name string
	p    Person
}{
	{
		name: "testperson1",
		p: Person{
			ID:           2,
			Name:         "Tamara Skibicki",
			EmailAddress: "pit@ackermannsehls.org",
			PhoneNumbers: []PhoneNumber{
				{Type: "work", Number: "(03867) 3074101"},
				{Type: "mobile", Number: "+49-077-1634655"},
				{Type: "other", Number: "2134"},
			},
			Comment:   "fake profile",
			ChangedAt: parseTime("2016-04-24T10:30:07+02:00"),
			CreatedAt: parseTime("2016-04-24T10:30:07+02:00"),
			Version:   23,
		},
	},
	{
		name: "testperson2",
		p: Person{
			ID:           23,
			Name:         "Mario Drees",
			EmailAddress: "bela_freigang@herweg.com",
			ChangedAt:    parseTime("2016-04-24T10:30:07+00:00"),
			CreatedAt:    parseTime("2016-04-24T10:30:07+00:00"),
			Version:      1,
		},
	},
}

func marshal(t *testing.T, item interface{}) []byte {
	buf, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("json.Marshal(): %v", err)
	}

	return buf
}

func TestPersonMarshal(t *testing.T) {
	for i, test := range testPersons {
		buf := marshal(t, test.p)

		golden := filepath.Join("test-fixtures", "TestPersonMarshal_"+test.name+".golden")
		if *update {
			err := ioutil.WriteFile(golden, buf, 0644)
			if err != nil {
				t.Fatalf("test %d: update golden file %v failed: %v", i, golden, err)
			}
		}

		expected, err := ioutil.ReadFile(golden)
		if err != nil {
			t.Errorf("test %d: unable to read golden file %v", i, golden)
		}
		if !bytes.Equal(buf, expected) {
			t.Errorf("test %d (%v) marshal to JSON failed:\n  want: %s\n   got: %s", i, test.name, expected, buf)
		}
	}
}

var testPersonValidate = []struct {
	name  string
	valid bool
	p     Person
}{
	{
		name:  "invalid1",
		valid: false,
		p: Person{
			Name: "",
		},
	},
}

func TestPersonValidate(t *testing.T) {
	for i, test := range testPersons {
		if err := test.p.Validate(); err != nil {
			t.Errorf("test %v (%v) failed: testPerson is invalid: %v", test.name, i, err)
		}
	}

	for i, test := range testPersonValidate {
		err := test.p.Validate()
		if test.valid && err != nil {
			t.Errorf("test %v (%v) failed: testPerson should be valid but is invalid: %v", test.name, i, err)
		}

		if !test.valid && err == nil {
			t.Errorf("test %v (%v) failed: testPerson should be invalid but is valid", test.name, i)
		}
	}
}

func fakePerson(t *testing.T) *Person {
	p, err := NewFakePerson("de")
	if err != nil {
		t.Fatalf("NewFakePerson(): %v", err)
	}
	p.ID = rand.Int63()
	return p
}

func TestPersonUpdate(t *testing.T) {
	p1 := fakePerson(t)
	p2 := fakePerson(t)

	p1.Update(PersonJSON{Name: &p2.Name})

	// create another fake person

}
