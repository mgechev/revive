package fixtures

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
