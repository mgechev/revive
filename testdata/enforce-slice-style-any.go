package fixtures

func somefn() {
	m1 := make([]string)
	m2 := []string{}
	m3 := make([]string, 10)
	m4 := make([]string, 0, 10)
	m5 := []string{"v1"}
	m6 := [...]string{}
	m7 := [...]string{"v1"}
	m8 := [1]string{"v1"}

	_ = m1
	_ = m2
	_ = m3
	_ = m4
	_ = m5
	_ = m6
	_ = m7
	_ = m8
}
