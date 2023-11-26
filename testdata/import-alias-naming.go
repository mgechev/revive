package fixtures

import (
	_ "strings" // _ aliases should be ignored
	. "dotimport" // . aliases should be ignored
	bar_foo "strings" // MATCH /import name (bar_foo) must match the regular expression: ^[a-z][a-z0-9]{0,}$/
	fooBAR "strings"  // MATCH /import name (fooBAR) must match the regular expression: ^[a-z][a-z0-9]{0,}$/
	v1 "strings"
	magical "magic/hat"
)

func somefunc() {
	fooBAR.Clone("")
	bar_foo.Clone("")
	v1.Clone("")
	magical.Clone("")
}
