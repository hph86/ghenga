package db

import (
	"fmt"
	"time"

	"github.com/elithrar/simple-scrypt"
	"github.com/jmoiron/modl"
)

// User is a user of the system in the database.
type User struct {
	ID           int64
	Login        string
	PasswordHash string

	ChangedAt time.Time
	CreatedAt time.Time
	Version   int64
}

// NewUser returns a new User initialized with the given password.
func NewUser(login, password string) (*User, error) {
	hash, err := scrypt.GenerateFromPassword([]byte(password), scrypt.DefaultParams)
	if err != nil {
		return nil, err
	}

	return &User{
		Login:        login,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		ChangedAt:    time.Now(),
	}, nil
}

// CheckPassword returns true iff the password matches the user's password hash.
func (u User) CheckPassword(password string) bool {
	err := scrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u User) String() string {
	return fmt.Sprintf("<User %v (%v)>", u.Login, u.ID)
}

func FindUser(db *modl.DbMap, login string) (*User, error) {
	var u User
	err := db.SelectOne(&u, "SELECT * FROM users WHERE login = ?", login)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
