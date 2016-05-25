package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ghenga/db"
	"net/http"

	"github.com/fd0/probe"
)

// ListUsers returns the list of all users.
func (c *Client) ListUsers() ([]db.User, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/api/user", nil)
	if err != nil {
		return nil, probe.Trace(err)
	}

	var list []db.User
	err = c.doJSON(req, http.StatusOK, &list)
	if err != nil {
		return nil, probe.Trace(err)
	}

	return list, nil
}

// FindUser returns the record of a single user, identified by the user ID.
func (c *Client) FindUser(id int) (db.User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/user/%d", c.BaseURL, id), nil)
	if err != nil {
		return db.User{}, probe.Trace(err)
	}

	var user db.User
	err = c.doJSON(req, http.StatusOK, &user)
	if err != nil {
		return db.User{}, probe.Trace(err, id)
	}

	return user, nil
}

// CreateUser creates a new user account.
func (c *Client) CreateUser(u db.User) (db.User, error) {
	data, err := json.Marshal(u)
	if err != nil {
		return db.User{}, probe.Trace(err, u)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/api/user", bytes.NewReader(data))
	if err != nil {
		return db.User{}, probe.Trace(err)
	}

	var resultingUser db.User
	err = c.doJSON(req, http.StatusCreated, &resultingUser)
	if err != nil {
		return db.User{}, probe.Trace(err, u)
	}

	return resultingUser, nil
}
