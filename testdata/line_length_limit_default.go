package fixtures

import "fmt"

func foo(a, b int) {
	fmt.Printf("loooooooooooooooooooooooooooooooooooooooooong line out of limit")
}

// MATCH:6 /line is 81 characters, out of limit 80/
