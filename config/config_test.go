package config

import (
	"path/filepath"
	"strings"
	"testing"

	goversion "github.com/hashicorp/go-version"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestGetConfig(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		for name, tc := range map[string]struct {
			confPath   string
			wantConfig lint.Config
		}{
			"default config": {
				wantConfig: func() lint.Config {
					c := defaultConfig()
					normalizeConfig(c)
					return *c
				}(),
			},
			"non-reg issue #470": {
				confPath: "issue-470.toml",
				wantConfig: lint.Config{
					Confidence: 0.8,
					Severity:   lint.SeverityWarning,
					Rules: lint.RulesConfig{
						"add-constant": {
							Severity: lint.SeverityWarning,
							Arguments: lint.Arguments{
								map[string]any{
									"maxLitCount": "3",
									"allowStrs":   `"`,
									"allowFloats": "0.0,1.0,1.,2.0,2.",
									"allowInts":   "0,1,2",
								},
							},
						},
					},
				},
			},
			"config from file issue #585": {
				confPath: "issue-585.toml",
				wantConfig: lint.Config{
					Confidence: 0.0,
					Severity:   lint.SeverityWarning,
				},
			},
			"config from file default confidence issue #585": {
				confPath: "issue-585-defaultConfidence.toml",
				wantConfig: lint.Config{
					Confidence: 0.8,
					Severity:   lint.SeverityWarning,
				},
			},
			"config from file goVersion": {
				confPath: "goVersion.toml",
				wantConfig: lint.Config{
					Confidence: 0.8,
					GoVersion:  goversion.Must(goversion.NewSemver("1.20.0")),
				},
			},
			"config from file ignoreGeneratedHeader": {
				confPath: "ignoreGeneratedHeader.toml",
				wantConfig: lint.Config{
					Confidence:            0.8,
					IgnoreGeneratedHeader: true,
				},
			},
			"config from file enableDefault": {
				confPath: "enableDefault.toml",
				wantConfig: lint.Config{
					Confidence:            0.8,
					IgnoreGeneratedHeader: false,
					EnableDefaultRules:    true,
					Rules: lint.RulesConfig{
						"blank-imports":        {},
						"context-as-argument":  {},
						"context-keys-type":    {},
						"dot-imports":          {},
						"empty-block":          {},
						"error-naming":         {},
						"error-return":         {},
						"error-strings":        {},
						"errorf":               {},
						"exported":             {},
						"increment-decrement":  {},
						"indent-error-flow":    {},
						"package-comments":     {},
						"range":                {},
						"receiver-naming":      {},
						"redefines-builtin-id": {},
						"superfluous-else":     {},
						"time-naming":          {},
						"unexported-return":    {},
						"unreachable-code":     {},
						"unused-parameter":     {},
						"var-declaration":      {},
						"var-naming":           {},
					},
				},
			},
		} {
			t.Run(name, func(t *testing.T) {
				var cfgPath string
				if tc.confPath != "" {
					cfgPath = filepath.Join("testdata", tc.confPath)
				}

				cfg, err := GetConfig(cfgPath)

				if err != nil {
					t.Fatalf("Unexpected error %v", err)
				}
				if cfg.IgnoreGeneratedHeader != tc.wantConfig.IgnoreGeneratedHeader {
					t.Errorf("IgnoreGeneratedHeader: expected %v, got %v", tc.wantConfig.IgnoreGeneratedHeader, cfg.IgnoreGeneratedHeader)
				}
				if cfg.Confidence != tc.wantConfig.Confidence {
					t.Errorf("Confidence: expected %v, got %v", tc.wantConfig.Confidence, cfg.Confidence)
				}
				if cfg.Severity != tc.wantConfig.Severity {
					t.Errorf("Severity: expected %v, got %v", tc.wantConfig.Severity, cfg.Severity)
				}
				if cfg.EnableAllRules != tc.wantConfig.EnableAllRules {
					t.Errorf("EnableAllRules: expected %v, got %v", tc.wantConfig.EnableAllRules, cfg.EnableAllRules)
				}
				if cfg.EnableDefaultRules != tc.wantConfig.EnableDefaultRules {
					t.Errorf("EnableDefaultRules: expected %v, got %v", tc.wantConfig.EnableDefaultRules, cfg.EnableDefaultRules)
				}
				if cfg.ErrorCode != tc.wantConfig.ErrorCode {
					t.Errorf("ErrorCode: expected %v, got %v", tc.wantConfig.ErrorCode, cfg.ErrorCode)
				}
				if cfg.WarningCode != tc.wantConfig.WarningCode {
					t.Errorf("WarningCode: expected %v, got %v", tc.wantConfig.WarningCode, cfg.WarningCode)
				}
				if !tc.wantConfig.GoVersion.Equal(cfg.GoVersion) {
					t.Errorf("GoVersion: expected %v, got %v", tc.wantConfig.GoVersion, cfg.GoVersion)
				}

				if len(cfg.Exclude) != len(tc.wantConfig.Exclude) {
					t.Errorf("Exclude length: expected %v, got %v", len(tc.wantConfig.Exclude), len(cfg.Exclude))
				} else {
					for i, exclude := range tc.wantConfig.Exclude {
						if cfg.Exclude[i] != exclude {
							t.Errorf("Exclude[%d]: expected %v, got %v", i, exclude, cfg.Exclude[i])
						}
					}
				}

				if len(cfg.Rules) != len(tc.wantConfig.Rules) {
					t.Errorf("Rules count: expected %v, got %v", len(tc.wantConfig.Rules), len(cfg.Rules))
				}
				for ruleName, wantRule := range tc.wantConfig.Rules {
					gotRule, exists := cfg.Rules[ruleName]
					if !exists {
						t.Errorf("Rule %q: expected to exist, but not found", ruleName)
						continue
					}
					if gotRule.Disabled != wantRule.Disabled {
						t.Errorf("Rule %q Disabled: expected %v, got %v", ruleName, wantRule.Disabled, gotRule.Disabled)
					}
					if gotRule.Severity != wantRule.Severity {
						t.Errorf("Rule %q Severity: expected %v, got %v", ruleName, wantRule.Severity, gotRule.Severity)
					}
					if len(gotRule.Arguments) != len(wantRule.Arguments) {
						t.Errorf("Rule %q Arguments length: expected %v, got %v", ruleName, len(wantRule.Arguments), len(gotRule.Arguments))
					}
					if len(gotRule.Exclude) != len(wantRule.Exclude) {
						t.Errorf("Rule %q Exclude length: expected %v, got %v", ruleName, len(wantRule.Exclude), len(gotRule.Exclude))
					} else {
						for i, wantExclude := range wantRule.Exclude {
							if gotRule.Exclude[i] != wantExclude {
								t.Errorf("Rule %q Exclude[%d]: expected %v, got %v", ruleName, i, wantExclude, gotRule.Exclude[i])
							}
						}
					}
				}
				// Check for unexpected rules in actual config
				for ruleName := range cfg.Rules {
					if _, exists := tc.wantConfig.Rules[ruleName]; !exists {
						t.Errorf("Rule %q: found in actual config but not expected", ruleName)
					}
				}

				if len(cfg.Directives) != len(tc.wantConfig.Directives) {
					t.Errorf("Directives count: expected %v, got %v", len(tc.wantConfig.Directives), len(cfg.Directives))
				}
				for directiveName, wantDirective := range tc.wantConfig.Directives {
					gotDirective, exists := cfg.Directives[directiveName]
					if !exists {
						t.Errorf("Directive %q: expected to exist, but not found", directiveName)
						continue
					}
					if gotDirective.Severity != wantDirective.Severity {
						t.Errorf("Directive %q Severity: expected %v, got %v", directiveName, wantDirective.Severity, gotDirective.Severity)
					}
				}
				// Check for unexpected directives in actual config
				for directiveName := range cfg.Directives {
					if _, exists := tc.wantConfig.Directives[directiveName]; !exists {
						t.Errorf("Directive %q: found in actual config but not expected", directiveName)
					}
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
	})

	t.Run("failure", func(t *testing.T) {
		for name, tc := range map[string]struct {
			confPath  string
			wantError string
		}{
			"unknown file": {
				confPath:  "unknown",
				wantError: "cannot read the config file",
			},
			"malformed file": {
				confPath:  "malformed.toml",
				wantError: "cannot parse the config file",
			},
		} {
			t.Run(name, func(t *testing.T) {
				_, err := GetConfig(filepath.Join("testdata", tc.confPath))

				if err != nil && !strings.Contains(err.Error(), tc.wantError) {
					t.Errorf("Unexpected error: want %q, got: %q", tc.wantError, err)
				}
			})
		}
	})
}

func TestGetLintingRules(t *testing.T) {
	tt := map[string]struct {
		confPath       string
		wantRulesCount int
		wantErr        string
	}{
		"no rules": {
			confPath:       "noRules.toml",
			wantRulesCount: 0,
		},
		"enableAllRules without disabled rules": {
			confPath:       "enableAll.toml",
			wantRulesCount: len(allRules),
		},
		"enableAllRules with 2 disabled rules": {
			confPath:       "enableAllBut2.toml",
			wantRulesCount: len(allRules) - 2,
		},
		"enableDefaultRules without disabled rules": {
			confPath:       "enableDefault.toml",
			wantRulesCount: len(defaultRules),
		},
		"enableDefaultRules with 2 disabled rules": {
			confPath:       "enableDefaultBut2.toml",
			wantRulesCount: len(defaultRules) - 2,
		},
		"enableDefaultRules plus 1 non-default rule": {
			confPath:       "enableDefaultPlus1.toml",
			wantRulesCount: len(defaultRules) + 1,
		},
		"enableAllRules and enableDefaultRules both set": {
			confPath:       "enableAllAndDefault.toml",
			wantRulesCount: len(allRules),
		},
		"enableDefaultRules plus rule already in defaults": {
			confPath:       "enableDefaultPlusDefaultRule.toml",
			wantRulesCount: len(defaultRules),
		},
		"enableAllRules plus rule already in all": {
			confPath:       "enableAllWithRule.toml",
			wantRulesCount: len(allRules),
		},
		"enable 2 rules": {
			confPath:       "enable2.toml",
			wantRulesCount: 2,
		},
		"var-naming configure error": {
			confPath: "varNamingConfigureError.toml",
			wantErr:  `cannot configure rule: "var-naming": invalid argument to the var-naming rule. Expecting a allowlist of type slice with initialisms, got string`,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			cfg, err := GetConfig(filepath.Join("testdata", tc.confPath))
			if err != nil {
				t.Fatalf("Unexpected error while loading conf: %v", err)
			}
			rules, err := GetLintingRules(cfg, []lint.Rule{})
			if tc.wantErr != "" {
				if err == nil || err.Error() != tc.wantErr {
					t.Fatalf("Expected error %q, got %q", tc.wantErr, err)
				}
				return
			}

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

func TestGetFormatter(t *testing.T) {
	t.Run("default formatter", func(t *testing.T) {
		formatter, err := GetFormatter("")
		if err != nil {
			t.Fatalf("Unexpected error %q", err)
		}
		if formatter == nil || formatter.Name() != "default" {
			t.Errorf("Expected formatter %q, got %v", "default", formatter)
		}
	})
	t.Run("unknown formatter", func(t *testing.T) {
		_, err := GetFormatter("unknown")
		if err == nil || err.Error() != "unknown formatter unknown" {
			t.Errorf("Expected error %q, got: %q", "unknown formatter unknown", err)
		}
	})
	t.Run("checkstyle formatter", func(t *testing.T) {
		formatter, err := GetFormatter("checkstyle")
		if err != nil {
			t.Fatalf("Unexpected error: %q", err)
		}
		if formatter == nil || formatter.Name() != "checkstyle" {
			t.Errorf("Expected formatter %q, got %v", "checkstyle", formatter)
		}
	})
}
