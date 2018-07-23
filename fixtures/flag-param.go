package fixtures

func foo(a bool, b int) { // MATCH /parameter 'a' seems to be a control flag, avoid control coupling/
	if a {

	}
}

func foo(a bool, b int) {
	str := mystruct{a, b}
}
