// Test for naming errors.

// Package foo ...
package foo

import (
	"errors"
	"fmt"
)

var unexp = errors.New("some unexported error") // MATCH /error var unexp should have name of the form errFoo/

// Exp ...
var Exp = errors.New("some exported error")

// MATCH:14 /error var Exp should have name of the form ErrFoo/

var (
	e1 = fmt.Errorf("blah %d", 4) // MATCH /error var e1 should have name of the form errFoo/
	// E2 ...
	E2 = fmt.Errorf("blah %d", 5) // MATCH /error var E2 should have name of the form ErrFoo/

	// check there is no false positive for blank identifier.
	// This pattern that can be found in benchmarks and examples should be allowed.
	// The fact that the error variable is not used is out of the scope of this rule.
	_ = errors.New("ok")
)

func f() {

	// the linter does not check local variables, this one is valid
	var whatever = errors.New("ok")
	_ = whatever

	// same as above with a blank identifier
	_ = errors.New("ok")
}
