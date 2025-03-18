package rule

import (
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestUnusedParamRule_Configure(t *testing.T) {
	tests := []struct {
		name           string
		arguments      lint.Arguments
		wantErr        error
		wantAllowRegex *regexp.Regexp
		wantFailureMsg string
	}{
		{
			name:           "no arguments",
			arguments:      lint.Arguments{},
			wantErr:        nil,
			wantAllowRegex: regexp.MustCompile("^_$"),
			wantFailureMsg: "parameter '%s' seems to be unused, consider removing or renaming it as _",
		},
		{
			name: "valid arguments",
			arguments: lint.Arguments{
				map[string]any{
					"allowRegex": "^_",
				},
			},
			wantErr:        nil,
			wantAllowRegex: regexp.MustCompile("^_"),
			wantFailureMsg: "parameter '%s' seems to be unused, consider removing or renaming it to match ^_",
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{
				map[string]any{
					"allowregex": "^_",
				},
			},
			wantErr:        nil,
			wantAllowRegex: regexp.MustCompile("^_"),
			wantFailureMsg: "parameter '%s' seems to be unused, consider removing or renaming it to match ^_",
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{
				map[string]any{
					"allow-regex": "^_",
				},
			},
			wantErr:        nil,
			wantAllowRegex: regexp.MustCompile("^_"),
			wantFailureMsg: "parameter '%s' seems to be unused, consider removing or renaming it to match ^_",
		},
		{
			name: "missed allowRegex value",
			arguments: lint.Arguments{
				map[string]any{
					"unknownKey": "123",
				},
			},
			wantErr:        nil,
			wantAllowRegex: regexp.MustCompile("^_$"),
			wantFailureMsg: "parameter '%s' seems to be unused, consider removing or renaming it as _",
		},
		{
			name: "invalid allowRegex: not a string",
			arguments: lint.Arguments{
				map[string]any{
					"allowRegex": 123,
				},
			},
			wantErr: errors.New("error configuring unused-parameter rule: allowRegex is not string but [int]"),
		},
		{
			name: "invalid allowRegex: not a valid regex",
			arguments: lint.Arguments{
				map[string]any{
					"allowRegex": "[",
				},
			},
			wantErr: errors.New("error configuring unused-parameter rule: allowRegex is not valid regex [[]: error parsing regexp: missing closing ]: " + "`[`"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rule UnusedParamRule

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
			if !reflect.DeepEqual(rule.allowRegex, tt.wantAllowRegex) {
				t.Errorf("unexpected allowRegex: got = %v, want %v", rule.allowRegex, tt.wantAllowRegex)
			}
			if rule.failureMsg != tt.wantFailureMsg {
				t.Errorf("unexpected failureMessage: got = %v, want %v", rule.failureMsg, tt.wantFailureMsg)
			}
		})
	}
}
