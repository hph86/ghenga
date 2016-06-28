package db

import (
	"testing"
	"time"
)

func TestSession(t *testing.T) {
	for i := 0; i < 20; i++ {
		s, err := NewSession("user", 300)
		if err != nil {
			t.Fatalf("unable to generate new token: %v", err)
		}

		token := s.Token

		if token == "" || len(token) != 2*tokenLength || token == "0000000000000000000000000000000000000000000000000000000000000000" {
			t.Fatalf("invalid token %q", token)
		}
	}
}

func TestSessionSave(t *testing.T) {
	db, cleanup := TestDBFilled(t, 10, 10)
	defer cleanup()

	var tokens []string

	for i := 0; i < 10; i++ {
		session, err := SaveNewSession(db, "user", time.Duration(10*(i-1))*time.Second)
		if err != nil {
			t.Fatalf("NewSession() error %v", err)
		}

		s, err := FindSession(db, session.Token)
		if err != nil {
			t.Fatalf("unable to find newly generated token in the session database: %v", err)
		}

		if s.Token != session.Token {
			t.Fatalf("FindSession returned a different token")
		}

		tokens = append(tokens, s.Token)
	}

	n, err := ExpireSessions(db)
	if err != nil {
		t.Fatalf("error expire sessions: %v", err)
	}

	if n != 2 {
		t.Errorf("expected 2 expired sessions, got %v", n)
	}

	if _, err = FindSession(db, tokens[0]); err == nil {
		t.Fatalf("expired session token %v still found in database", tokens[0])
	}
}
