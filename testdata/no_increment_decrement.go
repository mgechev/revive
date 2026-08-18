package fixtures

func standaloneIncrementDecrement() {
	x := 0
	x++ // MATCH /should replace x++ with x += 1/
	x-- // MATCH /should replace x-- with x -= 1/
}

func loopCounters() {
	for i := 0; i < 10; i++ {
		_ = i
	}
	for j := 10; j > 0; j-- {
		_ = j
	}
}

func incrementInLoopBody() {
	count := 0
	for range []int{1, 2, 3} {
		count++ // MATCH /should replace count++ with count += 1/
	}
	_ = count
}

func selectorOperand(s *struct{ n int }) {
	s.n++ // MATCH /should replace s.n++ with s.n += 1/
}
