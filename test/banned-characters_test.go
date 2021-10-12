package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestBannedCharacters(t *testing.T) {
	testRule(t, "banned-characters", &rule.BannedCharsRule{}, &lint.RuleConfig{
		Arguments: []interface{}{"Ω", "Σ", "σ"},
	})
}

