package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestContextAsArgument(t *testing.T) {
	testRule(t, "context-as-argument", &rule.ContextAsArgumentRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"allowTypesBefore": "AllowedBeforeType,AllowedBeforeStruct,*AllowedBeforePtrStruct,*testing.T",
			},
		},
	},
	)
}
