package lint

import "testing"

func TestFailureSeverityFor(t *testing.T) {
	for name, tc := range map[string]struct {
		config  *Config
		failure Failure
		want    Severity
	}{
		"rule configured as error": {
			config:  &Config{Rules: RulesConfig{"r": {Severity: SeverityError}}},
			failure: Failure{RuleName: "r"},
			want:    SeverityError,
		},
		"rule configured as warning": {
			config:  &Config{Rules: RulesConfig{"r": {Severity: SeverityWarning}}},
			failure: Failure{RuleName: "r"},
			want:    SeverityWarning,
		},
		"directive configured as error": {
			config:  &Config{Directives: DirectivesConfig{"d": {Severity: SeverityError}}},
			failure: Failure{RuleName: "d"},
			want:    SeverityError,
		},
		"rule without severity defaults to warning": {
			config:  &Config{Rules: RulesConfig{"r": {}}},
			failure: Failure{RuleName: "r"},
			want:    SeverityWarning,
		},
		"rule not in config defaults to warning": {
			config:  &Config{},
			failure: Failure{RuleName: "unknown"},
			want:    SeverityWarning,
		},
		"nil config defaults to warning": {
			config:  nil,
			failure: Failure{RuleName: "r"},
			want:    SeverityWarning,
		},
	} {
		t.Run(name, func(t *testing.T) {
			if got := tc.failure.SeverityFor(tc.config); got != tc.want {
				t.Errorf("SeverityFor: expected %q, got %q", tc.want, got)
			}
		})
	}
}
