// Package golint comment
package golint

// Test cases for enabling checks of exported methods of private types in exported rule
type private struct {
}

// MATCH /comment on exported method private.Method should be of the form "Method ..."/
func (p *private) Method() {
}
