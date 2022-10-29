package commentspacings

//revive:disable:cyclomatic
type some struct{}

//revive:disable:cyclomatic High complexity score but easy to understand
type some1 struct{}

// Some nice comment
//
// With an empty line in the middle will make the rule panic!
func itsATrap() {}

// This is a well formed comment
type hello struct {
	random `json:random`
}

// MATCH:21 /no space between comment delimiter and comment text/

//This comment does not respect the spacing rule!
var a string

//myOwnDirective: do something
