package revivelib

import (
	"testing"

	"github.com/mgechev/revive/config"
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

	if revive.lintingRules == nil || len(revive.lintingRules) == 0 {
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

type mockRule struct {
}

func (*mockRule) Name() string {
	return "mock-rule"
}

func (*mockRule) Apply(_ *lint.File, _ lint.Arguments) []lint.Failure {
	return nil
}

func getMockRevive(t *testing.T) *Revive {
	t.Helper()

	conf, err := config.GetConfig("../defaults.toml")
	if err != nil {
		t.Fatal(err)
	}

	revive, err := New(
		conf,
		true,
		2048,
		NewExtraRule(&mockRule{}, lint.RuleConfig{}),
	)
	if err != nil {
		t.Fatal(err.Error())
	}

	return revive
}
