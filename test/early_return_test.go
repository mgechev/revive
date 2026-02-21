package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestEarlyReturn(t *testing.T) {
	testRule(t, "early_return", &rule.EarlyReturnRule{})
	testRule(t, "early_return_scope", &rule.EarlyReturnRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"preserveScope"}})
	testRule(t, "early_return_scope", &rule.EarlyReturnRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"preserve-scope"}})
	testRule(t, "early_return_jump", &rule.EarlyReturnRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"allowJump"}})
	testRule(t, "early_return_jump", &rule.EarlyReturnRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"allow-jump"}})
	testRule(t, "early_return_jump_scope", &rule.EarlyReturnRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"allow-jump", "preserve-scope"}})
}
