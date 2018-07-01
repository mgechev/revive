package fixtures

import (
	"fmt"
)

func foo(a, b, c, d int) {
	fmt.Print("bad")        // MATCH /do not call fmt.Print, use logger/
	fmt.Println("bad")      // MATCH /do not call fmt.Println, use logger/
	fmt.Printf("%s", "bad") // MATCH /do not call fmt.Printf, use logger/
	fmt.Sprint("ss")
}
