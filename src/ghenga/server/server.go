package server

import "github.com/gorilla/mux"

// NewRouter returns a new router with the complete ghenga API already attached.
func NewRouter(env *Env) *mux.Router {
	router := mux.NewRouter()
	PeopleHandler(env, router)
	LoginHandler(env, router)
	return router
}
