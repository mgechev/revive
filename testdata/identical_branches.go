package fixtures

func identicalBranches() {
	if true { // MATCH /both branches of the if are identical/

	} else {

	}

	if true {

	}

	if true {
		print()
	} else {
	}

	if true { // MATCH /both branches of the if are identical/
		print()
	} else {
		print()
	}

	if true {
		if true { // MATCH /both branches of the if are identical/
			print()
		} else {
			print()
		}
	} else {
		println()
	}

	if true {
		println("something")
	} else {
		println("else")
	}

	if true { // MATCH /this if...else if chain has identical branches (lines [41 49])/
		print("something")
	} else if true {
		print("something else")
	} else if true {
		print("other thing")
	} else if false {
		println()
	} else {
		print("something")
	}

	if true { // MATCH /this if...else if chain has identical branches (lines [53 59])/
		print("something")
	} else if true {
		print("something else")
	} else if true {
		print("other thing")
	} else if false {
		print("something")
	} else {
		println()
	}

	if true {
		print("something")
	} else if true {
		print("something else")
		if true { // MATCH /this if...else if chain has identical branches (lines [69 71])/
			print("something")
		} else if false {
			print("something")
		} else {
			println()
		}
	}

	// Should not warn because even branches are identical, the err variable is not
	if err := something(); err != nil {
		println(err)
	} else if err := somethingElse(); err != nil {
		println(err)
	}
}
