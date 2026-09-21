package rule

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestIdenticalSwitchBranchesRule_Configure(t *testing.T) {
	tests := []struct {
		name                      string
		arguments                 lint.Arguments
		wantErr                   error
		wantAllowIdenticalDefault bool
	}{
		{
			name:                      "no arguments",
			arguments:                 lint.Arguments{},
			wantErr:                   nil,
			wantAllowIdenticalDefault: false,
		},
		{
			name: "allow-identical-default enabled",
			arguments: lint.Arguments{map[string]any{
				"allow-identical-default": true,
			}},
			wantErr:                   nil,
			wantAllowIdenticalDefault: true,
		},
		{
			name: "allow-identical-default disabled",
			arguments: lint.Arguments{map[string]any{
				"allow-identical-default": false,
			}},
			wantErr:                   nil,
			wantAllowIdenticalDefault: false,
		},
		{
			name: "camelCased argument",
			arguments: lint.Arguments{map[string]any{
				"allowIdenticalDefault": true,
			}},
			wantErr:                   nil,
			wantAllowIdenticalDefault: true,
		},
		{
			name: "lowercased argument",
			arguments: lint.Arguments{map[string]any{
				"allowidenticaldefault": true,
			}},
			wantErr:                   nil,
			wantAllowIdenticalDefault: true,
		},
		{
			name:      "invalid argument type",
			arguments: lint.Arguments{123},
			wantErr:   errors.New("invalid argument to the identical-switch-branches rule. Expecting a k,v map, got int"),
		},
		{
			name: "invalid allow-identical-default type",
			arguments: lint.Arguments{map[string]any{
				"allow-identical-default": "invalid",
			}},
			wantErr: errors.New(`invalid configuration value for "allow-identical-default" in identical-switch-branches rule; need bool but got string`),
		},
		{
			name: "unknown option",
			arguments: lint.Arguments{map[string]any{
				"unknown": true,
			}},
			wantErr: errors.New(`invalid argument "unknown" for rule identical-switch-branches; expected "allow-identical-default"`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rule IdenticalSwitchBranchesRule

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
			if rule.allowIdenticalDefault != tt.wantAllowIdenticalDefault {
				t.Errorf("unexpected allowIdenticalDefault: got = %v, want %v", rule.allowIdenticalDefault, tt.wantAllowIdenticalDefault)
			}
		})
	}
}
