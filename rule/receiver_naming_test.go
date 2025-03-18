package rule

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestReceiverNamingRule_Configure(t *testing.T) {
	tests := []struct {
		name          string
		arguments     lint.Arguments
		wantErr       error
		wantMaxLength int
	}{
		{
			name:          "no arguments",
			arguments:     lint.Arguments{},
			wantErr:       nil,
			wantMaxLength: -1,
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
			wantErr:       nil,
			wantMaxLength: 10,
		},
		{
			name: "valid lowercased argument",
			arguments: lint.Arguments{map[string]any{
				"maxlength": int64(10),
			}},
			wantErr:       nil,
			wantMaxLength: 10,
		},
		{
			name: "valid kebab-cased argument",
			arguments: lint.Arguments{map[string]any{
				"max-length": int64(10),
			}},
			wantErr:       nil,
			wantMaxLength: 10,
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
			var rule ReceiverNamingRule

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
			if rule.receiverNameMaxLength != tt.wantMaxLength {
				t.Errorf("unexpected maxLength: got = %v, want %v", rule.receiverNameMaxLength, tt.wantMaxLength)
			}
		})
	}
}
