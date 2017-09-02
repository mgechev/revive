package visitor

import (
	"go/token"

	"github.com/mgechev/revive/file"
	"github.com/mgechev/revive/rule"
)

// RuleVisitor defines a struct for a visitor.
type RuleVisitor struct {
	SyntaxVisitor
	RuleName      string
	RuleArguments rule.Arguments
	failures      []rule.Failure
	File          *file.File
}

// AddFailure adds a failure to the ist of failures.
func (w *RuleVisitor) AddFailure(failure rule.Failure) {
	w.failures = append(w.failures, failure)
}

// GetFailures returns the list of failures.
func (w *RuleVisitor) GetFailures() []rule.Failure {
	return w.failures
}

// GetPosition returns position by given start and end token.Pos.
func (w *RuleVisitor) GetPosition(start token.Pos, end token.Pos) rule.FailurePosition {
	s := w.File.ToPosition(start)
	e := w.File.ToPosition(end)
	return rule.FailurePosition{
		Start: s,
		End:   e,
	}
}
