package fixtures

func typeInferenceExample() {
	var x int = 42 // MATCH /should omit type int from declaration of var x; it will be inferred from the right-hand side/
	var y = 42     // No warning, type is inferred
}
