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
		confPath       string
		wantConfig     *lint.Config
		wantError      string
		wantConfidence float64
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
			wantConfig: func() *lint.Config {
				c := defaultConfig()
				normalizeConfig(c)
				return c
			}(),
			wantConfidence: defaultConfidence,
		},
		"config from file issue #585": {
			confPath:       "testdata/issue-585.toml",
			wantConfidence: 0.0,
		},
		"config from file default confidence issue #585": {
			confPath:       "testdata/issue-585-defaultConfidence.toml",
			wantConfidence: defaultConfidence,
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
			case tc.wantConfig != nil && !reflect.DeepEqual(cfg, tc.wantConfig):
				t.Fatalf("Expected config\n\t%+v\ngot:\n\t%+v", tc.wantConfig, cfg)
			case tc.wantConfig != nil && tc.wantConfidence != cfg.Confidence:
				t.Fatalf("Expected confidence\n\t%+v\ngot:\n\t%+v", tc.wantConfidence, cfg.Confidence)
			}
		})
	}

	t.Run("rule-level file filter excludes", func(t *testing.T) {
		cfg, err := GetConfig("testdata/rule-level-exclude-850.toml")
		if err != nil {
			t.Fatal("should be valid config")
		}
		r1 := cfg.Rules["r1"]
		if len(r1.Exclude) > 0 {
			t.Fatal("r1 should have empty excludes")
		}
		r2 := cfg.Rules["r2"]
		if len(r2.Exclude) != 1 {
			t.Fatal("r2 should have exclude set")
		}
		if !r2.MustExclude("some/file.go") {
			t.Fatal("r2 should be initialized and exclude some/file.go")
		}
		if r2.MustExclude("some/any-other.go") {
			t.Fatal("r2 should not exclude some/any-other.go")
		}
	})
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
			rules, err := GetLintingRules(cfg, []lint.Rule{})
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
			rules, err := GetLintingRules(cfg, []lint.Rule{})
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
