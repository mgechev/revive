package fixtures

import (
	_ "strings"
	bar_foo "strings" // MATCH /import name (bar_foo) must match the regular expression: ^[a-z]+$/
	fooBAR "strings"  // MATCH /import name (fooBAR) must match the regular expression: ^[a-z]+$/
	v1 "strings"      // MATCH /import name (v1) must match the regular expression: ^[a-z]+$/
	magical "magic/hat"
)

func somefunc() {
	fooBAR.Clone("")
	bar_foo.Clone("")
	v1.Clone("")
	magical.Clone("")
}
