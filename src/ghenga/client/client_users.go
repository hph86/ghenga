package client

import (
	"errors"
	"fmt"
	"ghenga/db"

	"github.com/fd0/probe"
)

// ListUsers returns the list of all users.
func (c *Client) ListUsers() ([]db.User, error) {
	var list []db.User
	err := c.getJSON(c.BaseURL+"/api/user", &list)
	if err != nil {
		return nil, probe.Trace(err)
	}

	return list, nil
}

// FindUser returns the record of a single user, identified by the user ID.
func (c *Client) FindUser(id int) (db.User, error) {
	var user db.User
	err := c.getJSON(fmt.Sprintf("%s/api/user/%d", c.BaseURL, id), &user)
	if err != nil {
		return db.User{}, probe.Trace(err, id)
	}

	return user, nil
}

// CreateUser creates a new user account.
func (c *Client) CreateUser(u db.User) (db.User, error) {
	return db.User{}, errors.New("not implemented")
}
