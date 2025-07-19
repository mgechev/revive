package fixtures

import "slices"

func one(a int) {
	a, b := 1, 2 // MATCH /parameter 'a' seems to be modified/
	a++          // MATCH /parameter 'a' seems to be modified/
}

func two(b, c float32) {
	if c > 0.0 {
		b = 1 // MATCH /parameter 'b' seems to be modified/
	}
}

type foo struct {
	a string
}

func three(s *foo) {
	s.a = "foooooo"
}

// non regression test for issue 355
func issue355(_ *foo) {
	_ = "foooooo"
}

func testSlicesDeleteAssigned(s []int) {
	s = slices.Delete(s, 0, 1)
	// MATCH:30 /parameter 's' seems to be modified/
	// MATCH:30 /parameter 's' seems to be modified by 'slices.Delete'/
	s = slices.DeleteFunc(s, func(e string) bool { return true })
	// MATCH:33 /parameter 's' seems to be modified/
	// MATCH:33 /parameter 's' seems to be modified by 'slices.DeleteFunc'/
	_ = slices.Delete(s, 0, 1)                     // MATCH /parameter 's' seems to be modified by 'slices.Delete'/
	_ = slices.DeleteFunc(s, func(e string) bool { // MATCH /parameter 's' seems to be modified by 'slices.DeleteFunc'/
		return true
	})
	s, b := slices.Delete(s, 0, 1), 2
	// MATCH:40 /parameter 's' seems to be modified/
	// MATCH:40 /parameter 's' seems to be modified by 'slices.Delete'/
	s := slices.Clone(s)
	// MATCH:43 /parameter 's' seems to be modified/
}

func testSlicesDeleteCloned(s []int) {
	s2 := slices.Clone(s)
	s2 = slices.Delete(s2, 0, 1)
	_ = slices.Delete(slices.Clone(s), 0, 1)
	_ = slices.DeleteFunc(slices.Clone(s), func(e string) bool {
		return e == "test"
	})
}

func testSlicesDeleteCopied(s []int) {
	c := make([]int, len(s))
	copy(c, s)
	c = slices.Delete(c, 0, 1)
}

func testMultipleParams(a, b, s []int) {
	a = slices.Delete(a, 0, 1)
	// MATCH:63 /parameter 'a' seems to be modified/
	// MATCH:63 /parameter 'a' seems to be modified by 'slices.Delete'/
	b = slices.Delete(b, 1, 2)
	// MATCH:66 /parameter 'b' seems to be modified/
	// MATCH:66 /parameter 'b' seems to be modified by 'slices.Delete'/
	s = []int{1, 2, 3} // MATCH /parameter 's' seems to be modified/
	s = slices.Delete(s, 0, 1)
	// MATCH:70 /parameter 's' seems to be modified/
	// MATCH:70 /parameter 's' seems to be modified by 'slices.Delete'/
}

func testAssignToNewVar(s []int) {
	newSlice := s
	newSlice = slices.Delete(newSlice, 0, 1)
}
