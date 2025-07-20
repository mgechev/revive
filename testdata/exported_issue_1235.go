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

// scope changed
// MATCH:23 /comment on exported method StatusHTTP.Qux should be of the form "Qux ..." to match its exported status, not "qux ..."/
// MATCH:33 /comment on exported var ErrNotFound should be of the form "ErrNotFound ..." to match its exported status, not "errNotFound ..."/

// case change
// MATCH:7 /comment on exported function SendJSON should be of the form "SendJSON ..." by using its correct casing, not "SendJson ..."/
// MATCH:12 /comment on exported var ErrInvalidJSON should be of the form "ErrInvalidJSON ..." by using its correct casing, not "ErrInvalidJson ..."/
// MATCH:18 /comment on exported method StatusHTTP.FooBar should be of the form "FooBar ..." by using its correct casing, not "Foobar ..."/
// MATCH:28 /comment on exported function SendJSON should be of the form "SendJSON ..." by using its correct casing, not "SendJson ..."/

// VeryLongCommentThatCouldBeCJKThatCannotBeSplitOnSpaces is about the function F.
func F() string {
	return "F"
}

// This one is a safeguard against future changes in the way the error is reported.
// Here we should never suggest something that would include VeryLongCommentThatCouldBeCJKThatCannotBeSplitOnSpaces in the error message.
// MATCH:46 /comment on exported function F should be of the form "F ..."/
