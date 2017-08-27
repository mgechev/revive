package rules

import (
	"go/token"

	"github.com/mgechev/golinter/file"
)

const (
	// FailureTypeWarning declares failures of type warning
	FailureTypeWarning = "warning"
	// FailureTypeError declares failures of type error.
	FailureTypeError = "error"
)

// FailureType is the type for the failure types.
type FailureType string

// Failure defines a struct for a linting failure.
type Failure struct {
	Failure  string
	Type     FailureType
	Position FailurePosition
	file     *file.File
}

// GetFilename returns the filename.
func (f *Failure) GetFilename() string {
	return f.Position.Start.Filename
}

// FailurePosition returns the failure position
type FailurePosition struct {
	Start token.Position
	End   token.Position
}

// RuleArguments is type used for the arguments of a rule.
type RuleArguments []string

type RuleConfig struct {
	Name      string
	Arguments RuleArguments
}

// Rule defines an abstract rule.
type Rule interface {
	Apply(file *file.File, args RuleArguments) []Failure
}
