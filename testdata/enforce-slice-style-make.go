package fixtures

func somefn() {
	m0 := make([]string, 10)
	m1 := make([]string, 0, 10)
	m2 := make([]string, 0)
	m3 := []string{} // MATCH /use make([]type) instead of []type{} (or declare nil slice)/
	m4 := []string{"v1", "v2"}
	m5 := [8]string{}
	m6 := [...]string{}

	_ = m0
	_ = m1
	_ = m2
	_ = m3
	_ = m4
	_ = m5
	_ = m6
}

type Slice []string

func somefn2() {
	m0 := make(Slice, 10)
	m1 := make(Slice, 0, 10)
	m2 := make(Slice, 0)
	m3 := Slice{} // MATCH /use make([]type) instead of []type{} (or declare nil slice)/
	m4 := Slice{"v1", "v2"}

	_ = m0
	_ = m1
	_ = m2
	_ = m3
	_ = m4
}

type SliceSlice Slice

func somefn3() {
	m0 := make(SliceSlice, 10)
	m1 := make(SliceSlice, 0, 10)
	m2 := make(SliceSlice, 0)
	m3 := SliceSlice{} // MATCH /use make([]type) instead of []type{} (or declare nil slice)/
	m4 := SliceSlice{"v1", "v2"}

	_ = m0
	_ = m1
	_ = m2
	_ = m3
	_ = m4
}

func somefn4() {
	m1 := [][]string{} // MATCH /use make([]type) instead of []type{} (or declare nil slice)/
	m1["el0"] = make([]string, 10)
	m1["el1"] = make([]string, 0, 10)
	m1["el2"] = make([]string, 0)
	m1["el3"] = []string{} // MATCH /use make([]type) instead of []type{} (or declare nil slice)/
	m1["el4"] = []string{"v1", "v2"}

	_ = m1
}
