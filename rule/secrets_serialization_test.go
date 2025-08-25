package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestSecretsSerializationConfigure(t *testing.T) {
	type argumentsTest struct {
		name      string
		arguments lint.Arguments
		wantErr   error
	}
	tests := []argumentsTest{
		{
			name: "String Slice",
			arguments: lint.Arguments{
				[]any{"foo", "Bar"}},
			wantErr: nil,
		},
		{
			name: "Int Slice",
			arguments: lint.Arguments{
				[]any{1, 2}},
			wantErr: errors.New("invalid argument to the secrets-serialization rule: expecting secretFieldIndicators of type slice of strings, got slice of type int"),
		},
		{
			name: "Not a Slice",
			arguments: lint.Arguments{
				"this is not a slice"},
			wantErr: errors.New("invalid argument to the secrets-serialization rule: expecting secretFieldIndicators of type slice of strings, got string"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.SecretsSerializationRule
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
