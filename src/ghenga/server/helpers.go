package server

// cleanupErr runs fn and sets err to the returned error if err is nil.
func cleanupErr(err *error, fn func() error) {
	e := fn()
	if *err == nil {
		*err = e
	}
}
