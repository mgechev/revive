package fixtures

func fooFlagP(a bool, b int) { // MATCH /parameter 'a' seems to be a control flag, avoid control coupling/
	if a {

	}
}

func barFlagP(a bool, b int) {
	str := mystruct{a, b}
}

// issue #1211
func bazFlagP(a int, b bool) {
	lBool := true
	if lBool {
		// do something
	}
}
