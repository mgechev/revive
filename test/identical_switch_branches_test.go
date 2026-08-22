package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestIdenticalSwitchBranches(t *testing.T) {
	testRule(t, "identical_switch_branches", &rule.IdenticalSwitchBranchesRule{})
	testRule(t, "identical_switch_branches_allow_identical_default", &rule.IdenticalSwitchBranchesRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"allow-identical-default"},
	})
}
