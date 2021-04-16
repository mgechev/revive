// Test string literal regex checks

package pkg

func stringRegexMethod1(a, b string) {

}

func stringRegexMethod2(a, b string, c struct {
	d string
}) {

}

type stringRegexMethods struct{}

func (s stringRegexMethods) Method3(a, b, c string) {

}

func stringRegex() {
	stringRegexMethod1("This string is fine", "")
	stringRegexMethod1("this string is not capitalized", "") // MATCH /string literal doesn't match user defined regex (must start with a capital letter)/
	stringRegexMethod2(s3, "", struct {
		d string
	}{
		d: "This string is capitalized, but ends with a period."}) // MATCH /string literal doesn't match user defined regex /[^\.]$//
	s := stringRegexMethods{}
	s.Method3("", "", "This string starts with th") // MATCH /string literal doesn't match user defined regex (must not start with 'th')/
}
