package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestBannedCharactersDefault(t *testing.T) {
	testRule(t, "banned_characters_default", &rule.BannedCharsRule{}, &lint.RuleConfig{})
}

// Test banned characters in a const, var and func names.
// One banned character is in the comment and should not be checked.
// One banned character from the list is not present in the fixture file.
func TestBannedCharacters(t *testing.T) {
	testRule(t, "banned_characters", &rule.BannedCharsRule{}, &lint.RuleConfig{
		Arguments: []any{"Ω", "Σ", "σ", "1"},
	})
}
