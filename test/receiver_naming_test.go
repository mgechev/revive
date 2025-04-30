package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestReceiverNamingTypeParams(t *testing.T) {
	testRule(t, "receiver_naming_issue_669", &rule.ReceiverNamingRule{})
}

func TestReceiverNamingMaxLength(t *testing.T) {
	testRule(t, "receiver_naming_issue_1040", &rule.ReceiverNamingRule{},
		&lint.RuleConfig{
			Arguments: []any{
				map[string]any{"maxLength": int64(2)},
			},
		},
	)
	testRule(t, "receiver_naming_issue_1040", &rule.ReceiverNamingRule{},
		&lint.RuleConfig{
			Arguments: []any{
				map[string]any{"max-length": int64(2)},
			},
		},
	)
}
