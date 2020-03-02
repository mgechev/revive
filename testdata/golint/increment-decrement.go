// Test for use of x++ and x--.

// Package pkg ...
package pkg

func addOne(x int) int {
	x += 1 // MATCH /should replace x += 1 with x++/
	return x
}

func subOneInLoop(y int) {
	for ; y > 0; y -= 1 { // MATCH /should replace y -= 1 with y--/
	}
}
