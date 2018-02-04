// Test that exported names have correct comments.

// Package pkg does something.
package pkg

import "time"

type T int // MATCH /exported type T should have comment or be unexported/

func (T) F() {} // MATCH /exported method T.F should have comment or be unexported/

// this is a nice type.
// MATCH /comment on exported type U should be of the form "U ..." (with optional leading article)/
type U string

// this is a neat function.
// MATCH /comment on exported method U.G should be of the form "G ..."/
func (U) G() {}

// A V is a string.
type V string

// V.H has a pointer receiver

func (*V) H() {} // MATCH /exported method V.H should have comment or be unexported/

var W = "foo" // MATCH /exported var W should have comment or be unexported/

const X = "bar" // MATCH /exported const X should have comment or be unexported/

var Y, Z int // MATCH /exported var Z should have its own declaration/

// Location should be okay, since the other var name is an underscore.
var Location, _ = time.LoadLocation("Europe/Istanbul") // not Constantinople

// this is improperly documented
// MATCH /comment on exported const Thing should be of the form "Thing ..."/
const Thing = "wonderful"
