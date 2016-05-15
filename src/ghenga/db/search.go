package db

import "github.com/jmoiron/modl"

// FuzzyFindPersons searches the database for persons related to the query
// string.
func FuzzyFindPersons(db *modl.DbMap, query string) ([]*Person, error) {
	var result []*Person

	err := db.Select(&result, "SELECT * FROM people WHERE name LIKE ?", "%"+query+"%")
	if err != nil {
		return nil, err
	}

	return result, nil
}
