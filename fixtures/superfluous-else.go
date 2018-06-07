// Test of return+else warning.

// Package pkg ...
package pkg

import (
	"fmt"
	"log"
)

func h(f func() bool) string {
	for {
		if ok := f(); ok {
			a := 1
			continue
		} else { // MATCH /if block ends with a continue statement, so drop this else and outdent its block (move short variable declaration to its own line if necessary)/
			return "it's NOT okay!"
		}
	}
}

func i(f func() bool) string {
	for {
		if f() {
			a := 1
			continue
		} else { // MATCH /if block ends with a continue statement, so drop this else and outdent its block/
			log.Printf("non-positive")
		}
	}

	return "ok"
}

func j(f func() bool) string {
	for {
		if f() {
			break
		} else { // MATCH /if block ends with a break statement, so drop this else and outdent its block/
			log.Printf("non-positive")
		}
	}

	return "ok"
}

func k() {
	var a = 10
	/* do loop execution */
LOOP:
	for a < 20 {
		if a == 15 {
			a = a + 1
			goto LOOP
		} else { // MATCH /if block ends with a goto statement, so drop this else and outdent its block/
			fmt.Printf("value of a: %d\n", a)
			a++
		}
	}
}

func j(f func() bool) string {
	for {
		if f() {
			a := 1
			fallthrough
		} else {
			log.Printf("non-positive")
		}
	}

	return "ok"
}
