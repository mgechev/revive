package fixtures

func unnecessaryConditional() bool {
	if cond {
		return true
	} else {
		return false
	}

	if cond {
		id = true
	} else {
		id = false
	}

	return false
}
