package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestFileLengthLimit(t *testing.T) {
	testRule(t, "file-length-limit-default", &rule.FileLengthLimitRule{}, &lint.RuleConfig{
		Arguments: []any{},
	})
	testRule(t, "file-length-limit-9", &rule.FileLengthLimitRule{}, &lint.RuleConfig{
		Arguments: []any{map[string]any{"max": int64(9)}},
	})
	testRule(t, "file-length-limit-7-skip-comments", &rule.FileLengthLimitRule{}, &lint.RuleConfig{
		Arguments: []any{map[string]any{"max": int64(7), "skipComments": true}},
	})
	testRule(t, "file-length-limit-6-skip-blank", &rule.FileLengthLimitRule{}, &lint.RuleConfig{
		Arguments: []any{map[string]any{"max": int64(6), "skipBlankLines": true}},
	})
	testRule(t, "file-length-limit-4-skip-comments-skip-blank", &rule.FileLengthLimitRule{}, &lint.RuleConfig{
		Arguments: []any{map[string]any{"max": int64(4), "skipComments": true, "skipBlankLines": true}},
	})
}
