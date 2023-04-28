// Test of empty-blocks.

package fixtures

import "net/http"

func f(x int) {} // Must not match

type foo struct{}

func (f foo) f(x *int)  {} // Must not match
func (f *foo) g(y *int) {} // Must not match

func h() {
	go http.ListenAndServe()
	select {} // Must not match
}

var (
	pkgDeclChan chan int
	pkgDefChan  = make(chan int)
)

func g(f func() bool, paramChan chan int) {
	{ // MATCH /this block is empty, you can remove it/
	}

	_ = func(e error) {} // Must not match

	if ok := f(); ok { // MATCH /this block is empty, you can remove it/
		// only a comment
	} else {
		println("it's NOT empty!")
	}

	if ok := f(); ok {
		println("it's NOT empty!")
	} else { // MATCH /this block is empty, you can remove it/

	}

	for i := 0; i < 10; i++ { // MATCH /this block is empty, you can remove it/

	}

	for { // MATCH /this block is empty, you can remove it/

	}

	for true { // MATCH /this block is empty, you can remove it/

	}

	// issue 386, then overwritten by issue 416
	var c = make(chan int)
	for range c { // Must not match
	}

	assignChan := make(chan int)
	for range assignChan { // Must not match
	}

	_, reAssignChan, _ := "foo", assignChan, "bar"
	for range reAssignChan { // Must not match
	}

	_, secondChan := twoValues()
	for range secondChan { // Must not match
	}

	for range paramChan { // Must not match
	}

	reAssignParamChan := paramChan
	for range reAssignParamChan { // Must not match
	}

	for range callChan() { // Must not match
	}

	for range func() <-chan int {
		c := make(chan int)
		return c
	}() { // Must not match
	}

	for range pkgDeclChan { // Must not match
	}

	for range pkgDefChan { // Must not match
	}

	var s = "a string"
	for range s { // MATCH /this block is empty, you can remove it/
	}

	select {
	case _, ok := <-c:
		if ok { // MATCH /this block is empty, you can remove it/
		}
	}

	// issue 810
	next := 0
	iter := func(v *int) bool {
		*v = next
		next++
		fmt.Println(*v)
		return next < 10
	}

	z := 0
	for iter(&z) { // Must not match
	}

	for process() { // Must not match
	}

	var it iterator
	for it.next() { // Must not match
	}
}

func process() bool {
	return false
}

type iterator struct{}

func (it *iterator) next() bool {
	return false
}

func callChan() <-chan int {
	return make(chan int)
}

func twoValues() (interface{}, <-chan int) {
	return "test", make(chan int)
}
