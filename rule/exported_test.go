package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestExportedRule_Configure(t *testing.T) {
	tests := []struct {
		name      string
		arguments lint.Arguments
		wantErr   error
	}{
		{
			name:      "default configuration",
			arguments: lint.Arguments{},
		},
		{
			name: "valid arguments",
			arguments: lint.Arguments{
				"checkPrivateReceivers",
				"disableStutteringCheck",
				"checkPublicInterface",
				"disableChecksOnConstants",
				"disableChecksOnFunctions",
				"disableChecksOnMethods",
				"disableChecksOnTypes",
				"disableChecksOnVariables",
			},
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{
				"checkprivatereceivers",
				"disablestutteringcheck",
				"checkpublicinterface",
				"disablechecksonconstants",
				"disablechecksonfunctions",
				"disablechecksonmethods",
				"disablechecksontypes",
				"disablechecksonvariables",
			},
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{
				"check-private-receivers",
				"disable-stuttering-check",
				"check-public-interface",
				"disable-checks-on-constants",
				"disable-checks-on-functions",
				"disable-checks-on-methods",
				"disable-checks-on-types",
				"disable-checks-on-variables",
			},
		},
		{
			name: "valid sayRepetitiveInsteadOfStutters",
			arguments: lint.Arguments{
				"sayRepetitiveInsteadOfStutters",
			},
		},
		{
			name: "valid lowercased sayRepetitiveInsteadOfStutters",
			arguments: lint.Arguments{
				"sayrepetitiveinsteadofstutters",
			},
		},
		{
			name: "valid kebab-cased sayRepetitiveInsteadOfStutters",
			arguments: lint.Arguments{
				"say-repetitive-instead-of-stutters",
			},
		},
		{
			name:      "unknown configuration flag",
			arguments: lint.Arguments{"unknownFlag"},
			wantErr:   errors.New("unknown configuration flag unknownFlag for exported rule"),
		},
		{
			name:      "invalid argument type",
			arguments: lint.Arguments{123},
			wantErr:   errors.New("invalid argument for the exported rule: expecting a string, got int"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.ExportedRule

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
