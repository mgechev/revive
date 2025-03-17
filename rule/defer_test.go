package rule

import (
	"errors"
	"reflect"
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestDeferRule_Configure(t *testing.T) {
	tests := []struct {
		name      string
		arguments lint.Arguments
		wantErr   error
		wantAllow map[string]bool
	}{
		{
			name:      "no arguments",
			arguments: lint.Arguments{},
			wantErr:   nil,
			wantAllow: map[string]bool{
				"loop":             true,
				"callchain":        true,
				"methodcall":       true,
				"return":           true,
				"recover":          true,
				"immediaterecover": true,
			},
		},
		{
			name: "valid arguments",
			arguments: lint.Arguments{
				[]any{
					"loop",
					"callChain",
					"methodCall",
					"return",
					"recover",
					"immediateRecover",
				},
			},
			wantErr: nil,
			wantAllow: map[string]bool{
				"loop":             true,
				"callchain":        true,
				"methodcall":       true,
				"return":           true,
				"recover":          true,
				"immediaterecover": true,
			},
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{
				[]any{
					"loop",
					"callchain",
					"methodcall",
					"return",
					"recover",
					"immediaterecover",
				},
			},
			wantErr: nil,
			wantAllow: map[string]bool{
				"loop":             true,
				"callchain":        true,
				"methodcall":       true,
				"return":           true,
				"recover":          true,
				"immediaterecover": true,
			},
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{
				[]any{
					"loop",
					"call-chain",
					"method-call",
					"return",
					"recover",
					"immediate-recover",
				},
			},
			wantErr: nil,
			wantAllow: map[string]bool{
				"loop":             true,
				"callchain":        true,
				"methodcall":       true,
				"return":           true,
				"recover":          true,
				"immediaterecover": true,
			},
		},
		{
			name: "invalid argument type",
			arguments: lint.Arguments{
				"invalid_argument",
			},
			wantErr: errors.New("invalid argument 'invalid_argument' for 'defer' rule. Expecting []string, got string"),
		},
		{
			name: "invalid subcase type",
			arguments: lint.Arguments{
				[]any{123},
			},
			wantErr: errors.New("invalid argument '123' for 'defer' rule. Expecting string, got int"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rule DeferRule

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
			if !reflect.DeepEqual(rule.allow, tt.wantAllow) {
				t.Errorf("unexpected allow: got = %v, want %v", rule.allow, tt.wantAllow)
			}
		})
	}
}
