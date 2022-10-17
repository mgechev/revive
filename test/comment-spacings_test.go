package test

import(
	"testing"
	
	"github.com/mgechev/revive/rule"
)

func TestCommentSpacings(t *testing.T) {
	testRule(t, "comment-spacings", &rule.CommentSpacingsRule{})
}