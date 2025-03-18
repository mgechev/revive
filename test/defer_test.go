package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestDefer(t *testing.T) {
	testRule(t, "defer", &rule.DeferRule{})
}

func TestDeferLoopDisabled(t *testing.T) {
	testRule(t, "defer_loop_disabled", &rule.DeferRule{}, &lint.RuleConfig{
		Arguments: []any{[]any{"return", "recover", "callChain", "methodCall"}},
	})
	testRule(t, "defer_loop_disabled", &rule.DeferRule{}, &lint.RuleConfig{
		Arguments: []any{[]any{"return", "recover", "call-chain", "method-call"}},
	})
}

func TestDeferOthersDisabled(t *testing.T) {
	testRule(t, "defer_only_loop_enabled", &rule.DeferRule{}, &lint.RuleConfig{
		Arguments: []any{[]any{"loop"}},
	})
}
