package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestLintFileHeader(t *testing.T) {
	testRule(t, "lint-file-header1", &rule.FileHeaderRule{}, &lint.RuleConfig{
		Arguments: []any{"foobar"},
	})

	testRule(t, "lint-file-header2", &rule.FileHeaderRule{}, &lint.RuleConfig{
		Arguments: []any{"foobar"},
	})

	testRule(t, "lint-file-header3", &rule.FileHeaderRule{}, &lint.RuleConfig{
		Arguments: []any{"foobar"},
	})

	testRule(t, "lint-file-header4", &rule.FileHeaderRule{}, &lint.RuleConfig{
		Arguments: []any{"^\\sfoobar$"},
	})

	testRule(t, "lint-file-header5", &rule.FileHeaderRule{}, &lint.RuleConfig{
		Arguments: []any{"^\\sfoo.*bar$"},
	})

	testRule(t, "lint-file-header6", &rule.FileHeaderRule{}, &lint.RuleConfig{
		Arguments: []any{"foobar"},
	})
}

func BenchmarkLintFileHeader(b *testing.B) {
	var t *testing.T
	for i := 0; i <= b.N; i++ {
		testRule(t, "lint-file-header1", &rule.FileHeaderRule{}, &lint.RuleConfig{
			Arguments: []any{"foobar"},
		})
	}
}
