package fixtures

import (
	. "dotimport" // . aliases should be ignored
	magical "magic/hat"
	_ "strings"       // _ aliases should be ignored
	bar_foo "strings" // MATCH /import name (bar_foo) must match the regular expression: ^[a-z][a-z0-9]{0,}$/
	fooBAR "strings"  // MATCH /import name (fooBAR) must match the regular expression: ^[a-z][a-z0-9]{0,}$/
	v1 "strings"
)

func somefunc() {
	fooBAR.Clone("")
	bar_foo.Clone("")
	v1.Clone("")
	magical.Clone("")
}
