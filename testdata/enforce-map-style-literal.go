package fixtures

func somefn() {
	m0 := make(map[string]string, 10)
	m1 := make(map[string]string) // MATCH /use map[type]type{} instead of make(map[type]type)/
	m2 := map[string]string{}
	m3 := map[string]string{"k1": "v1", "k2": "v2"}

	_ = m0
	_ = m1
	_ = m2
	_ = m3
}

type Map map[string]string

func somefn2() {
	m0 := make(Map, 10)
	m1 := make(Map) // MATCH /use map[type]type{} instead of make(map[type]type)/
	m2 := Map{}
	m3 := Map{"k1": "v1", "k2": "v2"}

	_ = m0
	_ = m1
	_ = m2
	_ = m3
}

type MapMap Map

func somefn3() {
	m0 := make(MapMap, 10)
	m1 := make(MapMap) // MATCH /use map[type]type{} instead of make(map[type]type)/
	m2 := MapMap{}
	m3 := MapMap{"k1": "v1", "k2": "v2"}

	_ = m0
	_ = m1
	_ = m2
	_ = m3
}

func somefn4() {
	m1 := map[string]map[string]string{}
	m1["el0"] = make(map[string]string, 10)
	m1["el1"] = make(map[string]string) // MATCH /use map[type]type{} instead of make(map[type]type)/
	m1["el2"] = map[string]string{}
	m1["el3"] = map[string]string{"k1": "v1", "k2": "v2"}

	_ = m1
}
