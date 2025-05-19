package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestErrorStringsRule_Configure(t *testing.T) {
	tests := []struct {
		name      string
		arguments lint.Arguments
		wantErr   error
	}{
		{
			name:      "Default configuration",
			arguments: lint.Arguments{},
		},
		{
			name:      "Valid custom functions",
			arguments: lint.Arguments{"mypkg.MyErrorFunc", "errors.New"},
		},
		{
			name:      "Argument not a string",
			arguments: lint.Arguments{123},
		},
		{
			name:      "Invalid package",
			arguments: lint.Arguments{".MyErrorFunc"},
			wantErr:   errors.New("found invalid custom function: .MyErrorFunc"),
		},
		{
			name:      "Invalid function",
			arguments: lint.Arguments{"errors."},
			//revive:disable-next-line // error-strings
			wantErr: errors.New("found invalid custom function: errors."), //nolint:revive // error-strings: it's ok for tests
		},
		{
			name:      "Invalid custom function",
			arguments: lint.Arguments{"invalidFunction"},
			wantErr:   errors.New("found invalid custom function: invalidFunction"),
		},
		{
			name:      "Mixed valid and invalid custom functions",
			arguments: lint.Arguments{"mypkg.MyErrorFunc", "invalidFunction", "invalidFunction2"},
			wantErr:   errors.New("found invalid custom function: invalidFunction,invalidFunction2"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.ErrorStringsRule

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
