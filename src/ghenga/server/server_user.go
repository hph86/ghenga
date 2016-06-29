package server

import (
	"encoding/json"
	"errors"
	"ghenga/db"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

// ListUsers handles listing users.
func ListUsers(ctx context.Context, env *Env, res http.ResponseWriter, req *http.Request) error {
	users, err := env.DB.ListUsers()
	if err != nil {
		return err
	}

	return httpWriteJSON(res, http.StatusOK, users)
}

// ShowUser returns a user record.
func ShowUser(ctx context.Context, env *Env, res http.ResponseWriter, req *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		return StatusError{Code: http.StatusBadRequest, Err: err}
	}

	u, err := env.DB.FindUser(int64(id))
	if err != nil {
		return StatusError{
			Err:  errors.New("user not found"),
			Code: http.StatusNotFound,
		}
	}

	return httpWriteJSON(res, http.StatusOK, u)
}

// CreateUser inserts a new person into the database. The request body must be valid JSON.
func CreateUser(ctx context.Context, env *Env, wr http.ResponseWriter, req *http.Request) (err error) {
	defer cleanupErr(&err, req.Body.Close)

	var ju db.UserJSON
	dec := json.NewDecoder(req.Body)
	if err = dec.Decode(&ju); err != nil {
		return err
	}

	var u db.User
	u.Update(ju)

	// overwrite fields we'd like to be set
	u.CreatedAt = time.Now()
	u.ChangedAt = time.Now()

	if err = u.Validate(); err != nil {
		return StatusError{Code: http.StatusBadRequest, Err: err}
	}

	err = env.DB.InsertUser(&u)
	if err != nil {
		return err
	}

	env.Debugf("created user %v", u)

	return httpWriteJSON(wr, http.StatusCreated, u)
}

// UpdateUser changes an existing user record. The request body must be valid JSON.
func UpdateUser(ctx context.Context, env *Env, wr http.ResponseWriter, req *http.Request) (err error) {
	defer cleanupErr(&err, req.Body.Close)

	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		return StatusError{Code: http.StatusBadRequest, Err: err}
	}

	var newUser db.UserJSON
	dec := json.NewDecoder(req.Body)
	if err = dec.Decode(&newUser); err != nil {
		return err
	}

	u, err := env.DB.FindUser(int64(id))
	if err != nil {
		env.Logf("unable to find person ID %v, error: %v", id, err)
		return err
	}

	if u.Version != newUser.Version {
		env.Debugf("person record is outdated, version %v != %v",
			u.Version, newUser.Version)
		return StatusError{
			Err:  errors.New("version field does not match"),
			Code: http.StatusConflict,
		}
	}

	// update the relevant fields
	u.Update(newUser)
	u.ChangedAt = time.Now()

	if err = u.Validate(); err != nil {
		return StatusError{Code: http.StatusBadRequest, Err: err}
	}

	if err := env.DB.UpdateUser(u); err != nil {
		env.Logf("unable update person %v, error: %v", u, err)
		return err
	}

	return httpWriteJSON(wr, http.StatusOK, u)
}

// DeleteUser removes a person from the database.
func DeleteUser(ctx context.Context, env *Env, wr http.ResponseWriter, req *http.Request) (err error) {
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		return StatusError{Code: http.StatusBadRequest, Err: err}
	}

	if err := env.DB.DeleteUser(int64(id)); err != nil {
		return err
	}

	return httpWriteJSON(wr, http.StatusOK, nil)
}

// UserHandler adds routes for the ghenga API in the given environment to r.
func UserHandler(ctx context.Context, env *Env, r *mux.Router) {
	r.Handle("/api/user", Handle(ctx, env, RequireAdmin(ListUsers))).Methods("GET")
	r.Handle("/api/user", Handle(ctx, env, RequireAdmin(CreateUser))).Methods("Post")
	r.Handle("/api/user/{id}", Handle(ctx, env, RequireAdmin(ShowUser))).Methods("GET")
	r.Handle("/api/user/{id}", Handle(ctx, env, RequireAdmin(UpdateUser))).Methods("PUT")
	r.Handle("/api/user/{id}", Handle(ctx, env, RequireAdmin(DeleteUser))).Methods("DELETE")
}
