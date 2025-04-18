package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestUnusedReceiver(t *testing.T) {
	testRule(t, "unused_receiver", &rule.UnusedReceiverRule{})
	testRule(t, "unused_receiver", &rule.UnusedReceiverRule{}, &lint.RuleConfig{Arguments: []any{}})
	testRule(t, "unused_receiver", &rule.UnusedReceiverRule{}, &lint.RuleConfig{Arguments: []any{
		map[string]any{"a": "^xxx"},
	}})
	testRule(t, "unused_receiver_custom_regex", &rule.UnusedReceiverRule{}, &lint.RuleConfig{Arguments: []any{
		map[string]any{"allowRegex": "^xxx"},
	}})
	testRule(t, "unused_receiver_custom_regex", &rule.UnusedReceiverRule{}, &lint.RuleConfig{Arguments: []any{
		map[string]any{"allow-regex": "^xxx"},
	}})
}

func BenchmarkUnusedReceiver(b *testing.B) {
	var t *testing.T
	for i := 0; i <= b.N; i++ {
		testRule(t, "unused_receiver", &rule.UnusedReceiverRule{})
	}
}
