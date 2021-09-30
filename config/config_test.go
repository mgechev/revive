package config

import (
	"reflect"
	"strings"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestGetConfig(t *testing.T) {
	tt := map[string]struct {
		confPath   string
		wantConfig *lint.Config
		wantError  string
	}{
		"non-reg issue #470": {
			confPath:  "testdata/issue-470.toml",
			wantError: "",
		},
		"unknown file": {
			confPath:  "unknown",
			wantError: "cannot read the config file",
		},
		"malformed file": {
			confPath:  "testdata/malformed.toml",
			wantError: "cannot parse the config file",
		},
		"default config": {
			wantConfig: defaultConfig(),
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			cfg, err := GetConfig(tc.confPath)
			switch {
			case err != nil && tc.wantError == "":
				t.Fatalf("Unexpected error\n\t%v", err)
			case err != nil && !strings.Contains(err.Error(), tc.wantError):
				t.Fatalf("Expected error\n\t%q\ngot:\n\t%v", tc.wantError, err)
			case tc.wantConfig != nil && reflect.DeepEqual(cfg, tc.wantConfig):
				t.Fatalf("Expected config\n\t%+v\ngot:\n\t%+v", tc.wantConfig, cfg)
			}

		})
	}
}

func TestGetLintingRules(t *testing.T) {
	tt := map[string]struct {
		confPath       string
		wantRulesCount int
	}{
		"no rules": {
			confPath:       "testdata/noRules.toml",
			wantRulesCount: 0,
		},
		"enableAllRules without disabled rules": {
			confPath:       "testdata/enableAll.toml",
			wantRulesCount: len(allRules),
		},
		"enableAllRules with 2 disabled rules": {
			confPath:       "testdata/enableAllBut2.toml",
			wantRulesCount: len(allRules) - 2,
		},
		"enable 2 rules": {
			confPath:       "testdata/enable2.toml",
			wantRulesCount: 2,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			cfg, err := GetConfig(tc.confPath)
			if err != nil {
				t.Fatalf("Unexpected error while loading conf: %v", err)
			}
			rules, err := GetLintingRules(cfg)
			switch {
			case err != nil:
				t.Fatalf("Unexpected error\n\t%v", err)
			case len(rules) != tc.wantRulesCount:
				t.Fatalf("Expected %v enabled linting rules got: %v", tc.wantRulesCount, len(rules))
			}

		})
	}
}

func TestGetGlobalSeverity(t *testing.T) {
	tt := map[string]struct {
		confPath               string
		wantGlobalSeverity     string
		particularRule         lint.Rule
		wantParticularSeverity string
	}{
		"enable 2 rules with one specific severity": {
			confPath:               "testdata/enable2OneSpecificSeverity.toml",
			wantGlobalSeverity:     "warning",
			particularRule:         &rule.CyclomaticRule{},
			wantParticularSeverity: "error",
		},
		"enableAllRules with one specific severity": {
			confPath:               "testdata/enableAllOneSpecificSeverity.toml",
			wantGlobalSeverity:     "error",
			particularRule:         &rule.DeepExitRule{},
			wantParticularSeverity: "warning",
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			cfg, err := GetConfig(tc.confPath)
			if err != nil {
				t.Fatalf("Unexpected error while loading conf: %v", err)
			}
			rules, err := GetLintingRules(cfg)
			if err != nil {
				t.Fatalf("Unexpected error while loading conf: %v", err)
			}
			for _, r := range rules {
				ruleName := r.Name()
				ruleCfg := cfg.Rules[ruleName]
				ruleSeverity := string(ruleCfg.Severity)
				switch ruleName {
				case tc.particularRule.Name():
					if tc.wantParticularSeverity != ruleSeverity {
						t.Fatalf("Expected Severity %v for rule %v, got %v", tc.wantParticularSeverity, ruleName, ruleSeverity)
					}
				default:
					if tc.wantGlobalSeverity != ruleSeverity {
						t.Fatalf("Expected Severity %v for rule %v, got %v", tc.wantGlobalSeverity, ruleName, ruleSeverity)
					}
				}
			}
		})
	}
}
