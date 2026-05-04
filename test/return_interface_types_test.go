package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestReturnInterfaceTypesDefault(t *testing.T) {
	testRule(t, "return_interface_types_default", &rule.ReturnInterfaceTypesRule{}, &lint.RuleConfig{})
}

func TestReturnInterfaceTypesSkip(t *testing.T) {
	testRule(t, "return_interface_types_skip", &rule.ReturnInterfaceTypesRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{
			"stopOnFirst":             true,
			"userDefinedIgnoredNames": []any{"fixtures.DummyResults"},
		}},
	})
}
