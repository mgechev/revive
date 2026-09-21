package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestFileLengthLimitRule_Configure(t *testing.T) {
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
			name: "valid arguments",
			arguments: lint.Arguments{map[string]any{
				"max":            int64(100),
				"skipComments":   true,
				"skipBlankLines": true,
			}},
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{map[string]any{
				"max":            int64(100),
				"skipcomments":   true,
				"skipblanklines": true,
			}},
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{map[string]any{
				"max":              int64(100),
				"skip-comments":    true,
				"skip-blank-lines": true,
			}},
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
			var r rule.FileLengthLimitRule

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
