// Test of empty-blocks.

package fixtures

import (
	"fmt"
	"net/http"
)

func f(x int) {} // Must not match

type foo struct{}

func (f foo) f(x *int)  {} // Must not match
func (f *foo) g(y *int) {} // Must not match

func h() {
	go http.ListenAndServe()
	select {} // Must not match
}

func g(f func() bool) {
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

	for _ = range "abc" { // MATCH /this block is empty, you can remove it/

	}

	for _, c := range "abc" { // should not match
		_ = c
	}

	// issue 386, then overwritten by issue 416
	var c = make(chan int)
	for range c { // Must not match
	}

	var s = "a string"
	for range s { // Must not match: we can live with false negatives in this case, as it's not a common pattern and we avoid false positives.
	}

	for range makeChan() { // Must not match
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

func makeChan() chan int {
	ch := make(chan int)
	close(ch)

	return ch
}
