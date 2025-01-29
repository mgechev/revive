package fixtures

func foo() string {
	customID := "result"    // ignore
	customVm := "result"    // json:{"MATCH": "var customVm should be customVM"}
	CUSTOMER_UP := "result" // json:{"MATCH": "don't use ALL_CAPS in Go names; use CamelCase", "Confidence": 0.8}
	return customId
}
