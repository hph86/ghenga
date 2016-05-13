package server

import (
	"errors"
	"ghenga/db"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// LoginResponseJSON is the structure returned by a login request.
type LoginResponseJSON struct {
	Token    string `json:"token"`
	ValidFor uint   `json:"valid_for"`
}

// Login allows users to log in and returns a token.
func Login(env *Env, res http.ResponseWriter, req *http.Request) error {
	username, password, ok := req.BasicAuth()
	if !ok {
		return StatusError{
			Code: http.StatusUnauthorized,
			Err:  errors.New("no login data present"),
		}
	}

	log.Printf("login user %v, password %v", username, password)

	u, err := db.FindUser(env.DbMap, username)
	if err != nil || !u.CheckPassword(password) {
		return StatusError{
			Code: http.StatusUnauthorized,
			Err:  errors.New("invalid username or password"),
		}
	}

	session, err := db.SaveNewSession(env.DbMap, username, env.Cfg.SessionDuration)
	if err != nil {
		return err
	}

	return httpWriteJSON(res, http.StatusOK, LoginResponseJSON{
		Token:    session.Token,
		ValidFor: uint(time.Now().Sub(session.ValidUntil) / time.Second),
	})
}

// LoginHandler adds routes to the for ghenga API in the given enviroment to r.
func LoginHandler(env *Env, r *mux.Router) {
	r.Handle("/login/token", Handler{H: Login, Env: env}).Methods("GET")
}
