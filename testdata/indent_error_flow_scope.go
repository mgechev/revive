package pkg

func fn4() {
	if cond {
		var x = fn2()
		fn3(x)
		return
	} else {
		y := fn2()
		fn3(y)
	}
	// Don't want to move the declaration of x here since it stays in scope afterward
	y := fn2()
	fn3(y)
}
