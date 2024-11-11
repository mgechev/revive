// Package foo ...
package foo

import (
	"errors"
	"fmt"
)

// Check for the error strings themselves.

func g(x int) error {
	var err error
	err = fmt.Errorf("This %d is too low", x)     // MATCH /error strings should not be capitalized or end with punctuation or a newline/
	err = fmt.Errorf("XML time")                  // ok
	err = fmt.Errorf("newlines are fun\n")        // MATCH /error strings should not be capitalized or end with punctuation or a newline/
	err = fmt.Errorf("Newlines are really fun\n") // MATCH /error strings should not be capitalized or end with punctuation or a newline/
	err = errors.New(`too much stuff.`)           // MATCH /error strings should not be capitalized or end with punctuation or a newline/
	err = errors.New("This %d is too low", x)     // MATCH /error strings should not be capitalized or end with punctuation or a newline/

	// Non-regression test for issue #610
	d.stack.Push(from)

	return err
}
