package fixtures

func redundantAssignment() {
	for a, b := range collection {
		a := a // MATCH /redundant assignment of range variable "a", use it directly/
		something(a, b)
	}

	for a, b := range collection {
		for a, b := range collection {
			a := a // MATCH /redundant assignment of range variable "a", use it directly/
			something(a, b)
		}
	}

	for a, b := range collection {
		b := b // MATCH /redundant assignment of range variable "b", use it directly/
		something(a, b)
	}

	for a := range collection {
		a := a // MATCH /redundant assignment of range variable "a", use it directly/
		something(a, b)
	}

	for _, a := range collection {
		a := a // MATCH /redundant assignment of range variable "a", use it directly/
		something(a, b)
	}

	// should not report
	for range collection {
		a := a + b
		something(a, b)
	}

	for a, b := range collection {
		a := a + b
		something(a, b)
	}
}
