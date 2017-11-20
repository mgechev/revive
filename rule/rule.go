package rule

import (
	"go/ast"
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
}

// AbstractRule defines an abstract rule.
type AbstractRule struct {
	failures []Failure
	File     *file.File
}

// AddFailures adds rule failures.
func (r *AbstractRule) AddFailures(failures ...Failure) {
	r.failures = append(r.failures, failures...)
}

// AddFailureAtNode adds rule failure at specific node.
func (r *AbstractRule) AddFailureAtNode(failure Failure, t ast.Node, file *file.File) {
	failure.Position = toFailurePosition(t.Pos(), t.End(), file)
	r.AddFailures(failure)
}

// Failures returns the rule failures.
func (r *AbstractRule) Failures() []Failure {
	return r.failures
}

func toFailurePosition(start token.Pos, end token.Pos, file *file.File) FailurePosition {
	return FailurePosition{
		Start: file.ToPosition(start),
		End:   file.ToPosition(end),
	}
}
