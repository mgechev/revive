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
		wantSkipConventionChecks       bool
		wantSkipTopLevelChecks         bool
		wantSkipDefaultBadNameChecks   bool
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
					"skipConventionChecks":       true,
					"skipTopLevelChecks":         true,
					"skipDefaultBadNameChecks":   true,
					"userDefinedBadNames":        []any{"helpers", "models"},
					"skipCollisionWithCommonStd": true,
					"checkCollisionWithAllStd":   true,
				},
			},
			wantErr:                        nil,
			wantSkipConventionChecks:       true,
			wantSkipTopLevelChecks:         true,
			wantSkipDefaultBadNameChecks:   true,
			wantUserDefinedBadNames:        map[string]struct{}{"helpers": {}, "models": {}},
			wantSkipCollisionWithCommonStd: true,
			wantCheckCollisionWithAllStd:   true,
		},
		{
			name: "valid arguments - lowercase",
			arguments: lint.Arguments{
				map[string]any{
					"skipconventionchecks":       true,
					"skiptoplevelchecks":         true,
					"skipdefaultbadnamechecks":   true,
					"userdefinedbadnames":        []any{"helpers", "models"},
					"skipcollisionwithcommonstd": true,
					"checkcollisionwithallstd":   true,
				},
			},
			wantErr:                        nil,
			wantSkipConventionChecks:       true,
			wantSkipTopLevelChecks:         true,
			wantSkipDefaultBadNameChecks:   true,
			wantUserDefinedBadNames:        map[string]struct{}{"helpers": {}, "models": {}},
			wantSkipCollisionWithCommonStd: true,
			wantCheckCollisionWithAllStd:   true,
		},
		{
			name: "valid arguments - kebab-case",
			arguments: lint.Arguments{
				map[string]any{
					"skip-convention-checks":         true,
					"skip-top-level-checks":          true,
					"skip-default-bad-name-checks":   true,
					"user-defined-bad-names":         []any{"helpers", "models"},
					"skip-collision-with-common-std": true,
					"check-collision-with-all-std":   true,
				},
			},
			wantErr:                        nil,
			wantSkipConventionChecks:       true,
			wantSkipTopLevelChecks:         true,
			wantSkipDefaultBadNameChecks:   true,
			wantUserDefinedBadNames:        map[string]struct{}{"helpers": {}, "models": {}},
			wantSkipCollisionWithCommonStd: true,
			wantCheckCollisionWithAllStd:   true,
		},
		{
			name: "partial arguments",
			arguments: lint.Arguments{
				map[string]any{
					"skip-convention-checks": true,
					"user-defined-bad-names": []any{"custom"},
				},
			},
			wantErr:                  nil,
			wantSkipConventionChecks: true,
			wantUserDefinedBadNames:  map[string]struct{}{"custom": {}},
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
			if rule.skipConventionChecks != tt.wantSkipConventionChecks {
				t.Errorf("unexpected skipConventionChecks: got = %v, want %v", rule.skipConventionChecks, tt.wantSkipConventionChecks)
			}
			if rule.skipTopLevelChecks != tt.wantSkipTopLevelChecks {
				t.Errorf("unexpected skipTopLevelChecks: got = %v, want %v", rule.skipTopLevelChecks, tt.wantSkipTopLevelChecks)
			}
			if rule.skipDefaultBadNameChecks != tt.wantSkipDefaultBadNameChecks {
				t.Errorf("unexpected skipDefaultBadNameChecks: got = %v, want %v", rule.skipDefaultBadNameChecks, tt.wantSkipDefaultBadNameChecks)
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

	t.Run("initializes alreadyCheckedNames", func(t *testing.T) {
		var rule PackageNamingRule

		err := rule.Configure(nil)
		if err != nil {
			t.Fatalf("unexpected error: got = %v, want = nil", err)
		}

		if rule.alreadyCheckedNames.elements == nil {
			t.Error("expected alreadyCheckedNames to be initialized, but got nil")
		}
	})
}
