package foo_test

// this is a test file in a _test package
//
// such files SHOULD NOT be linted by unexported_return rule
// because symbols defined in test files cannot be used in other packages

type foo struct{}

func NewFoo() foo {
	return foo{}
}
