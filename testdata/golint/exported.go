// Package golint comment
package golint

import "net/http"

//  GolintFoo is a dummy function
func GolintFoo() {} // MATCH /func name will be used as golint.GolintFoo by other packages, and that stutters; consider calling this Foo/

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

// Tests for common method names
//// Should NOT fail for methods
func (_) Error() string                                    { return "" }
func (_) String() string                                   { return "" }
func (_) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
func (_) Read(p []byte) (n int, err error)                 { return 0, nil }
func (_) Write(p []byte) (n int, err error)                { return 0, nil }
func (_) Unwrap(err error) error                           { return nil }

//// Should fail for functions

func Error() string                                    { return "" }     // MATCH /exported function Error should have comment or be unexported/
func String() string                                   { return "" }     // MATCH /exported function String should have comment or be unexported/
func ServeHTTP(w http.ResponseWriter, r *http.Request) {}                // MATCH /exported function ServeHTTP should have comment or be unexported/
func Read(p []byte) (n int, err error)                 { return 0, nil } // MATCH /exported function Read should have comment or be unexported/
func Write(p []byte) (n int, err error)                { return 0, nil } // MATCH /exported function Write should have comment or be unexported/
func Unwrap(err error) error                           { return nil }    // MATCH /exported function Unwrap should have comment or be unexported/

// The following cases are tests for issue 555
