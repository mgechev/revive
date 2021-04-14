// Test string literal regex checks

package pkg

func stringRegexMethod1(a int, b string) {

}

func stringRegex() {
	var (
		s1 = "This string is fine"
		s2 = "this string is not capitalized"                      // MATCH /string literal doesn't match user defined regex (must start with a capital letter)/
		s3 = "This string is capitalized, but ends with a period." // MATCH /string literal doesn't match user defined regex /[^\.]$//
	)

	stringRegexMethod1(0, s2) // MATCH /string literal doesn't match user defined regex (must start with a capital letter)/
}
