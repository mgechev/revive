// Test data for the superfluous-else rule with preserveScope option enabled

package fixtures

func fn1() {
	for {
		// No initializer, match as normal
		if cond {
			continue
		} else { // MATCH /if block ends with a continue statement, so drop this else and outdent its block/
			fn2()
		}
	}
}

func fn2() {
	for {
		// Moving the declaration of x here is fine since it goes out of scope either way
		if x := fn1(); x != nil {
			continue
		} else { // MATCH /if block ends with a continue statement, so drop this else and outdent its block (move short variable declaration to its own line if necessary)/
			fn2()
		}
	}
}

func fn3() {
	for {
		// Don't want to move the declaration of x here since it stays in scope afterward
		if x := fn1(); x != nil {
			continue
		} else {
			fn2()
		}
		x := fn2()
		fn3(x)
	}
}

func fn4() {
	for {
		if cond {
			continue
		} else {
			x := fn1()
			fn2(x)
		}
		// Don't want to move the declaration of x here since it stays in scope afterward
		y := fn2()
		fn3(y)
	}
}

func fn4() {
	for {
		if cond {
			continue
		} else { // MATCH /if block ends with a continue statement, so drop this else and outdent its block/
			x := fn1()
			fn2(x)
		}
		// Moving x here is fine since it goes out of scope anyway
	}
}
