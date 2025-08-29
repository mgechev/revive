// Test that blank imports in library packages are flagged.

// Package foo ...
package fixtures

func foo() (s1, s2 string) { // ok
	return "foo", "bar"
}

func bar() (s1 string, err error) { // ok
	return "foo", nil
}

func batz() (string, string, string) { // MATCH /when a function has more than two return values, only one should be named/
	return "a", "b", "c"
}

func qux() (x, y, z string) { // ok
	return "a", "b", "c"
}
