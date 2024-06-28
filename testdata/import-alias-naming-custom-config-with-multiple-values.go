package fixtures

import (
	_ "strings"
	bar_foo "strings"  // MATCH /import name (bar_foo) must match the regular expression: ^[a-z][a-z0-9]*$/
	fooBAR "strings"   // MATCH /import name (fooBAR) must match the regular expression: ^[a-z][a-z0-9]*$/
	v1 "strings"       // MATCH /import name (v1) must NOT match the regular expression: ^((v\d+)|(v\d+alpha\d+))$/
	v1alpha1 "strings" // MATCH /import name (v1alpha1) must NOT match the regular expression: ^((v\d+)|(v\d+alpha\d+))$/
	magical "magic/hat"
)

func somefunc() {
	fooBAR.Clone("")
	bar_foo.Clone("")
	v1.Clone("")
	magical.Clone("")
}
