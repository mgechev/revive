package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestUncheckedTypeAssertionRule_Configure(t *testing.T) {
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
			name: "valid acceptIgnoredAssertionResult argument",
			arguments: lint.Arguments{map[string]any{
				"acceptIgnoredAssertionResult": true,
			}},
		},
		{
			name: "valid lowercased argument",
			arguments: lint.Arguments{map[string]any{
				"acceptignoredassertionresult": true,
			}},
		},
		{
			name: "valid kebab-cased argument",
			arguments: lint.Arguments{map[string]any{
				"accept-ignored-assertion-result": true,
			}},
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
			var r rule.UncheckedTypeAssertionRule

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
