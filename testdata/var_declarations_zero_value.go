package fixtures

func zeroValueExample() {
	var y int                // No warning, zero value is implicit
	var z1 any = 0           // No warning, zero value for any is nil
	var z2 any = nil         // MATCH /should drop = nil from declaration of var z2; it is the zero value/
	var z3 interface{} = 0   // No warning, zero value for any is nil
	var z4 interface{} = nil // MATCH /should drop = nil from declaration of var z4; it is the zero value/
}
