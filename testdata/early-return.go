// Test of empty-blocks.

package fixtures

func earlyRet() bool {
	if cond { //   MATCH /if c {...} else {... return } can be simplified to if !c { ... return } .../
		println()
		println()
		println()
	} else {
		return false
	}

	if cond { //MATCH /if c {...} else {... return } can be simplified to if !c { ... return } .../
		println()
	} else {
		return false
	}

	if cond { //MATCH /if c { } else {... return} can be simplified to if !c { ... return }/
	} else {
		return false
	}

	if cond {
		println()
	} else if cond { //MATCH /if c { } else {... return} can be simplified to if !c { ... return }/
	} else {
		return false
	}

	if cond {
		println()
	} else if cond { //MATCH /if c {...} else {... return } can be simplified to if !c { ... return } .../
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

	if cond { //MATCH /if c {...} else {... return } can be simplified to if !c { ... return } .../
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
}
