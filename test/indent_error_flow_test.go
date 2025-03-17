package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestIndentErrorFlow(t *testing.T) {
	testRule(t, "indent_error_flow", &rule.IndentErrorFlowRule{})
	testRule(t, "indent_error_flow_scope", &rule.IndentErrorFlowRule{}, &lint.RuleConfig{Arguments: []any{"preserveScope"}})
	testRule(t, "indent_error_flow_scope", &rule.IndentErrorFlowRule{}, &lint.RuleConfig{Arguments: []any{"preserve-scope"}})
}
