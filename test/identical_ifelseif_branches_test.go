package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestIdenticalIfElseIfBranches(t *testing.T) {
	testRule(t, "identical_ifelseif_branches", &rule.IdenticalIfElseIfBranchesRule{})
}
