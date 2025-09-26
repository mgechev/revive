package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestRedundantAssignment(t *testing.T) {
	testRule(t, "redundant_assignment", &rule.RedundantAssignmentRule{}, &lint.RuleConfig{})
	testRule(t, "go1.22/redundant_assignment", &rule.RedundantAssignmentRule{}, &lint.RuleConfig{})
}
