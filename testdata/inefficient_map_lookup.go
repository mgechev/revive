package fixtures

import "fmt"

func inefficientMapLookup() {
	type aS struct {
		TagIDs map[int]string
	}
	a := aS{}
	someStaticValue := 1
	// use case from issue #1447
	for id := range a.TagIDs { // MATCH /inefficient lookup of map key/
		if id == someStaticValue {
			return
		}
	}

	// do not match if the loop body is something else than
	// just an if on the map key
	aMap := map[int]int{}
	for k := range aMap {
		fmt.Println(k)
		if k == 1 {
			return
		}
	}

	for k := range aMap {
		if k == 1 {
			return
		}
		fmt.Println(k)
	}
}
