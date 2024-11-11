package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestCognitiveComplexity(t *testing.T) {
	testRule(t, "cognitive_complexity", &rule.CognitiveComplexityRule{}, &lint.RuleConfig{
		Arguments: []any{int64(0)},
	})
}
