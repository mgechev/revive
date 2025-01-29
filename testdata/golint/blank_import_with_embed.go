// Test that blank import of "embed" is allowed

// Package foo ...
package foo

import (
	_ "embed"
	_ "fmt"
	/* MATCH:8 /a blank import should be only in a main or test package, or have a comment justifying it/ */
)

//go:embed .gitignore
var _ string

