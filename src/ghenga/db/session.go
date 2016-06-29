package db

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"time"
)

// Session contains the authentication token of a logged-in user.
type Session struct {
	Token      string
	User       string
	ValidUntil time.Time
}

func (s Session) String() string {
	return fmt.Sprintf("<Session %v, user %v (valid %v)>",
		s.Token[:8], s.User, s.ValidUntil.Sub(time.Now()))
}

const tokenLength = 32

// NewSession generates a new session for a user.
func NewSession(user string, valid time.Duration) (*Session, error) {
	buf := make([]byte, tokenLength)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return nil, err
	}

	s := &Session{
		Token:      hex.EncodeToString(buf),
		User:       user,
		ValidUntil: time.Now().Add(valid),
	}

	return s, nil
}

// SaveNewSession generates a new session for the user and saves it to the db.
func (db *DB) SaveNewSession(user string, valid time.Duration) (*Session, error) {
	s, err := NewSession(user, valid)
	if err != nil {
		return nil, err
	}

	err = db.dbmap.Insert(s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// FindSession searches the session with the given token in the database.
func (db *DB) FindSession(token string) (*Session, error) {
	var s Session
	err := db.dbmap.SelectOne(&s, "SELECT * FROM sessions WHERE token = $1", token)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// ExpireSessions removes expired sessions from the db.
func (db *DB) ExpireSessions() (sessionsRemoved int64, err error) {
	res := db.dbmap.Dbx.MustExec("DELETE FROM sessions WHERE valid_until < now()")
	return res.RowsAffected()
}

// Invalidate removes the session from the database.
func (db *DB) Invalidate(s *Session) error {
	_, err := db.dbmap.Delete(s)
	if err != nil {
		return err
	}

	return nil
}
