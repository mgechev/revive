package rule

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestCommentSpacingsRule_Configure(t *testing.T) {
	tests := []struct {
		name          string
		arguments     lint.Arguments
		wantErr       error
		wantAllowList []string
	}{
		{
			name:      "no arguments uses default allow list",
			arguments: lint.Arguments{},
			wantErr:   nil,
			wantAllowList: []string{
				"//#nosec",
			},
		},
		{
			name:      "valid arguments appended to default allow list",
			arguments: lint.Arguments{"mypragma:", "+optional"},
			wantErr:   nil,
			wantAllowList: []string{
				"//#nosec",
				"//mypragma:",
				"//+optional",
			},
		},
		{
			name:      "invalid argument type",
			arguments: lint.Arguments{123},
			wantErr:   errors.New("invalid argument 123 for comment-spacings; expected string but got int"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r CommentSpacingsRule

			err := r.Configure(tt.arguments)

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
				return
			}
			if len(r.allowList) != len(tt.wantAllowList) {
				t.Errorf("unexpected allowList length: got = %d, want = %d", len(r.allowList), len(tt.wantAllowList))
				return
			}
			for i, entry := range tt.wantAllowList {
				if r.allowList[i] != entry {
					t.Errorf("unexpected allowList[%d]: got = %q, want = %q", i, r.allowList[i], entry)
				}
			}
		})
	}
}
