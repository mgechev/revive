package linter

import (
	"go/ast"
	"go/token"
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
	Failure    string
	RuleName   string
	Type       FailureType
	Position   FailurePosition
	Node       ast.Node
	Confidence float64
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

// RulesConfig defiles the config for all rules.
type RulesConfig = map[string]Arguments

// Rule defines an abstract rule interaface
type Rule interface {
	Name() string
	Apply(*File, Arguments) []Failure
}

// AbstractRule defines an abstract rule.
type AbstractRule struct {
	Failures []Failure
}

// ToFailurePosition returns the failure position.
func ToFailurePosition(start token.Pos, end token.Pos, file *File) FailurePosition {
	return FailurePosition{
		Start: file.ToPosition(start),
		End:   file.ToPosition(end),
	}
}
