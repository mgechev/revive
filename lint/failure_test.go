package lint_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestFailureSeverityFor(t *testing.T) {
	for name, tc := range map[string]struct {
		config  *lint.Config
		failure lint.Failure
		want    lint.Severity
	}{
		"rule configured as error": {
			config:  &lint.Config{Rules: lint.RulesConfig{"r": {Severity: lint.SeverityError}}},
			failure: lint.Failure{RuleName: "r"},
			want:    lint.SeverityError,
		},
		"rule configured as warning": {
			config:  &lint.Config{Rules: lint.RulesConfig{"r": {Severity: lint.SeverityWarning}}},
			failure: lint.Failure{RuleName: "r"},
			want:    lint.SeverityWarning,
		},
		"directive configured as error": {
			config:  &lint.Config{Directives: lint.DirectivesConfig{"d": {Severity: lint.SeverityError}}},
			failure: lint.Failure{RuleName: "d"},
			want:    lint.SeverityError,
		},
		"rule without severity defaults to warning": {
			config:  &lint.Config{Rules: lint.RulesConfig{"r": {}}},
			failure: lint.Failure{RuleName: "r"},
			want:    lint.SeverityWarning,
		},
		"rule not in config defaults to warning": {
			config:  &lint.Config{},
			failure: lint.Failure{RuleName: "unknown"},
			want:    lint.SeverityWarning,
		},
		"nil config defaults to warning": {
			config:  nil,
			failure: lint.Failure{RuleName: "r"},
			want:    lint.SeverityWarning,
		},
	} {
		t.Run(name, func(t *testing.T) {
			if got := tc.failure.SeverityFor(tc.config); got != tc.want {
				t.Errorf("SeverityFor: expected %q, got %q", tc.want, got)
			}
		})
	}
}
