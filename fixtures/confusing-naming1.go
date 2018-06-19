// Test of confusing-naming rule.

// Package pkg ...
package pkg

type foo struct{}

func (t foo) aFoo() {
	return
}

func (t *foo) AFoo() { // MATCH /Method 'AFoo' differs only by capitalization to other method of 'foo'/
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

func AGlobal() { // MATCH /Function 'AGlobal' differs only by capitalization to other function in the same package/
}

func ABar() { // Should not warn

}

func aFoo() { // Should not warn

}

func (t foo) ABar() { // Should not warn
	return
}

func (t bar) ABar() { // MATCH /Method 'ABar' differs only by capitalization to other method of 'bar'/
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
