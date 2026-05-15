package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestReturnsInterfaceTypeDefault(t *testing.T) {
	testRule(t, "returns_interface_type_default", &rule.ReturnsInterfaceTypeRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{
			"ignoredNames": []any{
				"fixtures.DummyResults",
			},
		}},
	})
}
