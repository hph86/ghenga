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
			Name:         "Mario Drees",
			EmailAddress: "bela_freigang@herweg.com",
			ChangedAt:    parseTime("2016-04-24T10:30:07+00:00"),
			CreatedAt:    parseTime("2016-04-24T10:30:07+00:00"),
			Version:      1,
		},
	},
}

func TestPersonInsertSelect(t *testing.T) {
	db, cleanup := TestDB(t)
	defer cleanup()

	var ids []int64
	for _, test := range testPersons {
		err := db.Insert(&test.p)
		if err != nil {
			t.Errorf("saving %v failed: %v", test.name, err)
			continue
		}

		ids = append(ids, test.p.ID)
	}

	for i, test := range testPersons {
		var p Person
		err := db.SelectOne(&p, "SELECT * FROM people WHERE id=?", ids[i])
		if err != nil {
			t.Errorf("loading %v failed: %v", test.p.ID, err)
			continue
		}

		if err = p.LoadPhoneNumbers(db); err != nil {
			t.Errorf("error loading phone numbers: %v", err)
			continue
		}

		if p.ID == 0 {
			t.Errorf("ID of new person is zero")
		}

		if p.Version != test.p.Version+1 {
			t.Errorf("%v: wrong version loaded from db, want %v, got %v",
				test.name, test.p.Version+1, p.Version)
		}

		p.ID = test.p.ID
		p.Version = test.p.Version

		buf1 := marshal(t, test.p)
		buf2 := marshal(t, p)

		if !bytes.Equal(buf1, buf2) {
			t.Errorf("loading %v returned different data:\n  want: %s\n   got: %s",
				test.name, buf1, buf2)
			continue
		}
	}
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
			t.Errorf("test %d (%v) wrong JSON returned:\n  want: %s\n   got: %s", i, test.name, expected, buf)
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
