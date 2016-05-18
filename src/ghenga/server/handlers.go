// Package server contains the API server implementation and base
// functionality for ghenga.
package server

import (
	"encoding/json"
	"errors"
	"ghenga/db"
	"log"
	"net/http"

	"golang.org/x/net/context"

	"github.com/jmoiron/modl"
)

// Env is an environment for a handler function.
type Env struct {
	DbMap *modl.DbMap
	Cfg   Config
}

// HandleFunc is a function similar to http.HandleFunc, but extended by an
// explicit environment parameter. It may return an error.
type HandleFunc func(context.Context, *Env, http.ResponseWriter, *http.Request) error

// httpWriteJSON encodes the given structures as JSON and writes them to the
// ResponseWriter.
func httpWriteJSON(wr http.ResponseWriter, status int, data interface{}) error {
	wr.Header().Set("Content-Type", "application/json; charset=utf-8")
	wr.WriteHeader(status)

	if data == nil {
		_, err := wr.Write([]byte("{}\n"))
		return err
	}

	return json.NewEncoder(wr).Encode(data)
}

// jsonError is the struct for an error message returned by the API server.
type jsonError struct {
	Message string `json:"message,omitempty"`
}

// RecoverHandler recovers gracefully from panics that occur when running h.
func RecoverHandler(ctx context.Context, env *Env, wr http.ResponseWriter, req *http.Request, h HandleFunc) (err error) {
	defer func() {
		// catch panic that may have occurred while running the handler
		if r := recover(); r != nil {
			log.Printf("panic received!")
			log.Printf("r: %v", r)

			e := StatusError{Code: http.StatusInternalServerError}
			switch t := r.(type) {
			case string:
				e.Err = errors.New(t)
			case error:
				e.Err = t
			default:
				e.Err = errors.New("Unknown error")
			}

			err = e
		}
	}()

	return h(ctx, env, wr, req)
}

// RequireAuth ensures that only requests with a valid authentication token are
// passed to H, otherwise an error is returned.
func RequireAuth(h HandleFunc) HandleFunc {
	return func(ctx context.Context, env *Env, res http.ResponseWriter, req *http.Request) error {
		session, err := findSession(env, req)
		if err != nil {
			return err
		}

		ctx = db.NewContextWithSession(ctx, session)

		return h(ctx, env, res, req)
	}
}

// RequireAdmin ensures that only authenticated requests from a user which has
// the admin flag set are passed to H, otherwise an error is returned.
func RequireAdmin(h HandleFunc) HandleFunc {
	return func(ctx context.Context, env *Env, res http.ResponseWriter, req *http.Request) error {
		session, err := findSession(env, req)
		if err != nil {
			return err
		}

		u, err := db.FindUser(env.DbMap, session.User)
		if err != nil {
			return err
		}

		if !u.Admin {
			return StatusError{
				Code: http.StatusForbidden,
				Err:  errors.New("user is not admin"),
			}
		}

		ctx = db.NewContextWithSession(ctx, session)

		return h(ctx, env, res, req)
	}
}

// Handle takes a HandleFunc and returns an http.Handler.
func Handle(ctx context.Context, env *Env, h HandleFunc) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		err := RecoverHandler(ctx, env, wr, req, h)
		if err != nil {
			switch e := err.(type) {
			case Error:
				// return the error to the client as a nicely formatted json document.
				err = httpWriteJSON(wr, e.Status(), jsonError{Message: e.Error()})
				if err != nil {
					log.Printf("error writing error document to client: %v", err)
				}
			default:
				log.Printf("unhandled error: %#v", err)
				je := jsonError{Message: "internal server error"}

				if env.Cfg.Debug {
					je.Message = e.Error()
				}

				err = httpWriteJSON(wr, http.StatusInternalServerError, je)
				if err != nil {
					log.Printf("error writing error document to client: %v", err)
				}
				return
			}
		}
	})
}
