package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestUnusedReceiverRule_Configure(t *testing.T) {
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
			name: "argument is not a map",
			arguments: lint.Arguments{
				"invalid_argument",
			},
		},
		{
			name: "missing allowRegex key",
			arguments: lint.Arguments{
				map[string]any{},
			},
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
			var r rule.UnusedReceiverRule

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
