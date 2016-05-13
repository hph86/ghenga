package server

import "time"

// Config configures the ghenga server.
type Config struct {
	SessionDuration time.Duration
	Debug           bool
}
