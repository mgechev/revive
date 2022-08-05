package fixtures

import (
	pkgErrors "github.com/pkg/errors"
)

// Check for the error strings themselves.

func errorsStrings(x int) error {
	var err error
	return pkgErrors.Wrap(err, "This %d is too low") // MATCH /error strings should not be capitalized or end with punctuation or a newline/
}
