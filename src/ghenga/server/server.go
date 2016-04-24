package server

import (
	"ghenga/db"
	"log"
	"net/http"
	"os"
	"strconv"

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
	log.Printf("requested person %q", id)

	var person db.Person
	if err = env.DbMap.SelectOne(&person, "select * from people where id = ?", id); err != nil {
		return err
	}

	return httpWriteJSON(res, http.StatusOK, person)
}

// ListenAndServe starts a new ghenga API server with the given environment.
func ListenAndServe(env *Env) (err error) {
	r := mux.NewRouter()

	// API routes
	r.Handle("/api/person", Handler{HandleFunc: ListPeople, Env: env}).Methods("GET")
	r.Handle("/api/person/{id}", Handler{HandleFunc: ShowPerson, Env: env}).Methods("GET")

	// server static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(env.Public)))

	// activate logging to stdout
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, r))

	return http.ListenAndServe(env.ListenAddr, nil)
}
