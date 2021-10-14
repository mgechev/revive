package fixtures

const Ω = "Omega" // MATCH:3 /banned character found: Ω/

// func contains banned characters Ω // authorized banned chars in comment
func funcΣ() error { // MATCH:6 /banned character found: Σ/
	var charσhid string // MATCH:7 /banned character found: σ/
	return nil
}
