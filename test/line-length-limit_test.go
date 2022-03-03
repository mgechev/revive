package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/lint"
	"github.com/deepsourcelabs/revive/rule"
)

func TestLineLengthLimit(t *testing.T) {
	testRule(t, "line-length-limit", &rule.LineLengthLimitRule{}, &lint.RuleConfig{
		Arguments: []interface{}{int64(100)},
	})
}
