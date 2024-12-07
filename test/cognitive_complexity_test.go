package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestCognitiveComplexityDefault(t *testing.T) {
	testRule(t, "cognitive_complexity_default", &rule.CognitiveComplexityRule{}, &lint.RuleConfig{})
}

func TestCognitiveComplexity(t *testing.T) {
	testRule(t, "cognitive_complexity", &rule.CognitiveComplexityRule{}, &lint.RuleConfig{
		Arguments: []any{int64(0)},
	})
}
