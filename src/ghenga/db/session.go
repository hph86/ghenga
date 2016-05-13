package db

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/jmoiron/modl"
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
func SaveNewSession(db *modl.DbMap, user string, valid time.Duration) (*Session, error) {
	s, err := NewSession(user, valid)
	if err != nil {
		return nil, err
	}

	err = db.Insert(s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// FindSession searches the session with the given token in the database.
func FindSession(db *modl.DbMap, token string) (*Session, error) {
	var s Session
	err := db.SelectOne(&s, "SELECT * FROM sessions WHERE token = ?", token)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// ExpireSessions removes expired sessions from the db.
func ExpireSessions(db *modl.DbMap) (sessionsRemoved int64, err error) {
	res := db.Dbx.MustExec("DELETE FROM sessions WHERE valid_until < datetime('now', 'localtime')")
	return res.RowsAffected()
}
