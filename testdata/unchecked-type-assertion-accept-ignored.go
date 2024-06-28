package fixtures

var foo any = "foo"

func handleIgnoredIsOKByConfig() {
	// No lint here bacuse `acceptIgnoredAssertionResult` is set to `true`
	r, _ := foo.(int)
}

func handleSkippedStillFails() {
	r := foo.(int) // MATCH /type cast result is unchecked in foo.(int) - type assertion will panic if not matched/
}
