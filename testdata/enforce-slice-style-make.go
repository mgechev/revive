package fixtures

import "fmt"

func somefn() {
	m0 := make([]string, 10)
	m1 := make([]string, 0, 10)
	m2 := make([]string, 0)
	m3 := make([]string, 0, 0)
	m4 := []string{} // MATCH /use make([]type) instead of []type{} (or declare nil slice)/
	m5 := []string{"v1", "v2"}
	m6 := [8]string{}
	m7 := [...]string{}
	var m8 []string
	var m9 []string = make([]string, 0)
	var m10 = make([]string, 0)
	m11 := []string(nil)
	var m12 []string = []string{} // MATCH /use make([]type) instead of []type{} (or declare nil slice)/
	var m13 = []string{}          // MATCH /use make([]type) instead of []type{} (or declare nil slice)/
	var m14 = []string(nil)
	var m15 []string = nil
	var m16 []string
	m16 = make([]string, 0)
	m16 = []string{} // MATCH /use make([]type) instead of []type{} (or declare nil slice)/

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
	m2 := make(Slice, 0)
	m3 := make(Slice, 0, 0)
	m4 := Slice{} // MATCH /use make([]type) instead of []type{} (or declare nil slice)/
	m5 := Slice{"v1", "v2"}
	var m6 Slice
	var m7 Slice = make(Slice, 0)
	var m8 = make(Slice, 0)
	m9 := Slice(nil)
	var m10 Slice = Slice{} // MATCH /use make([]type) instead of []type{} (or declare nil slice)/
	var m11 = Slice{}       // MATCH /use make([]type) instead of []type{} (or declare nil slice)/
	var m12 = Slice(nil)
	var m13 Slice = nil
	var m14 Slice
	m14 = make(Slice, 0)
	m14 = Slice{} // MATCH /use make([]type) instead of []type{} (or declare nil slice)/

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
	m2 := make(SliceSlice, 0)
	m3 := make(SliceSlice, 0, 0)
	m4 := SliceSlice{} // MATCH /use make([]type) instead of []type{} (or declare nil slice)/
	m5 := SliceSlice{"v1", "v2"}
	var m6 SliceSlice
	var m7 SliceSlice = make(SliceSlice, 0)
	var m8 = make(SliceSlice, 0)
	m9 := SliceSlice(nil)
	var m10 SliceSlice = SliceSlice{} // MATCH /use make([]type) instead of []type{} (or declare nil slice)/
	var m11 = SliceSlice{}            // MATCH /use make([]type) instead of []type{} (or declare nil slice)/
	var m12 = SliceSlice(nil)
	var m13 SliceSlice = nil
	var m14 SliceSlice
	m14 = make(SliceSlice, 0)
	m14 = SliceSlice{} // MATCH /use make([]type) instead of []type{} (or declare nil slice)/

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
	m1 := [][]string{} // MATCH /use make([]type) instead of []type{} (or declare nil slice)/
	m1["el0"] = make([]string, 10)
	m1["el1"] = make([]string, 0, 10)
	m1["el2"] = make([]string, 0)
	m1["el3"] = make([]string, 0, 0)
	m1["el4"] = []string{} // MATCH /use make([]type) instead of []type{} (or declare nil slice)/
	m1["el5"] = []string{"v1", "v2"}
	m1["el6"] = nil

	_ = m1
}

func somefn5() {
	afunc([]string{}) // MATCH /use make([]type) instead of []type{} (or declare nil slice)/
	afunc(make([]string, 0))
	afunc([]string(nil))
}

func somefn6() {
	type s struct {
		a []string
	}
	s1 = s{a: []string{}} // MATCH /use make([]type) instead of []type{} (or declare nil slice)/
	s2 = s{a: make([]string, 0)}
	s3 = s{}
	s4 = s{a: nil}
}

func somefn7() {
	if m := []string(nil); len(m) == 0 {
		fmt.Println("Hello, 世界")
	}
	if m := make([]string, 0); len(m) == 0 {
		fmt.Println("Hello, 世界")
	}
	if m := []string{}; len(m) == 0 { // MATCH /use make([]type) instead of []type{} (or declare nil slice)/
		fmt.Println("Hello, 世界")
	}
}
