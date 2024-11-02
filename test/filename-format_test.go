package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestLintFilenameFormat(t *testing.T) {
	testRule(t, "filename-ok-default", &rule.FilenameFormatRule{}, &lint.RuleConfig{})

	testRule(t, "filenamе-with-non-ascii-char", &rule.FilenameFormatRule{}, &lint.RuleConfig{})

	testRule(t, "filename_with_underscores", &rule.FilenameFormatRule{}, &lint.RuleConfig{Arguments: []any{"^[A-Za-z][A-Za-z0-9]*.go$"}})
}
