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
	User     string `json:"user"`
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
	if err != nil {
		log.Printf("error finding user %q in database: %v", username, err)
	}

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
		User:     u.Login,
		Token:    session.Token,
		ValidFor: uint(env.Cfg.SessionDuration / time.Second),
	})
}

const authHeaderName = "X-Auth-Token"

// findSession returns a session for the request or an error if none is found.
func findSession(env *Env, req *http.Request) (*db.Session, error) {
	token := req.Header.Get(authHeaderName)
	if token == "" {
		return nil, StatusError{
			Code: http.StatusUnauthorized,
			Err:  errors.New("invalid session token"),
		}
	}

	session, err := db.FindSession(env.DbMap, token)
	if err != nil {
		log.Printf("error finding session with token %q in database: %v", token, err)
	}

	if err != nil || session.ValidUntil.Before(time.Now()) {
		return nil, StatusError{
			Code: http.StatusUnauthorized,
			Err:  errors.New("invalid session token"),
		}
	}

	log.Printf("found session for token %v: %v", token, session)

	return session, nil
}

// Info allows users to check whether a token is still valid and find the
// current username.
func Info(env *Env, res http.ResponseWriter, req *http.Request) error {
	session, err := findSession(env, req)
	if err != nil {
		return err
	}

	return httpWriteJSON(res, http.StatusOK, LoginResponseJSON{
		User:     session.User,
		Token:    session.Token,
		ValidFor: uint(session.ValidUntil.Sub(time.Now()) / time.Second),
	})
}

// Invalidate deletes a valid session token.
func Invalidate(env *Env, res http.ResponseWriter, req *http.Request) error {
	session, err := findSession(env, req)
	if err != nil {
		return err
	}

	return session.Invalidate(env.DbMap)
}

// LoginHandler adds routes to the for ghenga API in the given enviroment to r.
func LoginHandler(env *Env, r *mux.Router) {
	r.Handle("/api/login/token", Handler{H: Login, Env: env}).Methods("GET")
	r.Handle("/api/login/info", Handler{H: Info, Env: env}).Methods("GET")
	r.Handle("/api/login/invalidate", Handler{H: Invalidate, Env: env}).Methods("GET")
}
