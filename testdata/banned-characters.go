package fixtures

// func contains banned characters Ω
// MATCH:0 /banned character found: Ω/
func funcΣ() error { // MATCH:0 /banned character found: Σ/
	var charσhid string // MATCH:0 /banned character found: σ/
	return nil
}
