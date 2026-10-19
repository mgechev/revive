package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestLineLengthLimitRule_Configure(t *testing.T) {
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
			name:      "integer argument",
			arguments: lint.Arguments{int64(100)},
		},
		{
			name:      "negative integer argument",
			arguments: lint.Arguments{int64(-1)},
			wantErr:   errors.New(`invalid value passed as argument number to the "line-length-limit" rule`),
		},
		{
			name: "valid map arguments",
			arguments: lint.Arguments{map[string]any{
				"max":      int64(100),
				"excludes": []any{`^\s*//go:generate `, `https?://`},
			}},
		},
		{
			name: "valid capitalized max option",
			arguments: lint.Arguments{map[string]any{
				"Max": int64(100),
			}},
		},
		{
			name: "map without max keeps default",
			arguments: lint.Arguments{map[string]any{
				"excludes": []any{`https?://`},
			}},
		},
		{
			name:      "invalid argument type",
			arguments: lint.Arguments{"invalid"},
			wantErr:   errors.New(`invalid argument to the "line-length-limit" rule: expecting an integer or an options map, got string`),
		},
		{
			name: "invalid max type",
			arguments: lint.Arguments{map[string]any{
				"max": "invalid",
			}},
			wantErr: errors.New(`invalid value for the "max" option of the "line-length-limit" rule: expecting an integer, got string`),
		},
		{
			name: "negative max option",
			arguments: lint.Arguments{map[string]any{
				"max": int64(-1),
			}},
			wantErr: errors.New(`invalid value passed as argument number to the "line-length-limit" rule`),
		},
		{
			name: "invalid excludes type",
			arguments: lint.Arguments{map[string]any{
				"excludes": "invalid",
			}},
			wantErr: errors.New(`invalid value for the "excludes" option of the "line-length-limit" rule: expecting a slice of strings, got string`),
		},
		{
			name: "invalid excludes element type",
			arguments: lint.Arguments{map[string]any{
				"excludes": []any{123},
			}},
			wantErr: errors.New(`invalid value in the "excludes" option of the "line-length-limit" rule: expecting a string, got int`),
		},
		{
			name: "empty excludes pattern",
			arguments: lint.Arguments{map[string]any{
				"excludes": []any{""},
			}},
			wantErr: errors.New(`invalid value in the "excludes" option of the "line-length-limit" rule: regular expression must not be empty`),
		},
		{
			name: "excludes pattern that does not compile",
			arguments: lint.Arguments{map[string]any{
				"excludes": []any{"("},
			}},
			wantErr: errors.New("invalid value in the \"excludes\" option of the \"line-length-limit\" rule: regexp \"(\" does not compile: error parsing regexp: missing closing ): `(`"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.LineLengthLimitRule

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
