// Package server contains the API server implementation and base
// functionality for ghenga.
package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/jmoiron/modl"
)

// Env is an environment for a handler function.
type Env struct {
	DbMap *modl.DbMap
	Debug bool
}

// HandleFunc is a function similar to http.HandleFunc, but extended by an
// explicit environment parameter. It may return an error.
type HandleFunc func(*Env, http.ResponseWriter, *http.Request) error

// Handler is an http.Handler with an explicit error return value, bundled together with an environment.
type Handler struct {
	Env *Env
	H   HandleFunc
}

// httpWriteJSON encodes the given structures as JSON and writes them to the
// ResponseWriter.
func httpWriteJSON(wr http.ResponseWriter, status int, data interface{}) error {
	wr.Header().Set("Content-Type", "application/json; charset=utf-8")
	wr.WriteHeader(status)

	enc := json.NewEncoder(wr)

	return enc.Encode(data)
}

// jsonError is the struct for an error message returned by the API server.
type jsonError struct {
	Message string `json:"message,omitempty"`
}

// ServeHTTP allows a Handler to be used in place of http.Handler.
func (h Handler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	err := RecoverHandler(h.Env, wr, req, h.H)
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

			if h.Env.Debug {
				je.Message = e.Error()
			}

			err = httpWriteJSON(wr, http.StatusInternalServerError, je)
			if err != nil {
				log.Printf("error writing error document to client: %v", err)
			}
			return
		}
	}
}

// RecoverHandler recovers gracefully from panics that occur when running h.
func RecoverHandler(env *Env, wr http.ResponseWriter, req *http.Request, h HandleFunc) (err error) {
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

	return h(env, wr, req)
}
