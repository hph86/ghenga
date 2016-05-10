package db

import (
	"fmt"
	"time"
)

// User is a user of the system in the database.
type User struct {
	ID           int64
	Name         string
	Login        string
	PasswordHash string

	ChangedAt time.Time
	CreatedAt time.Time
	Version   int64
}

func (u User) String() string {
	return fmt.Sprintf("<User[%v] %v (%v)>", u.ID, u.Login, u.Name)
}
