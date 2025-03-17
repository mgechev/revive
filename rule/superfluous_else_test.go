package rule

import (
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestSuperfluousElseRule_Configure(t *testing.T) {
	tests := []struct {
		name              string
		arguments         lint.Arguments
		wantErr           error
		wantPreserveScope bool
	}{
		{
			name:              "no arguments",
			arguments:         lint.Arguments{},
			wantErr:           nil,
			wantPreserveScope: false,
		},
		{
			name: "valid arguments",
			arguments: lint.Arguments{
				"preserveScope",
			},
			wantErr:           nil,
			wantPreserveScope: true,
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{
				"preservescope",
			},
			wantErr:           nil,
			wantPreserveScope: true,
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{
				"preserve-scope",
			},
			wantErr:           nil,
			wantPreserveScope: true,
		},
		{
			name: "invalid arguments",
			arguments: lint.Arguments{
				"unknown",
			},
			wantErr:           nil,
			wantPreserveScope: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rule SuperfluousElseRule

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
			if rule.preserveScope != tt.wantPreserveScope {
				t.Errorf("unexpected preserveScope: got = %v, want = %v", rule.preserveScope, tt.wantPreserveScope)
			}
		})
	}
}
