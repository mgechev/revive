package fixtures

import "fmt"

func foo(a, b int) {
	fmt.Printf("single line characters out of limit") // MATCH /line is 105 characters, out of limit 100/
}

// revive:disable-next-line:line-length-limit
// The length of this comment line is over 80 characters, this is bad for readability.

// Warn: the testing framework does not allow to check for failures in comments

func toto() {
	// revive:disable-next-line:line-length-limit
	fmt.Println("This line is way too long. In my opinion, it should be shortened.")
}
