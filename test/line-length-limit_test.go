package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestLineLengthLimit(t *testing.T) {
	testRule(t, "line-length-limit", &rule.LineLengthLimitRule{}, &lint.RuleConfig{
		Arguments: []interface{}{int64(100)},
	})
}
