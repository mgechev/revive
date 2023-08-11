package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestUnusedParam(t *testing.T) {
	testRule(t, "unused-param", &rule.UnusedParamRule{})
	testRule(t, "unused-param-custom-regex", &rule.UnusedParamRule{}, &lint.RuleConfig{Arguments: []interface{}{
		map[string]interface{}{"allowRegex": "^xxx"},
	}})
}

func BenchmarkUnusedParam(b *testing.B) {
	var t *testing.T
	for i := 0; i <= b.N; i++ {
		testRule(t, "unused-param", &rule.UnusedParamRule{})
	}
}
