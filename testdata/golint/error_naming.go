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
)

func f() {
	var whatever = errors.New("ok") // ok
	_ = whatever
}
