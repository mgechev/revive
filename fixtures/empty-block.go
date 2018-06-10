// Test of empty-blocks.

package fixtures

func f(x int) bool {} // Must not match

func g(f func() bool) string {
	{ // MATCH /this block is empty, you can remove it/
	}

	a := func(e error){} // Must not match

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

}
