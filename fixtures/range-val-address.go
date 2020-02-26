package fixtures

func rangeValAddress() {
	m := map[string]*string{}

	mySlice := []string{"A", "B", "C"}
	for _, value := range mySlice {
		m["address"] = &value // MATCH /suspicious assignment of 'value'. range-loop variables always have the same address/
	}
}

func rangeValAddress2() {
	m := map[string]*string{}

	mySlice := []string{"A", "B", "C"}
	for i := range mySlice {
		m["address"] = &mySlice[i]
	}
}

func rangeValAddress3() {
	m := map[string]*string{}

	mySlice := []string{"A", "B", "C"}
	for _, value := range mySlice {
		a := &value // MATCH /suspicious assignment of 'value'. range-loop variables always have the same address/
		m["address"] = a
	}
}

func rangeValAddress4() {
	m := []*string{}

	mySlice := []string{"A", "B", "C"}
	for _, value := range mySlice {
		m = append(m, &value) // MATCH /suspicious assignment of 'value'. range-loop variables always have the same address/
	}
}

func rangeValAddress5() {
	m := map[*string]string{}

	mySlice := []string{"A", "B", "C"}
	for _, value := range mySlice {
		m[&value] = value // MATCH /suspicious assignment of 'value'. range-loop variables always have the same address/
	}
}
