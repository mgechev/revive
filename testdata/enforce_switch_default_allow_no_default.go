package fixtures

func enforceSwitchDefault() {

	switch expression {
	case condition:
	default:
	}

	switch expression {
	default: // MATCH /default case clause must be the last one/
	case condition:
	}

	switch expression {
	case condition:
	}
}
