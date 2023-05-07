// Test of empty-blocks.

package fixtures

import "os"

func earlyRet() bool {
	if cond { //   MATCH /if c { ... } else { ... return } can be simplified to if !c { ... return } .../
		println()
		println()
		println()
	} else {
		return false
	}

	if cond { //MATCH /if c { ... } else { ... return } can be simplified to if !c { ... return } .../
		println()
	} else {
		return false
	}

	if cond { //MATCH /if c { } else { ... return } can be simplified to if !c { ... return }/
	} else {
		return false
	}

	if cond {
		println()
	} else if cond { //MATCH /if c { } else { ... return } can be simplified to if !c { ... return }/
	} else {
		return false
	}

	// the first branch does not return, so we can't reduce nesting here
	if cond {
		println()
	} else if cond {
		println()
	} else {
		return false
	}

	// Case already covered by golint
	if cond {
		return true
	} else {
		return false
	}

	if cond { //MATCH /if c { ... } else { ... return } can be simplified to if !c { ... return } .../
		println()
		println()
		println()
	} else {
		return false
	}

	if cond {
		println()
		println()
		println()
	} else {
		println()
	}

	if cond {
		if cond { //MATCH /if c { ... } else { ... return } can be simplified to if !c { ... return } .../
			println()
		} else {
			return false
		}
	}

	if cond {
		println()
	} else {
		if cond { //MATCH /if c { ... } else { ... return } can be simplified to if !c { ... return } .../
			println()
		} else {
			return false
		}
	}

	if cond {
		println()
	} else if cond {
		println()
	} else {
		if cond { //MATCH /if c { ... } else { ... return } can be simplified to if !c { ... return } .../
			println()
		} else {
			return false
		}
	}

	for {
		if cond { //MATCH /if c { ... } else { ... continue } can be simplified to if !c { ... continue } .../
			println()
		} else {
			continue
		}
	}

	for {
		if cond { //MATCH /if c { ... } else { ... break } can be simplified to if !c { ... break } .../
			println()
		} else {
			break
		}
	}

	if cond { //MATCH /if c { ... } else { ... panic() } can be simplified to if !c { ... panic() } .../
		println()
	} else {
		panic("!")
	}

	if cond { //MATCH /if c { ... } else { ... goto } can be simplified to if !c { ... goto } .../
		println()
	} else {
		goto X
	}

	if x, ok := foo(); ok { //MATCH /if c { ... } else { ... return } can be simplified to if !c { ... return } ... (move short variable declaration to its own line if necessary)/
		println(x)
	} else {
		return false
	}

	if cond { //MATCH /if c { ... } else { ... os.Exit() } can be simplified to if !c { ... os.Exit() } .../
		println()
	} else {
		os.Exit(0)
	}
}
