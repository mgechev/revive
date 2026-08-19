package rule

import (
	"errors"
	"slices"
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestLineLengthLimitRule_Configure(t *testing.T) {
	tests := []struct {
		name         string
		arguments    lint.Arguments
		wantErr      error
		wantMax      int
		wantExcludes []string
	}{
		{
			name:      "no arguments",
			arguments: lint.Arguments{},
			wantErr:   nil,
			wantMax:   80,
		},
		{
			name:      "integer argument",
			arguments: lint.Arguments{int64(100)},
			wantErr:   nil,
			wantMax:   100,
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
			wantErr:      nil,
			wantMax:      100,
			wantExcludes: []string{`^\s*//go:generate `, `https?://`},
		},
		{
			name: "valid capitalized max option",
			arguments: lint.Arguments{map[string]any{
				"Max": int64(100),
			}},
			wantErr: nil,
			wantMax: 100,
		},
		{
			name: "map without max keeps default",
			arguments: lint.Arguments{map[string]any{
				"excludes": []any{`https?://`},
			}},
			wantErr:      nil,
			wantMax:      80,
			wantExcludes: []string{`https?://`},
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
			var rule LineLengthLimitRule

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
			if rule.max != tt.wantMax {
				t.Errorf("unexpected max: got = %v, want %v", rule.max, tt.wantMax)
			}

			gotExcludes := make([]string, 0, len(rule.excludes))
			for _, exclude := range rule.excludes {
				gotExcludes = append(gotExcludes, exclude.String())
			}
			if !slices.Equal(gotExcludes, tt.wantExcludes) {
				t.Errorf("unexpected excludes: got = %v, want %v", gotExcludes, tt.wantExcludes)
			}
		})
	}
}
