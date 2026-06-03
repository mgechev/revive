package test_test

import (
	"os"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestSpecifyDisableRule(t *testing.T) {
	l := lint.New(os.ReadFile, 0)

	filePath := "../testdata/specify_disable_rules.go"
	rules := []lint.Rule{&rule.ExportedRule{}}
	config := lint.Config{
		Rules: map[string]lint.RuleConfig{
			"exported": {},
		},
		Directives: lint.DirectivesConfig{
			"specify-disable-rule": {},
		},
	}

	ps, err := l.Lint([][]string{{filePath}}, rules, config)
	if err != nil {
		t.Fatalf("Linting %s: %v", filePath, err)
	}

	var failures []lint.Failure
	for f := range ps {
		failures = append(failures, f)
	}

	// Expect 3 failures: naked disable, naked disable-line, naked disable-next-line
	expectedCount := 3
	if len(failures) != expectedCount {
		t.Fatalf("Expected %d failures, got %d", expectedCount, len(failures))
	}

	for _, f := range failures {
		if f.RuleName != "specify-disable-rule" {
			t.Errorf("Expected rule name 'specify-disable-rule', got '%s'", f.RuleName)
		}
		if f.Failure != "rule name for lint disabling not found" {
			t.Errorf("Expected failure message 'rule name for lint disabling not found', got '%s'", f.Failure)
		}
	}
}

func TestSpecifyDisableRulesDisabled(t *testing.T) {
	// Without the directive config, naked disable should be allowed
	l := lint.New(os.ReadFile, 0)

	filePath := "../testdata/specify_disable_rules.go"
	rules := []lint.Rule{&rule.ExportedRule{}}
	config := lint.Config{
		Rules: map[string]lint.RuleConfig{
			"exported": {},
		},
	}

	ps, err := l.Lint([][]string{{filePath}}, rules, config)
	if err != nil {
		t.Fatalf("Linting %s: %v", filePath, err)
	}

	var failures []lint.Failure
	for f := range ps {
		failures = append(failures, f)
	}

	// No failures from the directive check - naked disables are allowed
	for _, f := range failures {
		if f.RuleName == "specify-disable-rule" {
			t.Error("specify-disable-rule failure should not appear when directive is not configured")
		}
	}
}
