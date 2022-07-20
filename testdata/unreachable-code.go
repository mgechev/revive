package fixtures

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func foo() int {
	log.Fatalf("%s", "About to fail") // ignore
	return 0                          // MATCH /unreachable code after this statement/
	return 1
	Println("unreachable")
}

func f() {
	fmt.Println("Hello, playground")
	if true {
		return // MATCH /unreachable code after this statement/
		Println("unreachable")
		os.Exit(2) // ignore
		Println("also unreachable")
	}
	return // MATCH /unreachable code after this statement/
	fmt.Println("Bye, playground")
}

func g() {
	fmt.Println("Hello, playground")
	if true {
		return // ignore if next stmt is labeled
	label:
		os.Exit(2) // ignore
	}

	fmt.Println("Bye, playground")
}

func TestA(t *testing.T) {
	tests := make([]int, 100)
	for i := range tests {
		println("i: ", i)
		if i == 0 {
			t.Fatal("i == 0") // MATCH /unreachable code after this statement/
			println("unreachable")
			continue
		}
		if i == 1 {
			t.Fatalf("i:%d", i) // MATCH /unreachable code after this statement/
			println("unreachable")
			continue
		}
		if i == 2 {
			t.FailNow() // MATCH /unreachable code after this statement/
			println("unreachable")
			continue
		}
	}
}
