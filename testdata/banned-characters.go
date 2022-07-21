package fixtures

const Ω = "Omega" // MATCH /banned character found: Ω/

// func contains banned characters Ω // authorized banned chars in comment
func funcΣ() error { // MATCH /banned character found: Σ/
	var charσhid string // MATCH /banned character found: σ/
	_ = charσhid        // MATCH /banned character found: σ/
	return nil
}
