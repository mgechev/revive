// Package golint comment
package golint

type (
	// O is a shortcut (alias) for map[string]interface{}, e.g. a JSON object.
	O = map[string]interface{}

	// A is shortcut for []O.
	A = []O

	// This Person type is simple
	Person = map[string]interface{}
)

type Foo struct{} // MATCH /exported type Foo should have comment or be unexported/

// The following cases are no-regression tests for issue 229

/* Bar something */
type Bar struct{}

/* Toto something */
func Toto() {}

/* FirstLetter something */
const FirstLetter = "A"

/*Bar2 something */
type Bar2 struct{}

/*Toto2 something */
func Toto2() {}

/*SecondLetter something */
const SecondLetter = "B"
