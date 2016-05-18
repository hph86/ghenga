package server

import (
	"encoding/json"
	"errors"
	"ghenga/db"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/gorilla/mux"
)

// ListPeople handles listing person records.
func ListPeople(ctx context.Context, env *Env, res http.ResponseWriter, req *http.Request) error {
	var people []*db.Person
	err := env.DbMap.Select(&people, "select * from people")
	if err != nil {
		return err
	}

	log.Printf("loaded %v person records", len(people))

	return httpWriteJSON(res, http.StatusOK, people)
}

// ShowPerson returns a Person record.
func ShowPerson(ctx context.Context, env *Env, res http.ResponseWriter, req *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		return StatusError{Code: http.StatusBadRequest, Err: err}
	}

	var person db.Person
	err = env.DbMap.SelectOne(&person, "select * from people where id = ?", id)
	if err != nil {
		return StatusError{
			Err:  errors.New("person not found"),
			Code: http.StatusNotFound,
		}
	}

	return httpWriteJSON(res, http.StatusOK, person)
}

// CreatePerson inserts a new person into the database. The request body must be valid JSON.
func CreatePerson(ctx context.Context, env *Env, wr http.ResponseWriter, req *http.Request) (err error) {
	defer cleanupErr(&err, req.Body.Close)

	var jp db.PersonJSON
	dec := json.NewDecoder(req.Body)
	if err = dec.Decode(&jp); err != nil {
		return err
	}

	var p db.Person
	p.Update(jp)

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
func UpdatePerson(ctx context.Context, env *Env, wr http.ResponseWriter, req *http.Request) (err error) {
	defer cleanupErr(&err, req.Body.Close)

	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		return StatusError{Code: http.StatusBadRequest, Err: err}
	}

	var newPerson db.PersonJSON
	dec := json.NewDecoder(req.Body)
	if err = dec.Decode(&newPerson); err != nil {
		return err
	}

	var p db.Person
	if err = env.DbMap.SelectOne(&p, "select id,created_at,version from people where id = ?", id); err != nil {
		log.Printf("unable to find person ID %v, sql error: %v", id, err)
		return err
	}

	log.Printf("loaded %v from db", p)

	if p.Version != newPerson.Version {
		log.Printf("person record is outdated, version %v != %v",
			p.Version, newPerson.Version)
		return StatusError{
			Err:  errors.New("version field does not match"),
			Code: http.StatusConflict,
		}
	}

	// update all fields except
	p.Update(newPerson)

	p.ChangedAt = time.Now()

	log.Printf("modified %v", p)
	if err = p.Validate(); err != nil {
		return StatusError{Code: http.StatusBadRequest, Err: err}
	}

	log.Printf("save %v", p)
	_, err = env.DbMap.Update(&p)
	if err != nil {
		log.Printf("unable update person %v, sql error: %v", p, err)
		return err
	}

	return httpWriteJSON(wr, http.StatusOK, p)
}

// DeletePerson removes a person from the database.
func DeletePerson(ctx context.Context, env *Env, wr http.ResponseWriter, req *http.Request) (err error) {
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		return StatusError{Code: http.StatusBadRequest, Err: err}
	}

	res := env.DbMap.Dbx.MustExec("delete from people where id = ?", id)

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n != 1 {
		return StatusError{
			Err:  errors.New("person not found"),
			Code: http.StatusNotFound,
		}
	}

	return httpWriteJSON(wr, http.StatusOK, nil)
}

// PeopleHandler adds routes for ghenga API in the given enviroment to r.
func PeopleHandler(ctx context.Context, env *Env, r *mux.Router) {
	r.Handle("/api/person", Handle(ctx, env, RequireAuth(ListPeople))).Methods("GET")
	r.Handle("/api/person", Handle(ctx, env, RequireAuth(CreatePerson))).Methods("POST")
	r.Handle("/api/person/{id}", Handle(ctx, env, RequireAuth(ShowPerson))).Methods("GET")
	r.Handle("/api/person/{id}", Handle(ctx, env, RequireAuth(UpdatePerson))).Methods("PUT")
	r.Handle("/api/person/{id}", Handle(ctx, env, RequireAuth(DeletePerson))).Methods("DELETE")
}
