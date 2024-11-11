// Test of return+else warning.

// Package pkg ...
package pkg

import "log"

func f(x int) bool {
	if x > 0 {
		return true
	} else { // MATCH /if block ends with a return statement, so drop this else and outdent its block/
		log.Printf("non-positive x: %d", x)
	}
	return false
}

func g(f func() bool) string {
	if ok := f(); ok {
		return "it's okay"
	} else { // MATCH /if block ends with a return statement, so drop this else and outdent its block (move short variable declaration to its own line if necessary)/
		return "it's NOT okay!"
	}
}

func h(f func() bool, x int) string {
	if err == author.ErrCourseNotFound {
		return
	} else if err == author.ErrCourseAccess {
		// side effect
	} else if err == author.AnotherError {
		return "okay"
	} else {
		if ok := f(); ok {
			return "it's okay"
		} else { // MATCH /if block ends with a return statement, so drop this else and outdent its block (move short variable declaration to its own line if necessary)/
			return "it's NOT okay!"
		}
	}
}

func i() string {
	if err == author.ErrCourseNotFound {
		return "not found"
	} else if err == author.AnotherError {
		return "something else"
	} else { // MATCH /if block ends with a return statement, so drop this else and outdent its block/
		return "okay"
	}
}
