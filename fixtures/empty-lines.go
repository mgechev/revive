// Test of empty-lines.

package fixtures

func f1(x *int) bool { // MATCH /extra empty line at the start of a block/

	return x > 2
}

func f2(x *int) bool {
	return x > 2 // MATCH /extra empty line at the end of a block/

}

func f3(x *int) bool { // MATCH /extra empty line at the start of a block/

	return x > 2 // MATCH /extra empty line at the end of a block/

}

func f4(x *int) bool {
	// This is fine.
	return x > 2
}

func f5(x *int) bool { // MATCH /extra empty line at the start of a block/

	// This is _not_ fine.
	return x > 2
}

func f6(x *int) bool {
	return x > 2
	// This is fine.
}

func f7(x *int) bool {
	return x > 2 // MATCH /extra empty line at the end of a block/
	// This is _not_ fine.

}

func f8(*int) bool {
	if x > 2 { // MATCH /extra empty line at the start of a block/

		return true
	}

	return false
}

func f9(*int) bool {
	if x > 2 {
		return true // MATCH /extra empty line at the end of a block/

	}

	return false
}
