package rule

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestExportedRule_Configure(t *testing.T) {
	tests := []struct {
		name               string
		arguments          lint.Arguments
		wantErr            error
		wantDisabledChecks disabledChecks
		wantStuttersMsg    string
	}{
		{
			name:               "default configuration",
			arguments:          lint.Arguments{},
			wantErr:            nil,
			wantDisabledChecks: disabledChecks{PrivateReceivers: true, PublicInterfaces: true},
			wantStuttersMsg:    "stutters",
		},
		{
			name:               "checkPrivateReceivers",
			arguments:          lint.Arguments{"checkprivatereceivers"},
			wantErr:            nil,
			wantDisabledChecks: disabledChecks{PrivateReceivers: false, PublicInterfaces: true},
			wantStuttersMsg:    "stutters",
		},
		{
			name:               "disableStutteringCheck",
			arguments:          lint.Arguments{"disablestutteringcheck"},
			wantErr:            nil,
			wantDisabledChecks: disabledChecks{PrivateReceivers: true, PublicInterfaces: true, Stuttering: true},
			wantStuttersMsg:    "stutters",
		},
		{
			name:               "sayRepetitiveInsteadOfStutters",
			arguments:          lint.Arguments{"sayrepetitiveinsteadofstutters"},
			wantErr:            nil,
			wantDisabledChecks: disabledChecks{PrivateReceivers: true, PublicInterfaces: true},
			wantStuttersMsg:    "is repetitive",
		},
		{
			name:               "checkPublicInterface",
			arguments:          lint.Arguments{"checkpublicinterface"},
			wantErr:            nil,
			wantDisabledChecks: disabledChecks{PrivateReceivers: true, PublicInterfaces: false},
			wantStuttersMsg:    "stutters",
		},
		{
			name:               "disableChecksOnConstants",
			arguments:          lint.Arguments{"disablechecksonconstants"},
			wantErr:            nil,
			wantDisabledChecks: disabledChecks{PrivateReceivers: true, PublicInterfaces: true, Const: true},
			wantStuttersMsg:    "stutters",
		},
		{
			name:               "disableChecksOnFunctions",
			arguments:          lint.Arguments{"disablechecksonfunctions"},
			wantErr:            nil,
			wantDisabledChecks: disabledChecks{PrivateReceivers: true, PublicInterfaces: true, Function: true},
			wantStuttersMsg:    "stutters",
		},
		{
			name:               "disableChecksOnMethods",
			arguments:          lint.Arguments{"disablechecksonmethods"},
			wantErr:            nil,
			wantDisabledChecks: disabledChecks{PrivateReceivers: true, PublicInterfaces: true, Method: true},
			wantStuttersMsg:    "stutters",
		},
		{
			name:               "disableChecksOnTypes",
			arguments:          lint.Arguments{"disablechecksontypes"},
			wantErr:            nil,
			wantDisabledChecks: disabledChecks{PrivateReceivers: true, PublicInterfaces: true, Type: true},
			wantStuttersMsg:    "stutters",
		},
		{
			name:               "disableChecksOnVariables",
			arguments:          lint.Arguments{"disablechecksonvariables"},
			wantErr:            nil,
			wantDisabledChecks: disabledChecks{PrivateReceivers: true, PublicInterfaces: true, Var: true},
			wantStuttersMsg:    "stutters",
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
			var rule ExportedRule

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
			if rule.disabledChecks != tt.wantDisabledChecks {
				t.Errorf("unexpected disabledChecks: got = %v, want %v", rule.disabledChecks, tt.wantDisabledChecks)
			}
			if rule.stuttersMsg != tt.wantStuttersMsg {
				t.Errorf("unexpected stuttersMsg: got = %v, want %v", rule.stuttersMsg, tt.wantStuttersMsg)
			}
		})
	}
}
