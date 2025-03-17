package rule

import (
	"errors"
	"regexp"
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestUnusedReceiverRule_Configure(t *testing.T) {
	tests := []struct {
		name           string
		arguments      lint.Arguments
		wantErr        error
		wantRegex      *regexp.Regexp
		wantFailureMsg string
	}{
		{
			name:           "no arguments",
			arguments:      lint.Arguments{},
			wantErr:        nil,
			wantRegex:      regexp.MustCompile("^_$"),
			wantFailureMsg: "method receiver '%s' is not referenced in method's body, consider removing or renaming it as _",
		},
		{
			name: "valid arguments",
			arguments: lint.Arguments{
				map[string]any{
					"allowRegex": "^_",
				},
			},
			wantErr:        nil,
			wantRegex:      regexp.MustCompile("^_"),
			wantFailureMsg: "method receiver '%s' is not referenced in method's body, consider removing or renaming it to match ^_",
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{
				map[string]any{
					"allowregex": "^_",
				},
			},
			wantErr:        nil,
			wantRegex:      regexp.MustCompile("^_"),
			wantFailureMsg: "method receiver '%s' is not referenced in method's body, consider removing or renaming it to match ^_",
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{
				map[string]any{
					"allow-regex": "^_",
				},
			},
			wantErr:        nil,
			wantRegex:      regexp.MustCompile("^_"),
			wantFailureMsg: "method receiver '%s' is not referenced in method's body, consider removing or renaming it to match ^_",
		},
		{
			name: "argument is not a map",
			arguments: lint.Arguments{
				"invalid_argument",
			},
			wantErr:        nil,
			wantFailureMsg: "method receiver '%s' is not referenced in method's body, consider removing or renaming it as _",
		},
		{
			name: "missing allowRegex key",
			arguments: lint.Arguments{
				map[string]any{},
			},
			wantErr:        nil,
			wantRegex:      allowBlankIdentifierRegex,
			wantFailureMsg: "method receiver '%s' is not referenced in method's body, consider removing or renaming it as _",
		},
		{
			name: "invalid allowRegex type",
			arguments: lint.Arguments{
				map[string]any{
					"allowRegex": 123,
				},
			},
			wantErr: errors.New("error configuring [unused-receiver] rule: allowRegex is not string but [int]"),
		},
		{
			name: "invalid allowRegex value",
			arguments: lint.Arguments{
				map[string]any{
					"allowRegex": "[invalid",
				},
			},
			wantErr: errors.New("error configuring [unused-receiver] rule: allowRegex is not valid regex [[invalid]: error parsing regexp: missing closing ]: `[invalid`"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rule UnusedReceiverRule

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
			if tt.wantRegex != nil && rule.allowRegex.String() != tt.wantRegex.String() {
				t.Errorf("unexpected allowRegex: got = %v, want = %v", rule.allowRegex.String(), tt.wantRegex.String())
			}
			if tt.wantFailureMsg != rule.failureMsg {
				t.Errorf("unexpected failureMsg: got = %v, want = %v", rule.failureMsg, tt.wantFailureMsg)
			}
		})
	}
}
