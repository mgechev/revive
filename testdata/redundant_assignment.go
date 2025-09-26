package fixtures

func redundantAssignment() {
	for a, b := range collection {
		a := a
		something(a, b)
	}

	for a, b := range collection {
		for a, b := range collection {
			a := a
			something(a, b)
		}
	}

	for a, b := range collection {
		b := b
		something(a, b)
	}

	for a := range collection {
		a := a
		something(a, b)
	}

	for _, a := range collection {
		a := a
		something(a, b)
	}

	for range collection {
		a := a + b
		something(a, b)
	}

	for a, b := range collection {
		a := a + b
		something(a, b)
	}
}
