package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestContextAsArgument(t *testing.T) {
	testRule(t, "context-as-argument", &rule.ContextAsArgumentRule{}, &lint.RuleConfig{
		Arguments: []interface{}{
			map[string]interface{}{
				"allowTypesBefore": []string{
					"AllowedBeforeType",
					"AllowedBeforeStruct",
					"*AllowedBeforePtrStruct",
					"*testing.T",
				},
			},
		},
	})
}
