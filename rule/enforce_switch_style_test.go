package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestEnforceSwitchStyleRule_Configure(t *testing.T) {
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
			name: "valid argument: allowNoDefault",
			arguments: lint.Arguments{
				"allowNoDefault",
			},
		},
		{
			name: "valid argument: allowDefaultNotLast",
			arguments: lint.Arguments{
				"allowDefaultNotLast",
			},
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{
				"allow-default-not-last",
				"allow-no-default",
			},
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{
				"allowdefaultnotlast",
				"allownodefault",
			},
		},
		{
			name: "unknown argument: unknown",
			arguments: lint.Arguments{
				"unknown",
			},
			wantErr: errors.New(`invalid argument "unknown" for rule enforce-switch-style; expected "allowNoDefault" or "allowDefaultNotLast"`),
		},
		{
			name: "unexpected type argument: 10",
			arguments: lint.Arguments{
				10,
			},
			wantErr: errors.New(`invalid argument for rule enforce-switch-style; expected string but got int`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.EnforceSwitchStyleRule

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
