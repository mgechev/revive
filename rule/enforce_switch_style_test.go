package rule

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestEnforceSwitchStyleRule_Configure(t *testing.T) {
	tests := []struct {
		name                    string
		arguments               lint.Arguments
		wantErr                 error
		wantAllowNoDefault      bool
		wantAllowDefaultNotLast bool
	}{
		{
			name:      "no arguments",
			arguments: lint.Arguments{},
			wantErr:   nil,
		},
		{
			name: "valid argument: allowNoDefault",
			arguments: lint.Arguments{
				"allowNoDefault",
			},
			wantErr:            nil,
			wantAllowNoDefault: true,
		},
		{
			name: "valid argument: allowDefaultNotLast",
			arguments: lint.Arguments{
				"allowDefaultNotLast",
			},
			wantErr:                 nil,
			wantAllowDefaultNotLast: true,
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{
				"allowdefaultnotlast",
				"allownodefault",
			},
			wantErr:                 nil,
			wantAllowNoDefault:      true,
			wantAllowDefaultNotLast: true,
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
			var rule EnforceSwitchStyleRule

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
			if tt.wantAllowNoDefault != rule.allowNoDefault {
				t.Errorf("unexpected allowNoDefault: want = %v, got %v", tt.wantAllowNoDefault, rule.allowNoDefault)
			}
			if tt.wantAllowDefaultNotLast != rule.allowDefaultNotLast {
				t.Errorf("unexpected funcRetValStyle: want = %v, got %v", tt.wantAllowDefaultNotLast, rule.allowDefaultNotLast)
			}
		})
	}
}
