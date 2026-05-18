package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestReturnsInterfaceTypeDefault(t *testing.T) {
	testRule(t, "returns_interface_type_default", &rule.ReturnsInterfaceTypeRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{}},
	})
	testRule(t, "returns_interface_type_report_all", &rule.ReturnsInterfaceTypeRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{
			"reportAll": true,
		}},
	})
	testRule(t, "returns_interface_type_report_all_search", &rule.ReturnsInterfaceTypeRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{
			"reportAll":      true,
			"searchingNames": []any{"error"},
		}},
	})
}
