package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
)

type TestFileFilterRule struct {
	WasApplyed bool
}

var _ lint.Rule = (*TestFileFilterRule)(nil)

func (*TestFileFilterRule) Name() string { return "test-file-filter" }
func (tfr *TestFileFilterRule) Apply(*lint.File, lint.Arguments) []lint.Failure {
	tfr.WasApplyed = true
	return nil
}

func TestFileExcludeFilterAtRuleLevel(t *testing.T) {
	t.Run("is called if no excludes", func(t *testing.T) {
		rule := &TestFileFilterRule{}
		testRule(t, "file-to-exclude", rule, &lint.RuleConfig{})
		if !rule.WasApplyed {
			t.Fatal("should call rule if no excludes")
		}
	})
	t.Run("is called if exclude not match", func(t *testing.T) {
		rule := &TestFileFilterRule{}
		cfg := &lint.RuleConfig{Exclude: []string{"no-matched.go"}}
		cfg.Initialize()
		testRule(t, "file-to-exclude", rule, cfg)
		if !rule.WasApplyed {
			t.Fatal("should call rule if no excludes")
		}
	})

	t.Run("not called if exclude not match", func(t *testing.T) {
		rule := &TestFileFilterRule{}
		cfg := &lint.RuleConfig{Exclude: []string{"file-to-exclude.go"}}
		cfg.Initialize()
		testRule(t, "file-to-exclude", rule, cfg)
		if rule.WasApplyed {
			t.Fatal("should not call rule if excluded")
		}
	})
}
