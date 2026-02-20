package rule

import (
	"errors"
	"reflect"
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestPackageNamingRule_Configure(t *testing.T) {
	tests := []struct {
		name                                  string
		arguments                             lint.Arguments
		wantErr                               error
		wantSkipPackageNameChecks             bool
		wantBadPackageNames                   map[string]struct{}
		wantSkipPackageNameCollisionWithGoStd bool
	}{
		{
			name:                      "no arguments",
			arguments:                 lint.Arguments{},
			wantErr:                   nil,
			wantSkipPackageNameChecks: false,
		},
		{
			name: "valid arguments",
			arguments: lint.Arguments{
				map[string]any{
					"skipPackageNameChecks":             true,
					"extraBadPackageNames":              []any{"helpers", "models"},
					"skipPackageNameCollisionWithGoStd": true,
				},
			},
			wantErr:                               nil,
			wantSkipPackageNameChecks:             true,
			wantBadPackageNames:                   map[string]struct{}{"helpers": {}, "models": {}},
			wantSkipPackageNameCollisionWithGoStd: true,
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{
				map[string]any{
					"skippackagenamechecks":             true,
					"extrabadpackagenames":              []any{"helpers", "models"},
					"skippackagenamecollisionwithgostd": true,
				},
			},
			wantErr:                               nil,
			wantSkipPackageNameChecks:             true,
			wantBadPackageNames:                   map[string]struct{}{"helpers": {}, "models": {}},
			wantSkipPackageNameCollisionWithGoStd: true,
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{
				map[string]any{
					"skip-initialism-name-checks":             true,
					"upper-case-const":                        true,
					"skip-package-name-checks":                true,
					"extra-bad-package-names":                 []any{"helpers", "models"},
					"skip-package-name-collision-with-go-std": true,
				},
			},
			wantErr:                               nil,
			wantSkipPackageNameChecks:             true,
			wantBadPackageNames:                   map[string]struct{}{"helpers": {}, "models": {}},
			wantSkipPackageNameCollisionWithGoStd: true,
		},
		{
			name: "invalid argument type",
			arguments: lint.Arguments{
				"invalid-arg",
			},
			wantErr: errors.New("invalid argument to the package-naming rule: expecting a k,v map, but got string"),
		},
		{
			name: "invalid extraBadPackageNames type",
			arguments: lint.Arguments{
				map[string]any{
					"extraBadPackageNames": "invalid-type",
				},
			},
			wantErr: errors.New("invalid argument to the package-naming rule: expecting extraBadPackageNames of type slice of strings, but got string"),
		},
		{
			name: "invalid extraBadPackageNames element type",
			arguments: lint.Arguments{
				map[string]any{
					"extraBadPackageNames": []any{"helpers", 123},
				},
			},
			wantErr: errors.New("invalid argument to the package-naming rule: expecting element 1 of extraBadPackageNames to be a string, but got 123(int)"),
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
			if !reflect.DeepEqual(rule.extraBadPackageNames, tt.wantBadPackageNames) {
				t.Errorf("unexpected extraBadPackageNames: got = %v, want %v", rule.extraBadPackageNames, tt.wantBadPackageNames)
			}
			if rule.skipPackageNameCollisionWithGoStd != tt.wantSkipPackageNameCollisionWithGoStd {
				t.Errorf("unexpected skipPackageNameCollisionWithGoStd: got = %v, want %v", rule.skipPackageNameCollisionWithGoStd, tt.wantSkipPackageNameCollisionWithGoStd)
			}
		})
	}
}

func TestPackageNamingRule_Configure_LoadStdPackages(t *testing.T) {
	t.Run("loads std packages", func(t *testing.T) {
		var rule PackageNamingRule

		err := rule.Configure(nil)
		if err != nil {
			t.Fatalf("unexpected error: got = %v, want = nil", err)
		}

		for _, pkg := range []string{"fmt", "http", "json", "rand", "io"} {
			if _, ok := rule.stdPackageNames[pkg]; !ok {
				t.Errorf("expected package %q to be loaded, but got empty", pkg)
			}
		}

		for _, pkg := range []string{"v2", "internal", "vendor", "std", "text", "go"} {
			if _, ok := rule.stdPackageNames[pkg]; ok {
				t.Errorf("expected package %q to be not loaded, but got loaded", pkg)
			}
		}
	})

	t.Run("does not reload std packages on subsequent Configure calls", func(t *testing.T) {
		var rule PackageNamingRule

		err := rule.Configure(nil)
		if err != nil {
			t.Fatalf("unexpected error: got = %v, want = nil", err)
		}
		firstLoadLen := len(rule.stdPackageNames)
		err = rule.Configure(nil)
		if err != nil {
			t.Fatalf("unexpected error: got = %v, want = nil", err)
		}
		if len(rule.stdPackageNames) != firstLoadLen {
			t.Errorf("expected stdPackageNames to be loaded only once, but got different lengths: first %d, second %d", firstLoadLen, len(rule.stdPackageNames))
		}
	})

	t.Run("skips loading std packages when collision check disabled", func(t *testing.T) {
		var rule PackageNamingRule

		err := rule.Configure(lint.Arguments{
			map[string]any{
				"skip-package-name-collision-with-go-std": true,
			},
		})
		if err != nil {
			t.Errorf("unexpected error: got = %v, want = nil", err)
		}

		if len(rule.stdPackageNames) != 0 {
			t.Errorf("expected stdPackageNames to be empty, but got %v", rule.stdPackageNames)
		}
	})
}
