package fixtures

import (
	"errors"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

// Testable package level example
func Example() {
	log.Fatal(errors.New("example"))
}

// Testable function example
func ExampleFoo() {
	log.Fatal(errors.New("example"))
}

// Testable method example
func ExampleBar_Qux() {
	log.Fatal(errors.New("example"))
}

// Not an example because it has an argument
func ExampleBar(int) {
	log.Fatal(errors.New("example")) // MATCH /calls to log.Fatal only in main() or init() functions/
}
