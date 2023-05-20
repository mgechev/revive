// Test data for the early-return rule with preserveScope option enabled

package fixtures

func fn1() {
	// No initializer, match as normal
	if cond { //   MATCH /if c { ... } else { ... return } can be simplified to if !c { ... return } .../
		fn2()
	} else {
		return
	}
}

func fn2() {
	// Moving the declaration of x here is fine since it goes out of scope either way
	if x := fn1(); x != nil { //   MATCH /if c { ... } else { ... return } can be simplified to if !c { ... return } ... (move short variable declaration to its own line if necessary)/
		fn2()
	} else {
		return
	}
}

func fn3() {
	// Don't want to move the declaration of x here since it stays in scope afterward
	if x := fn1(); x != nil {
		fn2()
	} else {
		return
	}
	x := fn2()
	fn3(x)
}

func fn4() {
	if cond {
		var x = fn2()
		fn3(x)
	} else {
		return
	}
	// Don't want to move the declaration of x here since it stays in scope afterward
	y := fn2()
	fn3(y)
}

func fn5() {
	if cond {
		x := fn2()
		fn3(x)
	} else {
		return
	}
	// Don't want to move the declaration of x here since it stays in scope afterward
	y := fn2()
	fn3(y)
}

func fn6() {
	if cond { //   MATCH /if c { ... } else { ... return } can be simplified to if !c { ... return } .../
		x := fn2()
		fn3(x)
	} else {
		return
	}
	// Moving x here is fine since it goes out of scope anyway
}
