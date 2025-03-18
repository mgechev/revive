package rule

import (
	"errors"
	"reflect"
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestContextAsArgumentRule_Configure(t *testing.T) {
	tests := []struct {
		name      string
		arguments lint.Arguments
		wantErr   error
		wantTypes map[string]struct{}
	}{
		{
			name: "valid arguments",
			arguments: lint.Arguments{
				map[string]any{
					"allowTypesBefore": "AllowedBeforeType,AllowedBeforeStruct,*AllowedBeforePtrStruct,*testing.T",
				},
			},
			wantErr: nil,
			wantTypes: map[string]struct{}{
				"context.Context":         {},
				"AllowedBeforeType":       {},
				"AllowedBeforeStruct":     {},
				"*AllowedBeforePtrStruct": {},
				"*testing.T":              {},
			},
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{
				map[string]any{
					"allowtypesbefore": "AllowedBeforeType,AllowedBeforeStruct,*AllowedBeforePtrStruct,*testing.T",
				},
			},
			wantErr: nil,
			wantTypes: map[string]struct{}{
				"context.Context":         {},
				"AllowedBeforeType":       {},
				"AllowedBeforeStruct":     {},
				"*AllowedBeforePtrStruct": {},
				"*testing.T":              {},
			},
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{
				map[string]any{
					"allow-types-before": "AllowedBeforeType,AllowedBeforeStruct,*AllowedBeforePtrStruct,*testing.T",
				},
			},
			wantErr: nil,
			wantTypes: map[string]struct{}{
				"context.Context":         {},
				"AllowedBeforeType":       {},
				"AllowedBeforeStruct":     {},
				"*AllowedBeforePtrStruct": {},
				"*testing.T":              {},
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
			var rule ContextAsArgumentRule

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
			if !reflect.DeepEqual(rule.allowTypes, tt.wantTypes) {
				t.Errorf("unexpected allowTypes: got = %v, want %v", rule.allowTypes, tt.wantTypes)
			}
		})
	}
}
