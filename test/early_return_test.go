package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestEarlyReturn(t *testing.T) {
	testRule(t, "early_return", &rule.EarlyReturnRule{})
	testRule(t, "early_return_scope", &rule.EarlyReturnRule{}, &lint.RuleConfig{Arguments: []any{"preserveScope"}})
	testRule(t, "early_return_scope", &rule.EarlyReturnRule{}, &lint.RuleConfig{Arguments: []any{"preserve-scope"}})
	testRule(t, "early_return_jump", &rule.EarlyReturnRule{}, &lint.RuleConfig{Arguments: []any{"allowJump"}})
	testRule(t, "early_return_jump", &rule.EarlyReturnRule{}, &lint.RuleConfig{Arguments: []any{"allow-jump"}})
	testRule(t, "early_return_jump_scope", &rule.EarlyReturnRule{}, &lint.RuleConfig{Arguments: []any{"allow-jump", "preserve-scope"}})
}
