package rule

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestFileLengthLimitRule_Configure(t *testing.T) {
	tests := []struct {
		name               string
		arguments          lint.Arguments
		wantErr            error
		wantMax            int
		wantSkipComments   bool
		wantSkipBlankLines bool
	}{
		{
			name:               "no arguments",
			arguments:          lint.Arguments{},
			wantErr:            nil,
			wantMax:            0,
			wantSkipComments:   false,
			wantSkipBlankLines: false,
		},
		{
			name: "valid arguments",
			arguments: lint.Arguments{map[string]any{
				"max":            int64(100),
				"skipComments":   true,
				"skipBlankLines": true,
			}},
			wantErr:            nil,
			wantMax:            100,
			wantSkipComments:   true,
			wantSkipBlankLines: true,
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{map[string]any{
				"max":            int64(100),
				"skipcomments":   true,
				"skipblanklines": true,
			}},
			wantErr:            nil,
			wantMax:            100,
			wantSkipComments:   true,
			wantSkipBlankLines: true,
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{map[string]any{
				"max":              int64(100),
				"skip-comments":    true,
				"skip-blank-lines": true,
			}},
			wantErr:            nil,
			wantMax:            100,
			wantSkipComments:   true,
			wantSkipBlankLines: true,
		},
		{
			name:      "invalid argument",
			arguments: lint.Arguments{123},
			wantErr:   errors.New(`invalid argument to the "file-length-limit" rule. Expecting a k,v map, got int`),
		},
		{
			name: "invalid max type",
			arguments: lint.Arguments{map[string]any{
				"max": "invalid",
			}},
			wantErr: errors.New(`invalid configuration value for max lines in "file-length-limit" rule; need positive int64 but got string`),
		},
		{
			name: "invalid skipComments type",
			arguments: lint.Arguments{map[string]any{
				"skipComments": "invalid",
			}},
			wantErr: errors.New(`invalid configuration value for skip comments in "file-length-limit" rule; need bool but got string`),
		},
		{
			name: "invalid skipBlankLines type",
			arguments: lint.Arguments{map[string]any{
				"skipBlankLines": "invalid",
			}},
			wantErr: errors.New(`invalid configuration value for skip blank lines in "file-length-limit" rule; need bool but got string`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rule FileLengthLimitRule

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
			if rule.skipComments != tt.wantSkipComments {
				t.Errorf("unexpected skipComments: got = %v, want %v", rule.skipComments, tt.wantSkipComments)
			}
			if rule.skipBlankLines != tt.wantSkipBlankLines {
				t.Errorf("unexpected skipBlankLines: got = %v, want %v", rule.skipBlankLines, tt.wantSkipBlankLines)
			}
		})
	}
}
