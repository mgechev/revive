package rule_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestIndentErrorFlowRule_Configure(t *testing.T) {
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
			arguments: lint.Arguments{
				"preserveScope",
			},
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{
				"preservescope",
			},
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{
				"preserve-scope",
			},
		},
		{
			name: "invalid arguments",
			arguments: lint.Arguments{
				"unknown",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.IndentErrorFlowRule

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
