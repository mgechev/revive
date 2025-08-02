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

	if true { // MATCH /"if...else if" chain with identical branches (lines 41 and 49)/
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

	if true { // MATCH /"if...else if" chain with identical branches (lines 53 and 59)/
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
		if true { // MATCH /"if...else if" chain with identical branches (lines 69 and 71)/
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
	// MATCH:86 /"if...else if" chain with identical branches (lines 86 and 90)/
	// MATCH:86 /"if...else if" chain with identical branches (lines 88 and 92)/

	if createFile() { // json:{"MATCH": "\"if...else if\" chain with identical branches (lines 98 and 102)","Confidence": 0.8}
		doSomething()
	} else if !delete() {
		return new("cannot delete file")
	} else if createFile() {
		doSomething()
	} else {
		return new("file error")
	}

	// Test confidence is reset
	if a { // json:{"MATCH": "\"if...else if\" chain with identical branches (lines 109 and 111)","Confidence": 1}
		foo()
	} else if b {
		foo()
	} else {
		bar()
	}

	switch a { // MATCH /"switch" with identical branches (lines 119 and 123)/
	// expected values
	case 1:
		foo()
	case 2:
		bar()
	case 3:
		foo()
	default:
		return newError("blah")
	}

	// MATCH:131 /"switch" with identical branches (lines 133 and 137)/
	// MATCH:131 /"switch" with identical branches (lines 135 and 139)/
	switch a {
	// expected values
	case 1:
		foo()
	case 2:
		bar()
	case 3:
		foo()
	default:
		bar()
	}

	switch a { // MATCH /"switch" with identical branches (lines 145 and 147)/
	// expected values
	case 1:
		foo()
	case 3:
		foo()
	default:
		if true { // MATCH /"if...else if" chain with identical branches (lines 150 and 152)/
			something()
		} else if true {
			something()
		} else {
			if true { // MATCH /both branches of the if are identical/
				print("identical")
			} else {
				print("identical")
			}
		}
	}

	// Skip untagged switch
	switch {
	case a > b:
		foo()
	default:
		foo()
	}

	// Do not warn on fallthrough
	switch a {
	case 1:
		foo()
		fallthrough
	case 2:
		fallthrough
	case 3:
		foo()
	case 4:
		fallthrough
	default:
		bar()
	}
}
