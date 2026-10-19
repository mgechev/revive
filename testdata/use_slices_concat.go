package fixtures

func useSlicesConcat(s1, s2, s3 []int, magic []byte) {
	_ = append(append([]int{}, s1...), s2...)
	_ = append(append(append([]int{}, s1...), s2...), s3...)
	_ = append(append([]byte{}, magic...), []byte("a payload")...)

	all := append([]int{}, s1...)
	all = append(all, s2...)
	_ = all

	switch len(s1) {
	case 0:
		inCase := append([]int{}, s1...)
		inCase = append(inCase, s2...)
		_ = inCase
	}

	done := make(chan int)
	select {
	case <-done:
		inComm := append([]int{}, s1...)
		inComm = append(inComm, s2...)
		_ = inComm
	}
}

func useSlicesConcatOK(s1, s2 []int, b []byte) {
	// mutate can modify s1 after the first append has copied it,
	// while slices.Concat copies the slices only once all of them are evaluated
	_ = append(append([]int{}, s1...), mutate(s1)...)

	mutated := append([]int{}, s1...)
	mutated = append(mutated, mutate(s1)...)
	_ = mutated

	_ = append([]int{}, s1...)            // a single append is a clone, not a concatenation
	_ = append(s1, s2...)                 // the appends are not based on an empty slice
	_ = append([]int{0}, s1...)           // the appends are not based on an empty slice
	_ = append(append(s1, s2...), s1...)  // the appends are not based on an empty slice
	_ = append(append([]int{}, s1...), 1) // the last append does not spread a slice

	_ = append(append([]byte{}, b...), "a string"...) // slices.Concat does not accept a string

	one := append([]int{}, s1...)
	one = append(one, 1) // the last append does not spread a slice
	_ = one

	other := append([]int{}, s1...)
	_ = other
	other = append(other, s2...) // the appends are not consecutive
	_ = other

	self := append([]int{}, s1...)
	self = append(self, self...) // the appended slice is the target itself
	_ = self
}

func mutate(s []int) []int {
	s[0] = 0

	return s
}
