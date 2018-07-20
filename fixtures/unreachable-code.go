package fixtures

import (
	"fmt"
	"log"
	"os"
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
