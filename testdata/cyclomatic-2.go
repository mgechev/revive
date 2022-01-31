// Test of cyclomatic complexity.

// Package pkg ...
package pkg

import "log"

func f(x int) bool { // MATCH /function f has cyclomatic complexity 4 (> max enabled 3)/
	if x > 0 && true || false {
		return true
	} else {
		log.Printf("non-positive x: %d", x)
	}
	return false
}

func g(f func() bool) string {
	if ok := f(); ok {
		return "it's okay"
	} else {
		return "it's NOT okay!"
	}
}
