package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestContextAsArgumentDefault(t *testing.T) {
	testRule(t, "context_as_argument_default", &rule.ContextAsArgumentRule{}, &lint.RuleConfig{})
}

func TestContextAsArgument(t *testing.T) {
	testRule(t, "context_as_argument", &rule.ContextAsArgumentRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"allowTypesBefore": "AllowedBeforeType,AllowedBeforeStruct,*AllowedBeforePtrStruct,*testing.T",
			},
		},
	})
	testRule(t, "context_as_argument", &rule.ContextAsArgumentRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"allow-types-before": "AllowedBeforeType,AllowedBeforeStruct,*AllowedBeforePtrStruct,*testing.T",
			},
		},
	})
}
