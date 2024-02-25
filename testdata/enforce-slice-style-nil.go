package fixtures

func somefn() {
	m0 := make([]string, 10)
	m1 := make([]string, 0, 10)
	m2 := make([]string, 0)    // MATCH /use nil slice declaration (e.g. var args []type) instead of make([]type, 0)/
	m3 := make([]string, 0, 0) // MATCH /use nil slice declaration (e.g. var args []type) instead of make([]type, 0)/
	m4 := []string{}           // MATCH /use nil slice declaration (e.g. var args []type) instead of []type{}/
	m5 := []string{"v1", "v2"}
	m6 := [8]string{}
	m7 := [...]string{}
	var m8 []string

	_ = m0
	_ = m1
	_ = m2
	_ = m3
	_ = m4
	_ = m5
	_ = m6
	_ = m7
	_ = m8
}

type Slice []string

func somefn2() {
	m0 := make(Slice, 10)
	m1 := make(Slice, 0, 10)
	m2 := make(Slice, 0)    // MATCH /use nil slice declaration (e.g. var args []type) instead of make([]type, 0)/
	m3 := make(Slice, 0, 0) // MATCH /use nil slice declaration (e.g. var args []type) instead of make([]type, 0)/
	m4 := Slice{}           // MATCH /use nil slice declaration (e.g. var args []type) instead of []type{}/
	m5 := Slice{"v1", "v2"}
	var m6 Slice

	_ = m0
	_ = m1
	_ = m2
	_ = m3
	_ = m4
	_ = m5
	_ = m6
}

type SliceSlice Slice

func somefn3() {
	m0 := make(SliceSlice, 10)
	m1 := make(SliceSlice, 0, 10)
	m2 := make(SliceSlice, 0)    // MATCH /use nil slice declaration (e.g. var args []type) instead of make([]type, 0)/
	m3 := make(SliceSlice, 0, 0) // MATCH /use nil slice declaration (e.g. var args []type) instead of make([]type, 0)/
	m4 := SliceSlice{}           // MATCH /use nil slice declaration (e.g. var args []type) instead of []type{}/
	m5 := SliceSlice{"v1", "v2"}
	var m6 SliceSlice

	_ = m0
	_ = m1
	_ = m2
	_ = m3
	_ = m4
	_ = m5
	_ = m6
}

func somefn4() {
	m1 := [][]string{} // MATCH /use nil slice declaration (e.g. var args []type) instead of []type{}/
	m1["el0"] = make([]string, 10)
	m1["el1"] = make([]string, 0, 10)
	m1["el2"] = make([]string, 0)    // MATCH /use nil slice declaration (e.g. var args []type) instead of make([]type, 0)/
	m1["el3"] = make([]string, 0, 0) // MATCH /use nil slice declaration (e.g. var args []type) instead of make([]type, 0)/
	m1["el4"] = []string{}           // MATCH /use nil slice declaration (e.g. var args []type) instead of []type{}/
	m1["el5"] = []string{"v1", "v2"}
	m1["el6"] = nil

	_ = m1
}

func somefn5() {
	afunc([]string{})        // MATCH /use nil slice declaration (e.g. var args []type) instead of []type{}/
	afunc(make([]string, 0)) // MATCH /use nil slice declaration (e.g. var args []type) instead of make([]type, 0)/
	afunc([]string(nil))
}

func somefn6() {
	type s struct {
		a []string
	}
	s1 = s{a: []string{}}        // MATCH /use nil slice declaration (e.g. var args []type) instead of []type{}/
	s2 = s{a: make([]string, 0)} // MATCH /use nil slice declaration (e.g. var args []type) instead of make([]type, 0)/
	s3 = s{}
	s4 = s{a: nil}
}
