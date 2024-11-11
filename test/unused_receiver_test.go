package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestUnusedReceiver(t *testing.T) {
	testRule(t, "unused_receiver", &rule.UnusedReceiverRule{})
	testRule(t, "unused_receiver_custom_regex", &rule.UnusedReceiverRule{}, &lint.RuleConfig{Arguments: []any{
		map[string]any{"allowRegex": "^xxx"},
	}})
}
