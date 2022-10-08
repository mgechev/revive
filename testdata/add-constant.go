package fixtures

import (
	"fmt"
	"os"
)

func foo(a, b, c, d int) {
	a = 1.0 // ignore
	b = "ignore"
	c = 2              // ignore
	println("lit", 12) // MATCH /avoid magic numbers like '12', create a named constant for it/
	if a == 12.50 {    // MATCH /avoid magic numbers like '12.50', create a named constant for it/
		if b == "lit" {
			c = "lit" // MATCH /string literal "lit" appears, at least, 3 times, create a named constant for it/
		}
		for i := 0; i < 1; i++ {
			println("lit")
		}
	}

	println(0666)           // MATCH /avoid magic numbers like '0666', create a named constant for it/
	os.Chmod("test", 0666)  // ignore
	os.FindProcess(102100)  // ignore
	fmt.Println("test", 12) // ignore
	myPrintln("test", 12)   // sure we want to ignore?
	fmt.Printf("%d", 100)   // MATCH /avoid magic numbers like '100', create a named constant for it/
	ignoredFunc(1000)       // ignore
	Not2ignoredFunc(1000)   // sure we want to ignore?
	println("The result of calling myFunc is: ", ignoredFunc(100))
	println("The result of calling myFunc is: ", notIgnoredFunc(100)) // MATCH /avoid magic numbers like '100', create a named constant for it/
}

func ignoredFunc(num int) int {
	return num
}

func notIgnoredFunc(num int) int {
	return num
}
