package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Login allows users to log in and returns a token.
func Login(env *Env, res http.ResponseWriter, req *http.Request) error {
	return httpWriteJSON(res, http.StatusOK, map[string]interface{}{
		"token":     "1234foobar",
		"valid_for": 7200,
	})
}

// LoginHandler adds routes to the for ghenga API in the given enviroment to r.
func LoginHandler(env *Env, r *mux.Router) {
	r.Handle("/login/token", Handler{H: Login, Env: env}).Methods("GET")
}
