package rule

import (
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestImportAliasNamingRule_Configure(t *testing.T) {
	tests := []struct {
		name           string
		arguments      lint.Arguments
		wantErr        error
		wantAllowRegex *regexp.Regexp
		wantDenyRegex  *regexp.Regexp
	}{
		{
			name:           "no arguments",
			arguments:      lint.Arguments{},
			wantErr:        nil,
			wantAllowRegex: regexp.MustCompile("^[a-z][a-z0-9]{0,}$"), //nolint:gocritic // regexpSimplify: backward compatibility
			wantDenyRegex:  nil,
		},
		{
			name:           "valid string argument",
			arguments:      lint.Arguments{"^[a-z][a-z0-9]*$"},
			wantErr:        nil,
			wantAllowRegex: regexp.MustCompile("^[a-z][a-z0-9]*$"),
			wantDenyRegex:  nil,
		},
		{
			name: "valid map arguments",
			arguments: lint.Arguments{map[string]any{
				"allowRegex": "^[a-z][a-z0-9]*$",
				"denyRegex":  "^v\\d+$",
			}},
			wantErr:        nil,
			wantAllowRegex: regexp.MustCompile("^[a-z][a-z0-9]*$"),
			wantDenyRegex:  regexp.MustCompile(`^v\d+$`),
		},
		{
			name: "valid map lowercased arguments",
			arguments: lint.Arguments{map[string]any{
				"allowregex": "^[a-z][a-z0-9]*$",
				"denyregex":  "^v\\d+$",
			}},
			wantErr:        nil,
			wantAllowRegex: regexp.MustCompile("^[a-z][a-z0-9]*$"),
			wantDenyRegex:  regexp.MustCompile(`^v\d+$`),
		},
		{
			name: "valid map kebab-cased arguments",
			arguments: lint.Arguments{map[string]any{
				"allow-regex": "^[a-z][a-z0-9]*$",
				"deny-regex":  "^v\\d+$",
			}},
			wantErr:        nil,
			wantAllowRegex: regexp.MustCompile("^[a-z][a-z0-9]*$"),
			wantDenyRegex:  regexp.MustCompile(`^v\d+$`),
		},
		{
			name:      "invalid argument type",
			arguments: lint.Arguments{123},
			wantErr:   errors.New(`invalid argument '123' for 'import-alias-naming' rule. Expecting string or map[string]string, got int`),
		},
		{
			name:      "invalid string argument regex",
			arguments: lint.Arguments{"["},
			wantErr:   errors.New("invalid argument to the import-alias-naming allowRegexp rule. Expecting \"[\" to be a valid regular expression, got: error parsing regexp: missing closing ]: `[`"),
		},
		{
			name: "invalid map key",
			arguments: lint.Arguments{
				map[string]any{
					"unknownKey": "value",
				},
			},
			wantErr: errors.New(`invalid map key for 'import-alias-naming' rule. Expecting 'allowRegex' or 'denyRegex', got unknownKey`),
		},
		{
			name: "invalid allowRegex type",
			arguments: lint.Arguments{map[string]any{
				"allowRegex": 123,
			}},
			wantErr: errors.New("invalid argument '123' for import-alias-naming allowRegexp rule. Expecting string, got int"),
		},
		{
			name: "invalid denyRegex type",
			arguments: lint.Arguments{map[string]any{
				"denyRegex": 123,
			}},
			wantErr: errors.New("invalid argument '123' for import-alias-naming denyRegexp rule. Expecting string, got int"),
		},
		{
			name: "invalid denyRegex regex",
			arguments: lint.Arguments{map[string]any{
				"denyRegex": "[",
			}},
			wantErr: errors.New(`invalid argument to the import-alias-naming denyRegexp rule. Expecting "[" to be a valid regular expression, got: error parsing regexp: missing closing ]: ` + "`[`"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rule ImportAliasNamingRule

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
			if !reflect.DeepEqual(rule.allowRegexp, tt.wantAllowRegex) {
				t.Errorf("unexpected allowRegex: got = %v, want %v", rule.allowRegexp, tt.wantAllowRegex)
			}
			if !reflect.DeepEqual(rule.denyRegexp, tt.wantDenyRegex) {
				t.Errorf("unexpected denyRegexp: got = %v, want %v", rule.denyRegexp, tt.wantDenyRegex)
			}
		})
	}
}
