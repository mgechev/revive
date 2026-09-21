package fixtures

func enforceSwitchStyle3() {

	switch expression {
	case condition:
	default:
	}

	switch expression {
	default: // MATCH /default case clause must be the last one/
	case condition:
	}

	switch expression { // MATCH /switch must have a default case clause/
	case condition:
	}

	// Must not fail when all branches jump
	switch expression {
	case condition:
		break
	case condition:
		print()
		return
	}

	// Type switch: default last (ok)
	switch v := expression.(type) {
	case int:
		print(v)
	default:
	}

	// Type switch: default not last
	switch v := expression.(type) {
	default: // MATCH /default case clause must be the last one/
		print(v)
	case int:
	}

	// Type switch: no default
	switch v := expression.(type) { // MATCH /switch must have a default case clause/
	case int:
		print(v)
	}

	// Type switch: must not fail when all branches jump
	switch v := expression.(type) {
	case int:
		print(v)
		break
	case string:
		print(v)
		return
	}
}
