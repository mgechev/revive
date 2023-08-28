package fixtures

import (
	_ "strings"
	bar_foo "strings" // MATCH /import name (bar_foo) must match the regular expression: ^[a-z][a-z0-9]$/
	fooBAR "strings"  // MATCH /import name (fooBAR) must match the regular expression: ^[a-z][a-z0-9]$/
	v1 "strings"
)

func somefunc() {
	fooBAR.Clone("")
	bar_foo.Clone("")
	v1.Clone("")
}
