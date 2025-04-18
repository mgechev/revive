package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestFileLengthLimit(t *testing.T) {
	testRule(t, "file_length_limit_disabled", &rule.FileLengthLimitRule{}, &lint.RuleConfig{
		Arguments: []any{},
	})
	testRule(t, "file_length_limit_disabled", &rule.FileLengthLimitRule{}, &lint.RuleConfig{
		Arguments: []any{map[string]any{"max": int64(0)}},
	})
	testRule(t, "file_length_limit_disabled", &rule.FileLengthLimitRule{}, &lint.RuleConfig{
		Arguments: []any{map[string]any{"skipComments": true, "skipBlankLines": true}},
	})
	testRule(t, "file_length_limit_9", &rule.FileLengthLimitRule{}, &lint.RuleConfig{
		Arguments: []any{map[string]any{"max": int64(9)}},
	})
	testRule(t, "file_length_limit_7_skip_comments", &rule.FileLengthLimitRule{}, &lint.RuleConfig{
		Arguments: []any{map[string]any{"max": int64(7), "skipComments": true}},
	})
	testRule(t, "file_length_limit_6_skip_blank", &rule.FileLengthLimitRule{}, &lint.RuleConfig{
		Arguments: []any{map[string]any{"max": int64(6), "skipBlankLines": true}},
	})
	testRule(t, "file_length_limit_4_skip_comments_skip_blank", &rule.FileLengthLimitRule{}, &lint.RuleConfig{
		Arguments: []any{map[string]any{"max": int64(4), "skipComments": true, "skipBlankLines": true}},
	})
	testRule(t, "file_length_limit_4_skip_comments_skip_blank", &rule.FileLengthLimitRule{}, &lint.RuleConfig{
		Arguments: []any{map[string]any{"max": int64(4), "skip-comments": true, "skip-blank-lines": true}},
	})
}
