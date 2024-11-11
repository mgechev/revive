package test

import (
	"testing"

	"github.com/mgechev/revive/internal/typeparams"
	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestReceiverNamingTypeParams(t *testing.T) {
	if !typeparams.Enabled() {
		t.Skip("type parameters are not enabled in the current build environment")
	}
	testRule(t, "receiver_naming_issue_669", &rule.ReceiverNamingRule{})
}

func TestReceiverNamingMaxLength(t *testing.T) {
	args := []any{map[string]any{
		"maxLength": int64(2),
	}}
	testRule(t, "receiver_naming_issue_1040", &rule.ReceiverNamingRule{},
		&lint.RuleConfig{Arguments: args})
}
