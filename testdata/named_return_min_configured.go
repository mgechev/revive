// Package foo ...
package fixtures

func batz() string { // MATCH /when a function has more than 1 return values, only one should be named/
	return "a", "b", "c"
}

func qux() (x string) { // ok
	return "a", "b", "c"
}
