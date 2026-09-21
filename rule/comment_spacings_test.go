package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestCommentSpacingsRule_Configure(t *testing.T) {
	tests := []struct {
		name      string
		arguments lint.Arguments
		wantErr   error
	}{
		{
			name:      "no arguments uses default allow list",
			arguments: lint.Arguments{},
		},
		{
			name:      "valid arguments appended to default allow list",
			arguments: lint.Arguments{"mypragma:", "+optional"},
		},
		{
			name:      "invalid argument type",
			arguments: lint.Arguments{123},
			wantErr:   errors.New("invalid argument 123 for comment-spacings; expected string but got int"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.CommentSpacingsRule

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
