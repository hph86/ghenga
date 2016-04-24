package server

import (
	"encoding/json"
	"ghenga/db"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
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
	if err = env.DbMap.SelectOne(&person, "select * from people where id = ?", id); err != nil {
		return err
	}

	return httpWriteJSON(res, http.StatusOK, person)
}

// CreatePerson inserts a new person into the database. The request body must be valid JSON.
func CreatePerson(env *Env, res http.ResponseWriter, req *http.Request) (err error) {
	defer cleanupErr(&err, req.Body.Close)

	var p db.Person
	dec := json.NewDecoder(req.Body)
	if err = dec.Decode(&p); err != nil {
		return err
	}

	p.CreatedAt = time.Now()
	p.ChangedAt = time.Now()

	if err = p.Validate(); err != nil {
		return StatusError{Code: http.StatusBadRequest, Err: err}
	}

	err = env.DbMap.Insert(&p)
	if err != nil {
		return err
	}

	return StatusError{Code: http.StatusCreated}
}

// ListenAndServe starts a new ghenga API server with the given environment.
func ListenAndServe(env *Env) (err error) {
	r := mux.NewRouter()

	// API routes
	r.Handle("/api/person", Handler{H: ListPeople, Env: env}).Methods("GET")
	r.Handle("/api/person/{id}", Handler{H: ShowPerson, Env: env}).Methods("GET")
	r.Handle("/api/person", Handler{H: CreatePerson, Env: env}).Methods("POST")

	// server static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(env.Public)))

	// activate logging to stdout
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, r))

	return http.ListenAndServe(env.ListenAddr, nil)
}
