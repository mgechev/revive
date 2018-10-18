package fixtures

import "fmt"

func foo(a, b int) {
	fmt.Printf("single line characters out of limit") // MATCH /line is 105 characters, out of limit 100/
}
