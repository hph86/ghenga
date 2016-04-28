package server

import (
	"encoding/json"
	"ghenga/db"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// ListPeople handles listing person records.
func ListPeople(env *Env, res http.ResponseWriter, req *http.Request) error {
	people := []db.Person{}
	err := env.DbMap.Select(&people, "select * from people")
	if err != nil {
		return err
	}
	log.Printf("loaded %v person records", len(people))

	return httpWriteJSON(res, http.StatusOK, people)
}

// ShowPerson returns a Person record.
func ShowPerson(env *Env, res http.ResponseWriter, req *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		return StatusError{Code: http.StatusBadRequest, Err: err}
	}
	log.Printf("requested person %v", id)

	var person db.Person
	err = env.DbMap.SelectOne(&person, "select * from people where id = ?", id)
	if err != nil {
		return err
	}

	return httpWriteJSON(res, http.StatusOK, person)
}

// CreatePerson inserts a new person into the database. The request body must be valid JSON.
func CreatePerson(env *Env, wr http.ResponseWriter, req *http.Request) (err error) {
	defer cleanupErr(&err, req.Body.Close)

	var p db.Person
	dec := json.NewDecoder(req.Body)
	if err = dec.Decode(&p); err != nil {
		return err
	}

	// overwrite fields we'd like to be set
	p.CreatedAt = time.Now()
	p.ChangedAt = time.Now()

	if err = p.Validate(); err != nil {
		return StatusError{Code: http.StatusBadRequest, Err: err}
	}

	err = env.DbMap.Insert(&p)
	if err != nil {
		return err
	}

	log.Printf("created person %v", p)

	return httpWriteJSON(wr, http.StatusCreated, p)
}

// UpdatePerson changes an existing person record. The request body must be valid JSON.
func UpdatePerson(env *Env, wr http.ResponseWriter, req *http.Request) (err error) {
	defer cleanupErr(&err, req.Body.Close)

	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		return StatusError{Code: http.StatusBadRequest, Err: err}
	}

	var p db.Person
	if err = env.DbMap.SelectOne(&p, "select * from people where id = ?", id); err != nil {
		return err
	}

	log.Printf("loaded %v from db", p)

	var newPerson db.PersonJSON
	dec := json.NewDecoder(req.Body)
	if err = dec.Decode(&newPerson); err != nil {
		return err
	}

	// only update a few fields from p
	if err = p.Update(newPerson); err != nil {
		return StatusError{Code: http.StatusBadRequest, Err: err}
	}

	p.ChangedAt = time.Now()

	log.Printf("modified %v", p)
	if err = p.Validate(); err != nil {
		return StatusError{Code: http.StatusBadRequest, Err: err}
	}

	log.Printf("save %v", p)
	_, err = env.DbMap.Update(&p)
	if err != nil {
		return err
	}

	return httpWriteJSON(wr, http.StatusCreated, p)
}

// DeletePerson removes a person from the database.
func DeletePerson(env *Env, wr http.ResponseWriter, req *http.Request) (err error) {
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		return StatusError{Code: http.StatusBadRequest, Err: err}
	}

	var p db.Person
	if err = env.DbMap.SelectOne(&p, "select * from people where id = ?", id); err != nil {
		return err
	}

	log.Printf("loaded %v from db", p)

	_, err = env.DbMap.Delete(&p)
	if err != nil {
		return err
	}

	return httpWriteJSON(wr, http.StatusOK, p)
}

// PeopleHandler adds routes to the for ghenga API in the given enviroment to r.
func PeopleHandler(env *Env, r *mux.Router) *mux.Router {
	if r == nil {
		panic("no router given")
	}

	r.Handle("/api/person", Handler{H: ListPeople, Env: env}).Methods("GET")
	r.Handle("/api/person", Handler{H: CreatePerson, Env: env}).Methods("POST")
	r.Handle("/api/person/{id}", Handler{H: ShowPerson, Env: env}).Methods("GET")
	r.Handle("/api/person/{id}", Handler{H: UpdatePerson, Env: env}).Methods("PUT")
	r.Handle("/api/person/{id}", Handler{H: DeletePerson, Env: env}).Methods("DELETE")

	return r
}
