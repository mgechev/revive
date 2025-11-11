package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestIdenticalSwitchBranches(t *testing.T) {
	testRule(t, "identical_switch_branches", &rule.IdenticalSwitchBranchesRule{})
}
