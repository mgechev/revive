package rule

import (
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestEarlyReturn_Configure(t *testing.T) {
	tests := []struct {
		name              string
		arguments         lint.Arguments
		wantErr           error
		wantPreserveScope bool
		wantAllowJump     bool
	}{
		{
			name:              "no arguments",
			arguments:         lint.Arguments{},
			wantErr:           nil,
			wantPreserveScope: false,
			wantAllowJump:     false,
		},
		{
			name: "valid arguments",
			arguments: lint.Arguments{
				"preserveScope",
				"allowJump",
			},
			wantErr:           nil,
			wantPreserveScope: true,
			wantAllowJump:     true,
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{
				"preservescope",
				"allowjump",
			},
			wantErr:           nil,
			wantPreserveScope: true,
			wantAllowJump:     true,
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{
				"preserve-scope",
				"allow-jump",
			},
			wantErr:           nil,
			wantPreserveScope: true,
			wantAllowJump:     true,
		},
		{
			name: "invalid arguments",
			arguments: lint.Arguments{
				"unknown",
			},
			wantErr:           nil,
			wantPreserveScope: false,
			wantAllowJump:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rule EarlyReturnRule

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
			if rule.allowJump != tt.wantAllowJump {
				t.Errorf("unexpected allowJump: got = %v, want = %v", rule.allowJump, tt.wantAllowJump)
			}
		})
	}
}
