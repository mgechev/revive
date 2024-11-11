package fixtures

import "fmt"

func foo() {
	mySlice := []string{"A", "B", "C"}
	for index, value := range mySlice {
		go func() {
			fmt.Printf("Index: %d\n", index) // MATCH /loop variable index captured by func literal/
			fmt.Printf("Value: %s\n", value) // MATCH /loop variable value captured by func literal/
		}()
	}

	myDict := make(map[string]int)
	myDict["A"] = 1
	myDict["B"] = 2
	myDict["C"] = 3
	for key, value := range myDict {
		defer func() {
			fmt.Printf("Index: %d\n", key)   // MATCH /loop variable key captured by func literal/
			fmt.Printf("Value: %s\n", value) // MATCH /loop variable value captured by func literal/
		}()
	}

	for i, newg := range groups {
		go func(newg int) {
			newg.run(m.opts.Context, i) // MATCH /loop variable i captured by func literal/
		}(newg)
	}

	for i, newg := range groups {
		newg := newg
		go func() {
			newg.run(m.opts.Context, i) // MATCH /loop variable i captured by func literal/
		}()
	}
}

func issue637() {
	for key := range m {
		myKey := key
		go func() {
			println(t{
				key:        myKey,
				otherField: (10 + key), // MATCH /loop variable key captured by func literal/
			})
		}()
	}
}
