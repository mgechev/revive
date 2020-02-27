package fixtures

func foo() string {
	customId := "result"
	customVm := "result" // MATCH /var customVm should be customVM/
	return customId
}
