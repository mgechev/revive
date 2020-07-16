package fixtures

var unexported string
var Exported string

func unexportednaming(
	S int, // MATCH /the symbol S can not be exported, its name should start with a lowercase letter/
	s int,
) (
	Result bool, // MATCH /the symbol Result can not be exported, its name should start with a lowercase letter/
	result bool,
) {
	var NotExportable int // MATCH /the symbol NotExportable can not be exported, its name should start with a lowercase letter/
	var local float32
	{
		OtherNotExportable := 0 // MATCH /the symbol OtherNotExportable can not be exported, its name should start with a lowercase letter/
	}
	const NotExportableConstant = "local" // MATCH /the symbol NotExportableConstant can not be exported, its name should start with a lowercase letter/

	return
}
