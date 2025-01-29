package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestIdenticalBranches(t *testing.T) {
	testRule(t, "identical_branches", &rule.IdenticalBranchesRule{})
}
