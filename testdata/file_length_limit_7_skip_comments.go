package fixtures

import "fmt"

// Foo is a function.
func Foo(a, b int) {
	// This
	/* is
	a
	*/
	// a comment.
	fmt.Println("Hello, world!")
	/*
		This is
		multiline
		comment.
	*/
}

// MATCH /file length is 8 lines, which exceeds the limit of 7/
