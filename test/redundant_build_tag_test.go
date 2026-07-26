package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestRedundantBuildTagRule(t *testing.T) {
	testRule(t, "redundant_build_tag", &rule.RedundantBuildTagRule{}, &lint.RuleConfig{})
}

func TestRedundantBuildTagRuleNoFailure(t *testing.T) {
	testRule(t, "redundant_build_tag_no_failure", &rule.RedundantBuildTagRule{}, &lint.RuleConfig{})
}

func TestRedundantBuildTagRuleGo116(t *testing.T) {
	testRule(t, "redundant_build_tag_go116", &rule.RedundantBuildTagRule{}, &lint.RuleConfig{})
}

func TestRedundantBuildTagRuleGo121(t *testing.T) {
	testRule(t, "go1.21/redundant_build_tag", &rule.RedundantBuildTagRule{}, &lint.RuleConfig{})
	testRule(t, "go1.21/redundant_build_tag_go120", &rule.RedundantBuildTagRule{}, &lint.RuleConfig{})
	testRule(t, "go1.21/redundant_build_tag_go122", &rule.RedundantBuildTagRule{}, &lint.RuleConfig{})
}
