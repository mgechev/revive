package fixtures

func foo(a, b, c, d int) {
	if bar == true { // MATCH /omit comparison with boolean constants/

	}
	for f() || false != yes { // MATCH /omit comparison with boolean constants/

	}
}

func bar() {
	a := 1
loop:
	for {
		switch a {
		case 1:
			a++
			println("one")
			break // MATCH /omit unnecessary break at the end of case clause/
		case 2:
			println("two")
			break loop
		default:
			println("default")
		}
	}

	return // MATCH /omit unnecessary return statement/
}
