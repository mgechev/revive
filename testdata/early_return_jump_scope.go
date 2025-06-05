// Test data for the early-return rule with allowJump option enabled

package fixtures

import (
	"os"
)

func fn1() {
	if cond { //MATCH /if c { ... } can be rewritten if !c { return } ... to reduce nesting/
		println()
		println()
		println()
	}
}

func fn3() {
	if a() {
		println()
		os.Exit(1)
	}
}

func fn4() {
	// No initializer, match as normal
	if cond { //   MATCH /if c { ... } else { ... return } can be simplified to if !c { ... return } .../
		fn2()
	} else {
		return
	}
}
