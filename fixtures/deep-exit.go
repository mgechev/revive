package fixtures

import (
	"log"
	"os"
	"syscall"
	"testing"
)

func foo0() {
	os.Exit(1) // MATCH /calls to os.Exit only in main() or init() functions/
}

func init() {
	log.Fatal("v ...interface{}")
}

func foo() {
	log.Fatalf(1) // MATCH /calls to log.Fatalf only in main() or init() functions/
}

func main() {
	log.Fatalln("v ...interface{}")
}

func bar() {
	log.Fatal(1) // MATCH /calls to log.Fatal only in main() or init() functions/
}

func bar2() {
	bar()
	syscall.Exit(1) // MATCH /calls to syscall.Exit only in main() or init() functions/
}

func TestMain(m *testing.M) {
	// must match because this is not a test file
	os.Exit(m.Run()) // MATCH /calls to os.Exit only in main() or init() functions/
}
