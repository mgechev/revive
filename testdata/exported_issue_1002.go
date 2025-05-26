// Package golint comment
package golint

// by default code below is valid,
// but if checkPublicInterface is switched on - it should check documentation in interfaces

// Some - some interface
type Some interface {
	Other // should not fail
	// Correct - should do all correct
	Correct()
	// MATCH /comment on exported interface method Some.SemiCorrect should be of the form "SemiCorrect ..."/
	SemiCorrect()
	NonCorrect() // MATCH /public interface method Some.NonCorrect should be commented/
}

// Other - just to check names compatibility
type Other interface{}

// for private interfaces it doesn't check docs anyway

type somePrivate interface {
	AllGood()
}
