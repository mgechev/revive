package fixtures // MATCH /the file has a comment density of 57% (4 comment lines for 3 code lines) but expected a minimum of 60%/

// func contains banned characters Ω // authorized banned chars in comment
func cd1() error {
	// the var
	var charσhid string
	/* the return */
	return nil
}
