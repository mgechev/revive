package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestReceiverNamingRule_Configure(t *testing.T) {
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
			name:      "invalid type",
			arguments: lint.Arguments{123},
			wantErr:   errors.New("unable to get arguments for rule receiver-naming. Expected object of key-value-pairs"),
		},
		{
			name: "valid maxLength argument",
			arguments: lint.Arguments{map[string]any{
				"maxLength": int64(10),
			}},
		},
		{
			name: "valid lowercased argument",
			arguments: lint.Arguments{map[string]any{
				"maxlength": int64(10),
			}},
		},
		{
			name: "valid kebab-cased argument",
			arguments: lint.Arguments{map[string]any{
				"max-length": int64(10),
			}},
		},
		{
			name: "invalid maxLength type",
			arguments: lint.Arguments{map[string]any{
				"maxLength": "10",
			}},
			wantErr: errors.New("invalid value 10 for argument maxLength of rule receiver-naming, expected integer value got string"),
		},
		{
			name: "unknown argument",
			arguments: lint.Arguments{map[string]any{
				"unknownKey": 10,
			}},
			wantErr: errors.New("unknown argument unknownKey for receiver-naming rule"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.ReceiverNamingRule

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
