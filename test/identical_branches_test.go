package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

// IdenticalBranches rule.
func TestIdenticalBranches(t *testing.T) {
	testRule(t, "identical_branches", &rule.IdenticalBranchesRule{})
}
