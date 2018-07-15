package fixtures

import (
	"fmt"
	"log"
	"os"
)

func f() {
	fmt.Println("Hello, playground")
	if true {
		return     // MATCH /unreachable code after this statement/
		os.Exit(2) // ignore
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
