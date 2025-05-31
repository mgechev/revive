package fixtures

func enforceSwitchStyle3() {

	switch expression {
	case condition:
	default:
	}

	switch expression {
	default:
	case condition:
	}

	switch expression {
	case condition:
	}
}
