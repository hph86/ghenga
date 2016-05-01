package db

import "testing"

var snakeCaseTests = []struct {
	before, after string
}{
	{"foobar", "foobar"},
	{"FooBar", "foo_bar"},
	{"fooBar", "foo_bar"},
	{"PersonID", "person_id"},
	{"PersonIDCollection", "person_id_collection"},
	{"OMGPerson", "omg_person"},
	{"OMG", "omg"},
}

func TestToSnakeCase(t *testing.T) {
	for i, test := range snakeCaseTests {
		got := ToSnakeCase(test.before)

		if got != test.after {
			t.Errorf("test %d of %q failed: wanted %q, got %q", i, test.before, test.after, got)
		}
	}
}
