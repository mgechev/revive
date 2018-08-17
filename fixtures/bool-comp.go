package fixtures

func foo(a, b, c, d int) {
	if bar == true { // MATCH /omit comparison with boolean constants/

	}
	for f() || false != yes { // MATCH /omit comparison with boolean constants/

	}

	if a == a { // MATCH /operands are the same on both sides of the binary expression/

	}

	if b > b { // MATCH /operands are the same on both sides of the binary expression/

	}

	if !b == b {

	}
}
