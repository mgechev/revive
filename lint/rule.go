package lint

import (
	"go/ast"
	"go/token"
)

const (
	// SeverityWarning declares failures of type warning
	SeverityWarning = "warning"
	// SeverityError declares failures of type error.
	SeverityError = "error"
)

// Severity is the type for the failure types.
type Severity string

// FailurePosition returns the failure position
type FailurePosition struct {
	Start token.Position
	End   token.Position
}

// Failure defines a struct for a linting failure.
type Failure struct {
	Failure    string
	RuleName   string
	Category   string
	Position   FailurePosition
	Node       ast.Node `json:"-"`
	Confidence float64
	URL        string
	// For future use
	ReplacementLine string
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
type Arguments = []interface{}

// RuleConfig is type used for the rule configuration.
type RuleConfig struct {
	Arguments Arguments
	Severity  Severity
}

// RulesConfig defines the config for all rules.
type RulesConfig = map[string]RuleConfig

// Config defines the config of the linter.
type Config struct {
	IgnoreGeneratedHeader bool `toml:"ignoreGeneratedHeader"`
	Confidence            float64
	Severity              Severity
	Rules                 RulesConfig `toml:"rule"`
	ErrorCode             int         `toml:"errorCode"`
	WarningCode           int         `toml:"warningCode"`
}

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
