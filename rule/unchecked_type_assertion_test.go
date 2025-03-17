package rule

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestUncheckedTypeAssertionRule_Configure(t *testing.T) {
	tests := []struct {
		name                             string
		arguments                        lint.Arguments
		wantErr                          error
		wantAcceptIgnoredAssertionResult bool
	}{
		{
			name:                             "no arguments",
			arguments:                        lint.Arguments{},
			wantErr:                          nil,
			wantAcceptIgnoredAssertionResult: false,
		},
		{
			name: "valid acceptIgnoredAssertionResult argument",
			arguments: lint.Arguments{map[string]any{
				"acceptIgnoredAssertionResult": true,
			}},
			wantErr:                          nil,
			wantAcceptIgnoredAssertionResult: true,
		},
		{
			name: "valid lowercased argument",
			arguments: lint.Arguments{map[string]any{
				"acceptignoredassertionresult": true,
			}},
			wantErr:                          nil,
			wantAcceptIgnoredAssertionResult: true,
		},
		{
			name: "valid kebab-cased argument",
			arguments: lint.Arguments{map[string]any{
				"accept-ignored-assertion-result": true,
			}},
			wantErr:                          nil,
			wantAcceptIgnoredAssertionResult: true,
		},
		{
			name:      "invalid type",
			arguments: lint.Arguments{123},
			wantErr:   errors.New("unable to get arguments. Expected object of key-value-pairs"),
		},
		{
			name: "invalid acceptIgnoredAssertionResult type",
			arguments: lint.Arguments{map[string]any{
				"acceptIgnoredAssertionResult": "true",
			}},
			wantErr: errors.New("unable to parse argument 'acceptIgnoredAssertionResult'. Expected boolean"),
		},
		{
			name: "unknown argument",
			arguments: lint.Arguments{map[string]any{
				"unknownKey": true,
			}},
			wantErr: errors.New("unknown argument: unknownKey"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rule UncheckedTypeAssertionRule

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
			if rule.acceptIgnoredAssertionResult != tt.wantAcceptIgnoredAssertionResult {
				t.Errorf("unexpected acceptIgnoredAssertionResult: got = %v, want %v", rule.acceptIgnoredAssertionResult, tt.wantAcceptIgnoredAssertionResult)
			}
		})
	}
}
