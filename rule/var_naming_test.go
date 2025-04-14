package rule

import (
	"errors"
	"reflect"
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestVarNamingRule_Configure(t *testing.T) {
	tests := []struct {
		name                      string
		arguments                 lint.Arguments
		wantErr                   error
		wantAllowList             []string
		wantBlockList             []string
		wantAllowUpperCaseConst   bool
		wantSkipPackageNameChecks bool
		wantBadPackageNames       map[string]struct{}
	}{
		{
			name:                      "no arguments",
			arguments:                 lint.Arguments{},
			wantErr:                   nil,
			wantAllowList:             nil,
			wantBlockList:             nil,
			wantAllowUpperCaseConst:   false,
			wantSkipPackageNameChecks: false,
		},
		{
			name: "valid arguments",
			arguments: lint.Arguments{
				[]any{"ID"},
				[]any{"VM"},
				[]any{map[string]any{
					"upperCaseConst":        true,
					"skipPackageNameChecks": true,
					"badPackageNames":       []string{"helpers", "models"},
				}},
			},
			wantErr:                   nil,
			wantAllowList:             []string{"ID"},
			wantBlockList:             []string{"VM"},
			wantAllowUpperCaseConst:   true,
			wantSkipPackageNameChecks: true,
			wantBadPackageNames:       map[string]struct{}{"helpers": {}, "models": {}},
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{
				[]any{"ID"},
				[]any{"VM"},
				[]any{map[string]any{
					"uppercaseconst":        true,
					"skippackagenamechecks": true,
					"badpackagenames":       []string{"helpers", "models"},
				}},
			},
			wantErr:                   nil,
			wantAllowList:             []string{"ID"},
			wantBlockList:             []string{"VM"},
			wantAllowUpperCaseConst:   true,
			wantSkipPackageNameChecks: true,
			wantBadPackageNames:       map[string]struct{}{"helpers": {}, "models": {}},
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{
				[]any{"ID"},
				[]any{"VM"},
				[]any{map[string]any{
					"upper-case-const":         true,
					"skip-package-name-checks": true,
					"bad-package-names":        []string{"helpers", "models"},
				}},
			},
			wantErr:                   nil,
			wantAllowList:             []string{"ID"},
			wantBlockList:             []string{"VM"},
			wantAllowUpperCaseConst:   true,
			wantSkipPackageNameChecks: true,
			wantBadPackageNames:       map[string]struct{}{"helpers": {}, "models": {}},
		},
		{
			name:      "invalid allowlist type",
			arguments: lint.Arguments{123},
			wantErr:   errors.New("invalid argument to the var-naming rule. Expecting a allowlist of type slice with initialisms, got int"),
		},
		{
			name:      "invalid allowlist value type",
			arguments: lint.Arguments{[]any{123}},
			wantErr:   errors.New("invalid 123 values of the var-naming rule. Expecting slice of strings but got element of type []interface {}"),
		},
		{
			name:      "invalid blocklist type",
			arguments: lint.Arguments{[]any{"ID"}, 123},
			wantErr:   errors.New("invalid argument to the var-naming rule. Expecting a blocklist of type slice with initialisms, got int"),
		},
		{
			name:      "invalid third argument type",
			arguments: lint.Arguments{[]any{"ID"}, []any{"VM"}, 123},
			wantErr:   errors.New("invalid third argument to the var-naming rule. Expecting a options of type slice, got int"),
		},
		{
			name:      "invalid third argument slice size",
			arguments: lint.Arguments{[]any{"ID"}, []any{"VM"}, []any{}},
			wantErr:   errors.New("invalid third argument to the var-naming rule. Expecting a options of type slice, of len==1, but 0"),
		},
		{
			name:      "invalid third argument first element type",
			arguments: lint.Arguments{[]any{"ID"}, []any{"VM"}, []any{123}},
			wantErr:   errors.New("invalid third argument to the var-naming rule. Expecting a options of type slice, of len==1, with map, but int"),
		},
		{
			name:      "invalid third argument badPackageNames",
			arguments: lint.Arguments{[]any{""}, []any{""}, []any{map[string]any{"badPackageNames": []int{1}}}},
			wantErr:   errors.New("invalid third argument to the var-naming rule. Expecting badPackageNames of type slice of strings, but []int"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rule VarNamingRule

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
			if rule.allowUpperCaseConst != tt.wantAllowUpperCaseConst {
				t.Errorf("unexpected allowUpperCaseConst: got = %v, want %v", rule.allowUpperCaseConst, tt.wantAllowUpperCaseConst)
			}
			if rule.skipPackageNameChecks != tt.wantSkipPackageNameChecks {
				t.Errorf("unexpected skipPackageNameChecks: got = %v, want %v", rule.skipPackageNameChecks, tt.wantSkipPackageNameChecks)
			}
			if !reflect.DeepEqual(rule.badPackageNames, tt.wantBadPackageNames) {
				t.Errorf("unexpected badPackageNames: got = %v, want %v", rule.badPackageNames, tt.wantBadPackageNames)
			}
		})
	}
}
