package fixtures

func enforceSwitchDefault() {

	switch expression {
	case condition:
	default:
	}

	switch expression {
	default:
	case condition:
	}

	switch expression { // MATCH /switch must have a default case clause/
	case condition:
	}
}
