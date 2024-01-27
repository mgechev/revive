package fixtures

func mcn() {
	if true {
		if true {
			if true { // MATCH /control flow nesting exceeds 2/

			}
		}
	} else {
		if true {
			if true { // MATCH /control flow nesting exceeds 2/
				if true {

				}
			}
		}
	}

	for {
		if true {
			for { // MATCH /control flow nesting exceeds 2/
			}
		}
	}

	switch {
	case false:
		if true {

		}
	case true:
		if true {
			for { // MATCH /control flow nesting exceeds 2/
			}
		}
	default:
	}
}
