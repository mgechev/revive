package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestMaxPublicStructs(t *testing.T) {
	testRule(t, "max_public_structs", &rule.MaxPublicStructsRule{}, &lint.RuleConfig{
		Arguments: []any{int64(1)},
	})
}

func TestMaxPublicStructsDefaultConfig(t *testing.T) {
	testRule(t, "max_public_structs_ok", &rule.MaxPublicStructsRule{}, &lint.RuleConfig{})
}
