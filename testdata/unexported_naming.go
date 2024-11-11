package fixtures

var unexported string
var Exported string

func unexportednaming(
	S int, // MATCH /the symbol S is local, its name should start with a lowercase letter/
	s int,
) (
	Result bool, // MATCH /the symbol Result is local, its name should start with a lowercase letter/
	result bool,
) {
	var NotExportable int // MATCH /the symbol NotExportable is local, its name should start with a lowercase letter/
	var local float32
	{
		OtherNotExportable := 0 // MATCH /the symbol OtherNotExportable is local, its name should start with a lowercase letter/
	}
	const NotExportableConstant = "local" // MATCH /the symbol NotExportableConstant is local, its name should start with a lowercase letter/

	return
}
