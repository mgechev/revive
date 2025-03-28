package rule

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestEnforceRepeatedArgTypeStyleRule_Configure(t *testing.T) {
	tests := []struct {
		name                string
		arguments           lint.Arguments
		wantErr             error
		wantFuncArgStyle    enforceRepeatedArgTypeStyleType
		wantFuncRetValStyle enforceRepeatedArgTypeStyleType
	}{
		{
			name:                "no arguments",
			arguments:           lint.Arguments{},
			wantErr:             nil,
			wantFuncArgStyle:    "any",
			wantFuncRetValStyle: "any",
		},
		{
			name: "valid arguments: short",
			arguments: lint.Arguments{
				"short",
			},
			wantErr:             nil,
			wantFuncArgStyle:    "short",
			wantFuncRetValStyle: "short",
		},
		{
			name: "valid arguments",
			arguments: lint.Arguments{
				map[string]any{
					"funcArgStyle":    "full",
					"funcRetValStyle": "short",
				},
			},
			wantErr:             nil,
			wantFuncArgStyle:    "full",
			wantFuncRetValStyle: "short",
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{
				map[string]any{
					"funcargstyle":    "full",
					"funcretvalstyle": "short",
				},
			},
			wantErr:             nil,
			wantFuncArgStyle:    "full",
			wantFuncRetValStyle: "short",
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{
				map[string]any{
					"func-arg-style":     "full",
					"func-ret-val-style": "short",
				},
			},
			wantErr:             nil,
			wantFuncArgStyle:    "full",
			wantFuncRetValStyle: "short",
		},
		{
			name: "unrecognized key",
			arguments: lint.Arguments{
				map[string]any{
					"unknownKey": "someValue",
				},
			},
			wantErr: errors.New("invalid map key for 'enforce-repeated-arg-type-style' rule. Expecting 'funcArgStyle' or 'funcRetValStyle', got unknownKey"),
		},
		{
			name: "invalid argument type",
			arguments: lint.Arguments{
				123,
			},
			wantErr: errors.New("invalid argument '123' for 'import-alias-naming' rule. Expecting string or map[string]string, got int"),
		},
		{
			name: "invalid argument when string",
			arguments: lint.Arguments{
				"invalid_argument",
			},
			wantErr: errors.New("invalid argument to the enforce-repeated-arg-type-style rule: invalid repeated arg type style: invalid_argument (expecting one of [any short full])"),
		},
		{
			name: "invalid funcArgStyle value",
			arguments: lint.Arguments{
				map[string]any{
					"funcArgStyle": 123,
				},
			},
			wantErr: errors.New("invalid map value type for 'enforce-repeated-arg-type-style' rule. Expecting string, got int"),
		},
		{
			name: "invalid funcRetValStyle value",
			arguments: lint.Arguments{
				map[string]any{
					"funcRetValStyle": 123,
				},
			},
			wantErr: errors.New("invalid map value '123' for 'enforce-repeated-arg-type-style' rule. Expecting string, got int"),
		},
		{
			name: "invalid funcArgStyle value: wrong string",
			arguments: lint.Arguments{
				map[string]any{
					"funcArgStyle": "invalid",
				},
			},
			wantErr: errors.New("invalid argument to the enforce-repeated-arg-type-style rule: invalid repeated arg type style: invalid (expecting one of [any short full])"),
		},
		{
			name: "invalid funcRetValStyle value: wrong string",
			arguments: lint.Arguments{
				map[string]any{
					"funcRetValStyle": "invalid",
				},
			},
			wantErr: errors.New("invalid argument to the enforce-repeated-arg-type-style rule: invalid repeated arg type style: invalid (expecting one of [any short full])"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rule EnforceRepeatedArgTypeStyleRule

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
			if rule.funcArgStyle != tt.wantFuncArgStyle {
				t.Errorf("unexpected funcArgStyle: got = %v, want %v", rule.funcArgStyle, tt.wantFuncArgStyle)
			}
			if rule.funcRetValStyle != tt.wantFuncRetValStyle {
				t.Errorf("unexpected funcRetValStyle: got = %v, want %v", rule.funcRetValStyle, tt.wantFuncRetValStyle)
			}
		})
	}
}
