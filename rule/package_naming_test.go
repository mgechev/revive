package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestPackageNamingRule_Configure(t *testing.T) {
	tests := []struct {
		name      string
		arguments lint.Arguments
		wantErr   error
	}{
		{
			name:      "no arguments",
			arguments: lint.Arguments{},
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
					"checkCollisionWithAllStd":   false,
				},
			},
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
					"checkcollisionwithallstd":   false,
				},
			},
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
					"check-collision-with-all-std":   false,
				},
			},
		},
		{
			name: "partial arguments",
			arguments: lint.Arguments{
				map[string]any{
					"skip-convention-name-check": true,
					"user-defined-bad-names":     []any{"custom"},
				},
			},
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
		},
		{
			name: "valid conventionNameCheckRegex",
			arguments: lint.Arguments{
				map[string]any{
					"convention-name-check-regex": "^[a-z][a-z0-9_]*$",
				},
			},
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
			wantErr: errors.New("invalid configuration for package-naming rule: skipConventionNameCheck and conventionNameCheckRegex cannot be both set"),
		},
		{
			name: "skipCollisionWithCommonStd with checkCollisionWithAllStd",
			arguments: lint.Arguments{
				map[string]any{
					"skip-collision-with-common-std": true,
					"check-collision-with-all-std":   true,
				},
			},
			wantErr: errors.New("invalid configuration for package-naming rule: skipCollisionWithCommonStd and checkCollisionWithAllStd cannot be both set"),
		},
		{
			name: "invalid skipTopLevelCheck type",
			arguments: lint.Arguments{
				map[string]any{
					"skip-top-level-check": "should-be-bool",
				},
			},
			wantErr: errors.New("invalid argument to the package-naming rule: expecting skipTopLevelCheck to be a boolean, but got string"),
		},
		{
			name: "invalid skipDefaultBadNameCheck type",
			arguments: lint.Arguments{
				map[string]any{
					"skip-default-bad-name-check": 42,
				},
			},
			wantErr: errors.New("invalid argument to the package-naming rule: expecting skipDefaultBadNameCheck to be a boolean, but got int"),
		},
		{
			name: "invalid skipCollisionWithCommonStd type",
			arguments: lint.Arguments{
				map[string]any{
					"skip-collision-with-common-std": []string{"invalid"},
				},
			},
			wantErr: errors.New("invalid argument to the package-naming rule: expecting skipCollisionWithCommonStd to be a boolean, but got []string"),
		},
		{
			name: "invalid checkCollisionWithAllStd type",
			arguments: lint.Arguments{
				map[string]any{
					"check-collision-with-all-std": 123.45,
				},
			},
			wantErr: errors.New("invalid argument to the package-naming rule: expecting checkCollisionWithAllStd to be a boolean, but got float64"),
		},
		{
			name: "valid checkExtraBadName set to true",
			arguments: lint.Arguments{
				map[string]any{
					"check-extra-bad-name": true,
				},
			},
		},
		{
			name: "valid checkExtraBadName set to false",
			arguments: lint.Arguments{
				map[string]any{
					"check-extra-bad-name": false,
				},
			},
		},
		{
			name: "invalid checkExtraBadName type",
			arguments: lint.Arguments{
				map[string]any{
					"check-extra-bad-name": "not-a-bool",
				},
			},
			wantErr: errors.New("invalid argument to the package-naming rule: expecting checkExtraBadName to be a boolean, but got string"),
		},
		{
			name:      "empty arguments slice",
			arguments: lint.Arguments{},
		},
		{
			name: "multiple arguments (more than 1)",
			arguments: lint.Arguments{
				map[string]any{"skip-convention-name-check": true},
				map[string]any{"skip-top-level-check": true},
			},
			wantErr: errors.New("invalid arguments to the package-naming rule: expected at most 1 argument, but got 2"),
		},
		{
			name: "unknown option key is silently ignored",
			arguments: lint.Arguments{
				map[string]any{
					"unknown-option":  true,
					"another-unknown": "value",
				},
			},
		},
		{
			name: "mixed known and unknown options",
			arguments: lint.Arguments{
				map[string]any{
					"skip-convention-name-check": true,
					"unknown-option":             "value",
				},
			},
		},
		{
			name: "all boolean options together",
			arguments: lint.Arguments{
				map[string]any{
					"skip-convention-name-check":     true,
					"skip-top-level-check":           true,
					"skip-default-bad-name-check":    true,
					"check-extra-bad-name":           true,
					"skip-collision-with-common-std": true,
					"check-collision-with-all-std":   false,
				},
			},
		},
		{
			name: "conventional names with empty userDefinedBadNames",
			arguments: lint.Arguments{
				map[string]any{
					"user-defined-bad-names": []any{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.PackageNamingRule

			err := r.Configure(tt.arguments)

			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("Configure() unexpected non-nil error %q", err)
				}
				return
			}
			if err == nil || err.Error() != tt.wantErr.Error() {
				t.Errorf("Configure() unexpected error: got %q, want %q", err, tt.wantErr)
			}
		})
	}
}
