package fixtures

func useNewOne(a int) *int { // MATCH /calls to "useNewOne(value)" can be replaced by a call to "new(value)"/
	return &a
}
