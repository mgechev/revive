package golint

import (
	"errors"
)

// SendJson sends a JSON object to the server.
func SendJSON(data interface{}) error {
	return nil
}

// ErrInvalidJson is returned when the JSON is invalid.
var ErrInvalidJSON = errors.New("invalid JSON")

// StatusHTTP represents an HTTP status code.
type StatusHTTP int

// Foobar blah blah
func (s StatusHTTP) FooBar() int {
	return int(s)
}

// qux was previously unexported, but now it is exported.
func (s StatusHTTP) Qux() int {
	return int(s)
}

// SendJson sends a JSON object to the server.
func SendJSON(data interface{}) error {
	return nil
}

// errNotFound is returned when the requested resource is not found.
var ErrNotFound = errors.New("not found")

// NEW FORMAT scope changed
// MATCH:23 /comment on exported method StatusHTTP.Qux should start with "Qux ..." using correct scope not "qux ..."/

// NEW FORMAT case change
// MATCH:7 /comment on exported function SendJSON should start with "SendJSON ..." using correct capitalization not "SendJson ..."/
// MATCH:18 /comment on exported method StatusHTTP.FooBar should start with "FooBar ..." using correct capitalization not "Foobar ..."/
// MATCH:28 /comment on exported function SendJSON should start with "SendJSON ..." using correct capitalization not "SendJson ..."/

// TODO FIX OLD FORMAT
// MATCH:12 /comment on exported var ErrInvalidJSON should be of the form "ErrInvalidJSON ..."/
// MATCH:33 /comment on exported var ErrNotFound should be of the form "ErrNotFound ..."/
