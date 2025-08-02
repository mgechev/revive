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
}
