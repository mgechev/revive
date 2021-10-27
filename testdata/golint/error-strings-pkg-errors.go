// Package foo ...
package foo

import (
	"github.com/pkg/errors"
)

// Check for the error strings themselves.

func errorsStrings(x int) error {
	var err error
	err = errors.Wrap(err, "This %d is too low")            // MATCH /error strings should not be capitalized or end with punctuation or a newline/
	err = errors.New("This %d is too low")                  // MATCH /error strings should not be capitalized or end with punctuation or a newline/
	err = errors.Wrapf(err, "This %d is too low", x)        // MATCH /error strings should not be capitalized or end with punctuation or a newline/
	err = errors.WithMessage(err, "This %d is too low")     // MATCH /error strings should not be capitalized or end with punctuation or a newline/
	err = errors.WithMessagef(err, "This %d is too low", x) // MATCH /error strings should not be capitalized or end with punctuation or a newline/
	return err
}
