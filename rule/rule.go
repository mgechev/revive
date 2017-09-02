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

type Config struct {
	Name      string
	Arguments Arguments
}

// Rule defines an abstract rule.
type Rule interface {
	GetName() string
	Apply(file *file.File, args Arguments) []Failure
}
