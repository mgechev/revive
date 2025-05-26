package fixtures

import (
	magical "magic/hat"
	_ "strings"
	bar_foo "strings"
	fooBAR "strings"
	v1 "strings"       // MATCH /import name (v1) must NOT match the regular expression: ^((v\d+)|(v\d+alpha\d+))$/
	v1alpha1 "strings" // MATCH /import name (v1alpha1) must NOT match the regular expression: ^((v\d+)|(v\d+alpha\d+))$/
)

func somefunc() {
	fooBAR.Clone("")
	bar_foo.Clone("")
	v1.Clone("")
	magical.Clone("")
}
