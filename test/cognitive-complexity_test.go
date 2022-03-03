package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/lint"
	"github.com/deepsourcelabs/revive/rule"
)

func TestCognitiveComplexity(t *testing.T) {
	testRule(t, "cognitive-complexity", &rule.CognitiveComplexityRule{}, &lint.RuleConfig{
		Arguments: []interface{}{int64(0)},
	})
}
