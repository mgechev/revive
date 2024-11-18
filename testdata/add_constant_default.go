package fixtures

func foo() {
	a = "ignore"
	b = "ignore"

	c = "match"
	d = "match"
	e = "match" // MATCH /string literal "match" appears, at least, 3 times, create a named constant for it/

	f = 5 // MATCH /avoid magic numbers like '5', create a named constant for it/
}
