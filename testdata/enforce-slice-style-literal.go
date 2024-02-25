package fixtures

func somefn() {
	m0 := make([]string, 10)
	m1 := make([]string, 0, 10)
	m2 := make([]string, 0)    // MATCH /use []type{} instead of make([]type, 0) (or declare nil slice)/
	m3 := make([]string, 0, 0) // MATCH /use []type{} instead of make([]type, 0) (or declare nil slice)/
	m4 := []string{}
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
	m2 := make(Slice, 0)    // MATCH /use []type{} instead of make([]type, 0) (or declare nil slice)/
	m3 := make(Slice, 0, 0) // MATCH /use []type{} instead of make([]type, 0) (or declare nil slice)/
	m4 := Slice{}
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
	m2 := make(SliceSlice, 0)    // MATCH /use []type{} instead of make([]type, 0) (or declare nil slice)/
	m3 := make(SliceSlice, 0, 0) // MATCH /use []type{} instead of make([]type, 0) (or declare nil slice)/
	m4 := SliceSlice{}
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
	m1 := [][]string{}
	m1["el0"] = make([]string, 10)
	m1["el1"] = make([]string, 0, 10)
	m1["el2"] = make([]string, 0)    // MATCH /use []type{} instead of make([]type, 0) (or declare nil slice)/
	m1["el3"] = make([]string, 0, 0) // MATCH /use []type{} instead of make([]type, 0) (or declare nil slice)/
	m1["el4"] = []string{}
	m1["el5"] = []string{"v1", "v2"}
	m1["el6"] = nil

	_ = m1
}
