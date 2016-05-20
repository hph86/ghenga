package client

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ParseError returns the error encoded in JSON in the http response.
func ParseError(r *http.Response) Error {
	ct := r.Header.Get("Content-Type")
	if strings.Split(ct, ";")[0] != "application/json" {
		return Error{
			Message: "invalid content type for error message: " + ct,
		}
	}

	var e Error
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&e)
	if err != nil {
		return Error{
			Message: "response body contained invalid JSON: " + err.Error(),
		}
	}

	return e
}

// Error is an error as returned by the ghenga API.
type Error struct {
	Message string `json:"message"`
}

func (e Error) String() string {
	return e.Message
}

func (e Error) Error() string {
	return e.Message
}
