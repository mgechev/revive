package fixtures

func foo() (a, b, c, d) { // MATCH /maximum number of return results per function exceeded; max 3 but got 4/
	var a, b, c, d int
}

func bar(a, b int) {

}

func baz(a string, b int) {

}

func qux() (string, string, int, string, int) { // MATCH /maximum number of return results per function exceeded; max 3 but got 5/

}
