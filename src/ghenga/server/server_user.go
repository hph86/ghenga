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
	var users []*db.User
	err := env.DbMap.Select(&users, "select * from users")
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

	var user db.User
	err = env.DbMap.SelectOne(&user, "select * from users where id = $1", id)
	if err != nil {
		return StatusError{
			Err:  errors.New("user not found"),
			Code: http.StatusNotFound,
		}
	}

	return httpWriteJSON(res, http.StatusOK, user)
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

	err = env.DbMap.Insert(&u)
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

	var u db.User
	if err = env.DbMap.SelectOne(&u, "select id,created_at,version from users where id = $1", id); err != nil {
		env.Logf("unable to find person ID %v, sql error: %v", id, err)
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

	_, err = env.DbMap.Update(&u)
	if err != nil {
		env.Logf("unable update person %v, sql error: %v", u, err)
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

	res := env.DbMap.Dbx.MustExec("delete from users where id = $1", id)

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n != 1 {
		return StatusError{
			Err:  errors.New("person not found"),
			Code: http.StatusNotFound,
		}
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
