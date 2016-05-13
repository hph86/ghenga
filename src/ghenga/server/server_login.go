package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// LoginResponseJSON is the structure returned by a login request.
type LoginResponseJSON struct {
	Token    string `json:"token"`
	ValidFor uint   `json:"valid_for"`
}

// Login allows users to log in and returns a token.
func Login(env *Env, res http.ResponseWriter, req *http.Request) error {
	return httpWriteJSON(res, http.StatusOK, LoginResponseJSON{
		Token:    "1234foobar",
		ValidFor: 7200,
	})
}

// LoginHandler adds routes to the for ghenga API in the given enviroment to r.
func LoginHandler(env *Env, r *mux.Router) {
	r.Handle("/login/token", Handler{H: Login, Env: env}).Methods("GET")
}
