package fixtures

func enforceSwitchStyle3() {

	switch expression { // skipt tagged switch
	case value:
	default:
	}

	switch {
	case a > 0, a < 0:
	case a == 0:
	case a < 0: // MATCH /case clause at line 11 has the same condition/
	default:
	}

	switch {
	case a > 0, a < 0, a > 0: // MATCH /case clause at line 18 has the same condition/
	case a == 0:
	case a < 0: // MATCH /case clause at line 18 has the same condition/
	default:
	}

	switch something {
	case 1:
		switch {
		case a > 0, a < 0, a > 0: // MATCH /case clause at line 27 has the same condition/
		case a == 0:
		}
	default:
	}

	switch {
	case a == 0:
		switch {
		case a > 0, a < 0, a > 0: // MATCH /case clause at line 36 has the same condition/
		case a == 0:
		}
	default:
	}
}
