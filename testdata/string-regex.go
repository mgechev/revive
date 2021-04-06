// Test forced string capitalization

package pkg

func stringFormatCapitalization() {
	var (
		s1 = "This string is fine"
		s2 = "this string is not capitalized"                      // MATCH /message doesn't match user defined regex (must be capitalized)/
		s3 = "This string is capitalized, but ends with a period." // MATCH /message doesn't match user defined regex /[^\.]$//
	)
}
