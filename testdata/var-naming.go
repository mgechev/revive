package fixtures

func foo() string {
	customId := "result"
	customVm := "result"  // MATCH /var customVm should be customVM/
	customIds := "result" // MATCH /var customIds should be customIDs/
	return customId
}
