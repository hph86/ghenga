package db

import "golang.org/x/net/context"

// key is used as the key to store values in a context.Context.
type ctxkey int

// sessionKey is used to store a session in a context.Context.
const sessionKey ctxkey = 0

// NewContextWithSession returns a new context which contains the session.
func NewContextWithSession(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}

// SessionFromContext returns the session from the context. If no session is
// available, ok is set to false.
func SessionFromContext(ctx context.Context) (session *Session, ok bool) {
	session, ok = ctx.Value(sessionKey).(*Session)
	return
}
