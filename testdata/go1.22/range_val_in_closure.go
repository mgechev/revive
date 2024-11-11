package fixtures

import "fmt"

func foo() {
	mySlice := []string{"A", "B", "C"}
	for index, value := range mySlice {
		go func() {
			fmt.Printf("Index: %d\n", index) // Shall not match
			fmt.Printf("Value: %s\n", value) // Shall not match
		}()
	}

	myDict := make(map[string]int)
	myDict["A"] = 1
	myDict["B"] = 2
	myDict["C"] = 3
	for key, value := range myDict {
		defer func() {
			fmt.Printf("Index: %d\n", key)   // Shall not match
			fmt.Printf("Value: %s\n", value) // Shall not match
		}()
	}

	for i, newg := range groups {
		go func(newg int) {
			newg.run(m.opts.Context, i) // Shall not match
		}(newg)
	}

	for i, newg := range groups {
		newg := newg
		go func() {
			newg.run(m.opts.Context, i) // Shall not match
		}()
	}
}

func issue637() {
	for key := range m {
		myKey := key
		go func() {
			println(t{
				key:        myKey,
				otherField: (10 + key), // Shall not match
			})
		}()
	}
}
