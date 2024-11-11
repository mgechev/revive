package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
)

type TestFileFilterRule struct {
	WasApplied bool
}

var _ lint.Rule = (*TestFileFilterRule)(nil)

func (*TestFileFilterRule) Name() string { return "test-file-filter" }
func (tfr *TestFileFilterRule) Apply(*lint.File, lint.Arguments) []lint.Failure {
	tfr.WasApplied = true
	return nil
}

func TestFileExcludeFilterAtRuleLevel(t *testing.T) {
	t.Run("is called if no excludes", func(t *testing.T) {
		rule := &TestFileFilterRule{}
		testRule(t, "file_to_exclude", rule, &lint.RuleConfig{})
		if !rule.WasApplied {
			t.Fatal("should call rule if no excludes")
		}
	})
	t.Run("is called if exclude not match", func(t *testing.T) {
		rule := &TestFileFilterRule{}
		cfg := &lint.RuleConfig{Exclude: []string{"no_matched.go"}}
		cfg.Initialize()
		testRule(t, "file_to_exclude", rule, cfg)
		if !rule.WasApplied {
			t.Fatal("should call rule if no excludes")
		}
	})

	t.Run("not called if exclude not match", func(t *testing.T) {
		rule := &TestFileFilterRule{}
		cfg := &lint.RuleConfig{Exclude: []string{"../testdata/file_to_exclude.go"}}
		cfg.Initialize()
		testRule(t, "file_to_exclude", rule, cfg)
		if rule.WasApplied {
			t.Fatal("should not call rule if excluded")
		}
	})
}
