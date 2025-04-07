package rule

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestExportedRule_Configure(t *testing.T) {
	tests := []struct {
		name                string
		arguments           lint.Arguments
		wantErr             error
		wantDisabledChecks  disabledChecks
		wantIsRepetitiveMsg string
	}{
		{
			name:      "default configuration",
			arguments: lint.Arguments{},
			wantErr:   nil,
			wantDisabledChecks: disabledChecks{
				PrivateReceivers: true,
				PublicInterfaces: true,
			},
			wantIsRepetitiveMsg: "stutters",
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
			wantErr: nil,
			wantDisabledChecks: disabledChecks{
				PrivateReceivers: false,
				PublicInterfaces: false,
				Const:            true,
				Function:         true,
				Method:           true,
				RepetitiveNames:  true,
				Type:             true,
				Var:              true,
			},
			wantIsRepetitiveMsg: "stutters",
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
			wantErr: nil,
			wantDisabledChecks: disabledChecks{
				PrivateReceivers: false,
				PublicInterfaces: false,
				Const:            true,
				Function:         true,
				Method:           true,
				RepetitiveNames:  true,
				Type:             true,
				Var:              true,
			},
			wantIsRepetitiveMsg: "stutters",
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
			wantErr: nil,
			wantDisabledChecks: disabledChecks{
				PrivateReceivers: false,
				PublicInterfaces: false,
				Const:            true,
				Function:         true,
				Method:           true,
				RepetitiveNames:  true,
				Type:             true,
				Var:              true,
			},
			wantIsRepetitiveMsg: "stutters",
		},
		{
			name: "valid sayRepetitiveInsteadOfStutters",
			arguments: lint.Arguments{
				"sayRepetitiveInsteadOfStutters",
			},
			wantErr: nil,
			wantDisabledChecks: disabledChecks{
				PrivateReceivers: true,
				PublicInterfaces: true,
			},
			wantIsRepetitiveMsg: "is repetitive",
		},
		{
			name: "valid lowercased sayRepetitiveInsteadOfStutters",
			arguments: lint.Arguments{
				"sayrepetitiveinsteadofstutters",
			},
			wantErr: nil,
			wantDisabledChecks: disabledChecks{
				PrivateReceivers: true,
				PublicInterfaces: true,
			},
			wantIsRepetitiveMsg: "is repetitive",
		},
		{
			name: "valid kebab-cased sayRepetitiveInsteadOfStutters",
			arguments: lint.Arguments{
				"say-repetitive-instead-of-stutters",
			},
			wantErr: nil,
			wantDisabledChecks: disabledChecks{
				PrivateReceivers: true,
				PublicInterfaces: true,
			},
			wantIsRepetitiveMsg: "is repetitive",
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
				t.Errorf("unexpected disabledChecks: got = %+v, want %+v", rule.disabledChecks, tt.wantDisabledChecks)
			}
			if rule.isRepetitiveMsg != tt.wantIsRepetitiveMsg {
				t.Errorf("unexpected stuttersMsg: got = %v, want %v", rule.isRepetitiveMsg, tt.wantIsRepetitiveMsg)
			}
		})
	}
}
