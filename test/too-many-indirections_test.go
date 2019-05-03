package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestTooManyIndirectionsRule(t *testing.T) {
	args := []interface{}{int64(2)}

	testRule(t, "too-many-indirections", &rule.TooManyIndirectionsRule{}, &lint.RuleConfig{Arguments: args})
}
