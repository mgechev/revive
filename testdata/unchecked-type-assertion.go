package fixtures

import "fmt"

var (
	foo  any = "foo"
	bars     = []any{1, 42, "some", "thing"}
)

func handleIgnored() {
	r, _ := foo.(int) // MATCH /type cast result is unchecked in foo.(int) - type assertion result ignored/
}

func handleSkipped() {
	r := foo.(int) // MATCH /type cast result is unchecked in foo.(int) - type assertion will panic if not matched/
}

func handleReturn() int {
	return foo.(int) // MATCH /type cast result is unchecked in foo.(int) - type assertion will panic if not matched/
}

func handleSwitch() {
	switch foo.(int) { // MATCH /type cast result is unchecked in foo.(int) - type assertion will panic if not matched/
	case 0:
	case 1:
		//
	}
}

func handleRange() {
	var some any = bars
	for _, x := range some.([]string) { // MATCH /type cast result is unchecked in some.([]string) - type assertion will panic if not matched/
		fmt.Println(x)
	}
}

func handleTypeSwitch() {
	// Should not be a lint
	switch foo.(type) {
	case string:
	case int:
		//
	}
}

func handleTypeSwitchWithAssignment() {
	// Should not be a lint
	switch n := foo.(type) {
	case string:
	case int:
		//
	}
}

func handleTypeSwitchReturn() {
	// Should not be a lint
	return foo.(type)
}
