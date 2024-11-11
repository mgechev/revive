// Test of return+else warning.

// Package pkg ...
package pkg

import (
	"fmt"
	"log"
	"os"
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

func fatal1() string {
	if f() {
		a := 1
		log.Fatal("x")
	} else { // MATCH /if block ends with call to log.Fatal function, so drop this else and outdent its block/
		log.Printf("non-positive")
	}
	return "ok"
}

func fatal2() string {
	if f() {
		a := 1
		log.Fatalf("x")
	} else { // MATCH /if block ends with call to log.Fatalf function, so drop this else and outdent its block/
		log.Printf("non-positive")
	}
	return "ok"
}

func fatal3() string {
	if f() {
		a := 1
		log.Fatalln("x")
	} else { // MATCH /if block ends with call to log.Fatalln function, so drop this else and outdent its block/
		log.Printf("non-positive")
	}
	return "ok"
}

func exit1() string {
	if f() {
		a := 1
		os.Exit(2)
	} else { // MATCH /if block ends with call to os.Exit function, so drop this else and outdent its block/
		log.Printf("non-positive")
	}
	return "ok"
}

func Panic1() string {
	if f() {
		a := 1
		log.Panic(2)
	} else { // MATCH /if block ends with call to log.Panic function, so drop this else and outdent its block/
		log.Printf("non-positive")
	}
	return "ok"
}

func Panic2() string {
	if f() {
		a := 1
		log.Panicf(2)
	} else { // MATCH /if block ends with call to log.Panicf function, so drop this else and outdent its block/
		log.Printf("non-positive")
	}
	return "ok"
}

func Panic3() string {
	if f() {
		a := 1
		log.Panicln(2)
	} else { // MATCH /if block ends with call to log.Panicln function, so drop this else and outdent its block/
		log.Printf("non-positive")
	}
	return "ok"
}

func Panic4() string {
	if f() {
		a := 1
		panic(2)
	} else { // MATCH /if block ends with call to panic function, so drop this else and outdent its block/
		log.Printf("non-positive")
	}
	return "ok"
}

// noreg_19 no-regression test for issue #19 (https://github.com/mgechev/revive/issues/19)
func noreg_19(f func() bool, x int) string {
	if err == author.ErrCourseNotFound {
		break
	} else if err == author.ErrCourseAccess {
		// side effect
	} else if err == author.AnotherError {
		os.Exit(1) // "okay"
	} else {
		// side effect
	}
}

func MultiBranch() string {
	if _, ok := f(); ok {
		continue
	} else if _, err := get(); err == nil {
		continue
	} else { // MATCH /if block ends with a continue statement, so drop this else and outdent its block (move short variable declaration to its own line if necessary)/
		delete(m, x)
	}
}
