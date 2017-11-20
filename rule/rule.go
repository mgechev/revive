package rule

import (
	"go/token"

	"github.com/mgechev/revive/file"
)

const (
	// FailureTypeWarning declares failures of type warning
	FailureTypeWarning = "warning"
	// FailureTypeError declares failures of type error.
	FailureTypeError = "error"
)

// FailureType is the type for the failure types.
type FailureType string

// FailurePosition returns the failure position
type FailurePosition struct {
	Start token.Position
	End   token.Position
}

// Failure defines a struct for a linting failure.
type Failure struct {
	Failure  string
	RuleName string
	Type     FailureType
	Position FailurePosition
	file     *file.File
}

// GetFilename returns the filename.
func (f *Failure) GetFilename() string {
	return f.Position.Start.Filename
}

// DisabledInterval contains a single disabled interval and the associated rule name.
type DisabledInterval struct {
	From     token.Position
	To       token.Position
	RuleName string
}

// Arguments is type used for the arguments of a rule.
type Arguments []string

// Config contains the rule configuration.
type Config struct {
	Name      string
	Arguments Arguments
}

// Rule defines an abstract rule interaface
type Rule interface {
	Name() string
	Apply(*file.File, Arguments) []Failure
	AddFailures(...Failure)
	Failures() []Failure
	Position(token.Pos, token.Pos) FailurePosition
}

// AbstractRule defines an abstract rule.
type AbstractRule struct {
	failures []Failure
	File     *file.File
}

// Apply must be overridden by the successor struct.
func (r *AbstractRule) Apply(file *file.File, args Arguments) {
	panic("Apply not implemented")
}

// AddFailures adds rule failures.
func (r *AbstractRule) AddFailures(failures ...Failure) {
	r.failures = append(r.failures, failures...)
}

// Failures returns the rule failures.
func (r *AbstractRule) Failures() []Failure {
	return r.failures
}

// Position returns position by given start and end token.Pos.
func (r *AbstractRule) Position(start token.Pos, end token.Pos) FailurePosition {
	s := r.File.ToPosition(start)
	e := r.File.ToPosition(end)
	return FailurePosition{
		Start: s,
		End:   e,
	}
}
