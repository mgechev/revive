package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestUnusedParamRule_Configure(t *testing.T) {
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
			name: "valid arguments",
			arguments: lint.Arguments{
				map[string]any{
					"allowRegex": "^_",
				},
			},
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{
				map[string]any{
					"allowregex": "^_",
				},
			},
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{
				map[string]any{
					"allow-regex": "^_",
				},
			},
		},
		{
			name: "missed allowRegex value",
			arguments: lint.Arguments{
				map[string]any{
					"unknownKey": "123",
				},
			},
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
			var r rule.UnusedParamRule

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
