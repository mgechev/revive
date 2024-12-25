package lint

import (
	"go/ast"
	"go/token"
)

const (
	FailureCategoryArgOrder            = "arg-order"
	FailureCategoryBadPractice         = "bad practice"
	FailureCategoryCodeStyle           = "code-style"
	FailureCategoryComments            = "comments"
	FailureCategoryComplexity          = "complexity"
	FailureCategoryContent             = "content"
	FailureCategoryErrors              = "errors"
	FailureCategoryImports             = "imports"
	FailureCategoryLogic               = "logic"
	FailureCategoryMaintenance         = "maintenance"
	FailureCategoryNaming              = "naming"
	FailureCategoryOptimization        = "optimization"
	FailureCategoryStyle               = "style"
	FailureCategoryTime                = "time"
	FailureCategoryTypeInference       = "type-inference"
	FailureCategoryUnaryOp             = "unary-op"
	FailureCategoryUnexportedTypeInAPI = "unexported-type-in-api"
	FailureCategoryZeroValue           = "zero-value"

	failureCategoryInternal = "REVIVE_INTERNAL"
	failureCategoryValidity = "validity"
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
	// For future use
	ReplacementLine string
}

// GetFilename returns the filename.
func (f *Failure) GetFilename() string {
	return f.Position.Start.Filename
}

// IsInternal returns true if this failure is internal, false otherwise.
func (f *Failure) IsInternal() bool {
	return f.Category == failureCategoryInternal
}

// NewInternalFailure yields an internal failure with the given message as failure message.
func NewInternalFailure(message string) Failure {
	return Failure{
		Category: failureCategoryInternal,
		Failure:  message,
	}
}
