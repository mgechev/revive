// Test of confusing-naming rule.
// Package pkg ...
package pkg

type foo struct{}

func (t foo) aFoo() {
	return
}

func (t *foo) AFoo() { // MATCH /Method 'AFoo' differs only by capitalization to method 'aFoo' in the same source file/
	return
}

type bar struct{}

func (t *bar) aBar() {
	return
}

func (t *bar) aFoo() { // Should not warn
	return
}

func aGlobal() {

}

func AGlobal() { // MATCH /Method 'AGlobal' differs only by capitalization to function 'aGlobal' in the same source file/
}

func ABar() { // Should not warn

}

func aFoo() { // Should not warn

}

func (t foo) ABar() { // Should not warn
	return
}

func (t bar) ABar() { // MATCH /Method 'ABar' differs only by capitalization to method 'aBar' in the same source file/
	return
}

func x() {}

type tFoo struct {
	asd      string
	aSd      int  // MATCH /Field 'aSd' differs only by capitalization to other field in the struct type tFoo/
	qwe, asD bool // MATCH /Field 'asD' differs only by capitalization to other field in the struct type tFoo/
	zxc      float32
}

type tBar struct {
	asd string
	qwe bool
	zxc float32
}

// issue #864
type x[T any] struct{}

func (x[T]) method() {
}

type y[T any] struct{}

func (y[T]) method() {
}

// issue #982
type a[T any] struct{}

func (x *a[T]) method() {
}

type b[T any] struct{}

func (x *b[T]) method() {
}
