package visitors

import (
	"go/token"
)

type RuleArguments []string

const DefaultLength = 1

type Failure struct {
	Failure  string
	Position token.Pos
}

type RuleVisitor struct {
	SyntaxVisitor
	ruleName      string
	ruleArguments RuleArguments
	failures      []Failure
}

func New(ruleName string, ruleArguments RuleArguments) *RuleVisitor {
	result := RuleVisitor{ruleName: ruleName, ruleArguments: ruleArguments}
	result.failures = make([]Failure, DefaultLength)
	return &result
}

func (w *RuleVisitor) AddFailure(failure Failure) {
	w.failures = append(w.failures, failure)
}

func (w *RuleVisitor) GetFailures() []Failure {
	return w.failures
}
