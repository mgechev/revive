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

	if true { // MATCH /"if...else if" chain with identical branches (lines [41 49])/
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

	if true { // MATCH /"if...else if" chain with identical branches (lines [53 59])/
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
		if true { // MATCH /"if...else if" chain with identical branches (lines [69 71])/
			print("something")
		} else if false {
			print("something")
		} else {
			println()
		}
	}

	// Should not warn because even if branches are identical, the err variable is not
	if err := something(); err != nil {
		println(err)
	} else if err := somethingElse(); err != nil {
		println(err)
	}

	// Multiple identical pair of branches
	if a {
		foo()
	} else if b {
		bar()
	} else if c {
		foo()
	} else if d {
		bar()
	}
	// MATCH:86 /"if...else if" chain with identical branches (lines [86 90])/
	// MATCH:86 /"if...else if" chain with identical branches (lines [88 92])/

	if createFile() { // json:{"MATCH": "\"if...else if\" chain with identical branches (lines [98 102])","Confidence": 0.8}
		doSomething()
	} else if !delete() {
		return new("cannot delete file")
	} else if createFile() {
		doSomething()
	} else {
		return new("file error")
	}

	// Test confidence is reset
	if a { // json:{"MATCH": "\"if...else if\" chain with identical branches (lines [109 111])","Confidence": 1}
		foo()
	} else if b {
		foo()
	} else {
		bar()
	}
}
