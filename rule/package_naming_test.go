package rule

import (
	"errors"
	"reflect"
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestPackageNamingRule_Configure(t *testing.T) {
	tests := []struct {
		name                           string
		arguments                      lint.Arguments
		wantErr                        error
		wantSkipConventionNameCheck    bool
		wantSkipTopLevelCheck          bool
		wantSkipDefaultBadNameCheck    bool
		wantUserDefinedBadNames        map[string]struct{}
		wantSkipCollisionWithCommonStd bool
		wantCheckCollisionWithAllStd   bool
	}{
		{
			name:      "no arguments",
			arguments: lint.Arguments{},
			wantErr:   nil,
		},
		{
			name: "valid arguments - camelCase",
			arguments: lint.Arguments{
				map[string]any{
					"skipConventionNameCheck":    true,
					"skipTopLevelCheck":          true,
					"skipDefaultBadNameCheck":    true,
					"userDefinedBadNames":        []any{"helpers", "models"},
					"skipCollisionWithCommonStd": true,
					"checkCollisionWithAllStd":   true,
				},
			},
			wantErr:                        nil,
			wantSkipConventionNameCheck:    true,
			wantSkipTopLevelCheck:          true,
			wantSkipDefaultBadNameCheck:    true,
			wantUserDefinedBadNames:        map[string]struct{}{"helpers": {}, "models": {}},
			wantSkipCollisionWithCommonStd: true,
			wantCheckCollisionWithAllStd:   true,
		},
		{
			name: "valid arguments - lowercase",
			arguments: lint.Arguments{
				map[string]any{
					"skipconventionnamecheck":    true,
					"skiptoplevelcheck":          true,
					"skipdefaultbadnamecheck":    true,
					"userdefinedbadnames":        []any{"helpers", "models"},
					"skipcollisionwithcommonstd": true,
					"checkcollisionwithallstd":   true,
				},
			},
			wantErr:                        nil,
			wantSkipConventionNameCheck:    true,
			wantSkipTopLevelCheck:          true,
			wantSkipDefaultBadNameCheck:    true,
			wantUserDefinedBadNames:        map[string]struct{}{"helpers": {}, "models": {}},
			wantSkipCollisionWithCommonStd: true,
			wantCheckCollisionWithAllStd:   true,
		},
		{
			name: "valid arguments - kebab-case",
			arguments: lint.Arguments{
				map[string]any{
					"skip-convention-name-check":     true,
					"skip-top-level-check":           true,
					"skip-default-bad-name-check":    true,
					"user-defined-bad-names":         []any{"helpers", "models"},
					"skip-collision-with-common-std": true,
					"check-collision-with-all-std":   true,
				},
			},
			wantErr:                        nil,
			wantSkipConventionNameCheck:    true,
			wantSkipTopLevelCheck:          true,
			wantSkipDefaultBadNameCheck:    true,
			wantUserDefinedBadNames:        map[string]struct{}{"helpers": {}, "models": {}},
			wantSkipCollisionWithCommonStd: true,
			wantCheckCollisionWithAllStd:   true,
		},
		{
			name: "partial arguments",
			arguments: lint.Arguments{
				map[string]any{
					"skip-convention-name-check": true,
					"user-defined-bad-names":     []any{"custom"},
				},
			},
			wantErr:                     nil,
			wantSkipConventionNameCheck: true,
			wantUserDefinedBadNames:     map[string]struct{}{"custom": {}},
		},
		{
			name: "invalid argument type",
			arguments: lint.Arguments{
				"invalid-arg",
			},
			wantErr: errors.New("invalid argument to the package-naming rule: expecting a k,v map, but got string"),
		},
		{
			name: "invalid userDefinedBadNames type",
			arguments: lint.Arguments{
				map[string]any{
					"user-defined-bad-names": "invalid-type",
				},
			},
			wantErr: errors.New("invalid argument to the package-naming rule: expecting userDefinedBadNames of type slice of strings, but got string"),
		},
		{
			name: "invalid userDefinedBadNames element type",
			arguments: lint.Arguments{
				map[string]any{
					"user-defined-bad-names": []any{"helpers", 123},
				},
			},
			wantErr: errors.New("invalid argument to the package-naming rule: expecting element 1 of userDefinedBadNames to be a string, but got 123(int)"),
		},
		{
			name: "empty string in userDefinedBadNames",
			arguments: lint.Arguments{
				map[string]any{
					"user-defined-bad-names": []any{"helpers", ""},
				},
			},
			wantErr: errors.New("invalid argument to the package-naming rule: userDefinedBadNames cannot contain empty string (index 1)"),
		},
		{
			name: "userDefinedBadNames with uppercase",
			arguments: lint.Arguments{
				map[string]any{
					"user-defined-bad-names": []any{"HELPERS", "Models"},
				},
			},
			wantErr:                 nil,
			wantUserDefinedBadNames: map[string]struct{}{"helpers": {}, "models": {}},
		},
		{
			name: "valid conventionNameCheckRegex",
			arguments: lint.Arguments{
				map[string]any{
					"convention-name-check-regex": "^[a-z][a-z0-9_]*$",
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid conventionNameCheckRegex type",
			arguments: lint.Arguments{
				map[string]any{
					"convention-name-check-regex": 42,
				},
			},
			wantErr: errors.New("invalid argument to the package-naming rule: expecting conventionNameCheckRegex to be a string, but got int"),
		},
		{
			name: "invalid conventionNameCheckRegex pattern",
			arguments: lint.Arguments{
				map[string]any{
					"convention-name-check-regex": "[",
				},
			},
			wantErr: errors.New("invalid argument to the package-naming rule: invalid regex for conventionNameCheckRegex: error parsing regexp: missing closing ]: `[`"),
		},
		{
			name: "skipConventionNameCheck with conventionNameCheckRegex",
			arguments: lint.Arguments{
				map[string]any{
					"skip-convention-name-check":  true,
					"convention-name-check-regex": "^[a-z]+$",
				},
			},
			wantErr: errors.New("invalid configuration for package-naming rule: skipConventionNameCheck and overrideConventionNameCheck cannot be both set"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rule PackageNamingRule

			err := rule.Configure(tt.arguments)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("unexpected error: got = nil, want = %v", tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("unexpected error: got = %v, want = %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: got = %v, want = nil", err)
			}
			if rule.skipConventionNameCheck != tt.wantSkipConventionNameCheck {
				t.Errorf("unexpected skipConventionNameCheck: got = %v, want %v", rule.skipConventionNameCheck, tt.wantSkipConventionNameCheck)
			}
			if rule.skipTopLevelCheck != tt.wantSkipTopLevelCheck {
				t.Errorf("unexpected skipTopLevelCheck: got = %v, want %v", rule.skipTopLevelCheck, tt.wantSkipTopLevelCheck)
			}
			if rule.skipDefaultBadNameCheck != tt.wantSkipDefaultBadNameCheck {
				t.Errorf("unexpected skipDefaultBadNameCheck: got = %v, want %v", rule.skipDefaultBadNameCheck, tt.wantSkipDefaultBadNameCheck)
			}
			if !reflect.DeepEqual(rule.userDefinedBadNames, tt.wantUserDefinedBadNames) {
				t.Errorf("unexpected userDefinedBadNames: got = %v, want %v", rule.userDefinedBadNames, tt.wantUserDefinedBadNames)
			}
			if rule.skipCollisionWithCommonStd != tt.wantSkipCollisionWithCommonStd {
				t.Errorf("unexpected skipCollisionWithCommonStd: got = %v, want %v", rule.skipCollisionWithCommonStd, tt.wantSkipCollisionWithCommonStd)
			}
			if rule.checkCollisionWithAllStd != tt.wantCheckCollisionWithAllStd {
				t.Errorf("unexpected checkCollisionWithAllStd: got = %v, want %v", rule.checkCollisionWithAllStd, tt.wantCheckCollisionWithAllStd)
			}
			if tt.wantErr == nil && rule.conventionNameCheckRegex != nil {
				if rule.conventionNameCheckRegex.String() == "" {
					t.Error("unexpected conventionNameCheckRegex: got empty string")
				}
			}
		})
	}
}

func TestPackageNamingRule_Configure_LoadStdPackages(t *testing.T) {
	t.Run("loads std packages when checkCollisionWithAllStd is true", func(t *testing.T) {
		var rule PackageNamingRule

		err := rule.Configure(lint.Arguments{
			map[string]any{
				"check-collision-with-all-std": true,
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: got = %v, want = nil", err)
		}

		if rule.allStdNames == nil {
			t.Fatal("expected allStdNames to be loaded, but got nil")
		}
		if len(rule.allStdNames) == 0 {
			t.Fatal("expected allStdNames to be populated, but got empty")
		}

		for _, pkg := range []string{"fmt", "http", "json", "rand", "io"} {
			if _, ok := rule.allStdNames[pkg]; !ok {
				t.Errorf("expected package %q to be loaded, but got empty", pkg)
			}
		}
	})

	t.Run("does not reload std packages on subsequent Configure calls", func(t *testing.T) {
		var rule PackageNamingRule

		err := rule.Configure(lint.Arguments{
			map[string]any{
				"check-collision-with-all-std": true,
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: got = %v, want = nil", err)
		}
		firstLoadLen := len(rule.allStdNames)
		err = rule.Configure(lint.Arguments{
			map[string]any{
				"check-collision-with-all-std": true,
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: got = %v, want = nil", err)
		}
		if len(rule.allStdNames) != firstLoadLen {
			t.Errorf("expected allStdNames to be loaded only once, but got different lengths: first %d, second %d", firstLoadLen, len(rule.allStdNames))
		}
	})

	t.Run("skips loading std packages when checkCollisionWithAllStd is false", func(t *testing.T) {
		var rule PackageNamingRule

		err := rule.Configure(lint.Arguments{
			map[string]any{
				"check-collision-with-all-std": false,
			},
		})
		if err != nil {
			t.Errorf("unexpected error: got = %v, want = nil", err)
		}

		if rule.allStdNames != nil {
			t.Errorf("expected allStdNames to be nil, but got %v", rule.allStdNames)
		}
	})
}
