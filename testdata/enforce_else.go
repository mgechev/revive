package fixtures

func enforceElse() {
	if true {
		// something
	} else if true { // MATCH /"if ... else if" constructs should end with "else" clauses/
		// something else
	}

	if true {
		return
	} else if true {
		return
	}

	if true {
		break
	} else if true {
		break
	}

	if true {
		return
	} else if true { // MATCH /"if ... else if" constructs should end with "else" clauses/
		// something else
	}

	if true {
		// something
	} else if true { // MATCH /"if ... else if" constructs should end with "else" clauses/
		return
	}

	if true {
		// something
	} else if true { // MATCH /"if ... else if" constructs should end with "else" clauses/
		if true {
			// something
		} else if true { // MATCH /"if ... else if" constructs should end with "else" clauses/
			return
		}
	}

	if true {
		if true {
			// something
		} else if true { // MATCH /"if ... else if" constructs should end with "else" clauses/
			return
		}
	} else if true { // MATCH /"if ... else if" constructs should end with "else" clauses/
		// something
	}
}
