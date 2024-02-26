package fixtures

import "fmt"

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
	var m9 []string = make([]string, 0) // MATCH /use []type{} instead of make([]type, 0) (or declare nil slice)/
	var m10 = make([]string, 0)         // MATCH /use []type{} instead of make([]type, 0) (or declare nil slice)/
	m11 := []string(nil)
	var m12 []string = []string{}
	var m13 = []string{}
	var m14 = []string(nil)
	var m15 []string = nil
	var m16 []string
	m16 = make([]string, 0) // MATCH /use []type{} instead of make([]type, 0) (or declare nil slice)/
	m16 = []string{}

	_ = m0
	_ = m1
	_ = m2
	_ = m3
	_ = m4
	_ = m5
	_ = m6
	_ = m7
	_ = m8
	_ = m9
	_ = m10
	_ = m11
	_ = m12
	_ = m13
	_ = m14
	_ = m15
	_ = m16
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
	var m7 Slice = make(Slice, 0) // MATCH /use []type{} instead of make([]type, 0) (or declare nil slice)/
	var m8 = make(Slice, 0)       // MATCH /use []type{} instead of make([]type, 0) (or declare nil slice)/
	m9 := Slice(nil)
	var m10 Slice = Slice{}
	var m11 = Slice{}
	var m12 = Slice(nil)
	var m13 Slice = nil
	var m14 Slice
	m14 = make(Slice, 0) // MATCH /use []type{} instead of make([]type, 0) (or declare nil slice)/
	m14 = Slice{}

	_ = m0
	_ = m1
	_ = m2
	_ = m3
	_ = m4
	_ = m5
	_ = m6
	_ = m7
	_ = m8
	_ = m9
	_ = m10
	_ = m11
	_ = m12
	_ = m13
	_ = m14
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
	var m7 SliceSlice = make(SliceSlice, 0) // MATCH /use []type{} instead of make([]type, 0) (or declare nil slice)/
	var m8 = make(SliceSlice, 0)            // MATCH /use []type{} instead of make([]type, 0) (or declare nil slice)/
	m9 := SliceSlice(nil)
	var m10 SliceSlice = SliceSlice{}
	var m11 = SliceSlice{}
	var m12 = SliceSlice(nil)
	var m13 SliceSlice = nil
	var m14 SliceSlice
	m14 = make(SliceSlice, 0) // MATCH /use []type{} instead of make([]type, 0) (or declare nil slice)/
	m14 = SliceSlice{}

	_ = m0
	_ = m1
	_ = m2
	_ = m3
	_ = m4
	_ = m5
	_ = m6
	_ = m7
	_ = m8
	_ = m9
	_ = m10
	_ = m11
	_ = m12
	_ = m13
	_ = m14
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

func somefn5() {
	afunc([]string{})
	afunc(make([]string, 0)) // MATCH /use []type{} instead of make([]type, 0) (or declare nil slice)/
	afunc([]string(nil))
}

func somefn6() {
	type s struct {
		a []string
	}
	s1 = s{a: []string{}}
	s2 = s{a: make([]string, 0)} // MATCH /use []type{} instead of make([]type, 0) (or declare nil slice)/
	s3 = s{}
	s4 = s{a: nil}
}

func somefn7() {
	if m := []string(nil); len(m) == 0 {
		fmt.Println("Hello, 世界")
	}
	if m := make([]string, 0); len(m) == 0 { // MATCH /use []type{} instead of make([]type, 0) (or declare nil slice)/
		fmt.Println("Hello, 世界")
	}
	if m := []string{}; len(m) == 0 {
		fmt.Println("Hello, 世界")
	}
}
