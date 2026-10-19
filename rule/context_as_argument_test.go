package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestContextAsArgumentRule_Configure(t *testing.T) {
	tests := []struct {
		name      string
		arguments lint.Arguments
		wantErr   error
	}{
		{
			name: "valid arguments",
			arguments: lint.Arguments{
				map[string]any{
					"allowTypesBefore": "AllowedBeforeType,AllowedBeforeStruct,*AllowedBeforePtrStruct,*testing.T",
				},
			},
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{
				map[string]any{
					"allowtypesbefore": "AllowedBeforeType,AllowedBeforeStruct,*AllowedBeforePtrStruct,*testing.T",
				},
			},
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{
				map[string]any{
					"allow-types-before": "AllowedBeforeType,AllowedBeforeStruct,*AllowedBeforePtrStruct,*testing.T",
				},
			},
		},
		{
			name: "invalid argument type",
			arguments: lint.Arguments{
				"invalid_argument",
			},
			wantErr: errors.New("invalid argument to the context-as-argument rule. Expecting a k,v map, got string"),
		},
		{
			name: "invalid allowTypesBefore value",
			arguments: lint.Arguments{
				map[string]any{
					"allowTypesBefore": 123,
				},
			},
			wantErr: errors.New("invalid argument to the context-as-argument.allowTypesBefore rule. Expecting a string, got int"),
		},
		{
			name: "unrecognized key",
			arguments: lint.Arguments{
				map[string]any{
					"unknownKey": "someValue",
				},
			},
			wantErr: errors.New("invalid argument to the context-as-argument rule. Unrecognized key unknownKey"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.ContextAsArgumentRule

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
