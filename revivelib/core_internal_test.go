package revivelib

import (
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestReviveCreateInstance(t *testing.T) {
	revive := getMockRevive(t)

	if revive.config == nil {
		t.Fatal("Could not load config.")
	}

	if revive.maxOpenFiles != 2048 {
		t.Fatal("Expected MaxOpenFiles to be 2048")
	}

	if len(revive.lintingRules) == 0 {
		t.Fatal("Linting rules not loaded.")
	}

	rules := map[string]lint.Rule{}
	for _, rule := range revive.lintingRules {
		rules[rule.Name()] = rule
	}

	if _, ok := rules["mock-rule"]; !ok {
		t.Fatal("Didn't load mock rule.")
	}

	if revive.config.ErrorCode != 1 && revive.config.WarningCode != 1 {
		t.Fatal("Didn't set the codes in the config instance.")
	}
}
