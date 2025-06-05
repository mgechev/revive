package main_test

// this is a test file in the main_test package
//
// such files SHOULD NOT be linted by unexported_return rule
// because symbols defined in main package cannot be imported in other packages

type foo struct{}

func NewFoo() foo {
	return foo{}
}
