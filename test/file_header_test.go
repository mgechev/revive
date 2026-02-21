package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestLintFileHeaderDefault(t *testing.T) {
	testRule(t, "lint_file_header_default", &rule.FileHeaderRule{}, &lint.RuleConfig{})
}

func TestLintFileHeader(t *testing.T) {
	testRule(t, "lint_file_header1", &rule.FileHeaderRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"foobar"},
	})

	testRule(t, "lint_file_header2", &rule.FileHeaderRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"foobar"},
	})

	testRule(t, "lint_file_header3", &rule.FileHeaderRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"foobar"},
	})

	testRule(t, "lint_file_header4", &rule.FileHeaderRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"^\\sfoobar$"},
	})

	testRule(t, "lint_file_header5", &rule.FileHeaderRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"^\\sfoo.*bar$"},
	})

	testRule(t, "lint_file_header6", &rule.FileHeaderRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"foobar"},
	})
}

func BenchmarkLintFileHeader(b *testing.B) {
	var t *testing.T
	for i := 0; i <= b.N; i++ {
		testRule(t, "lint_file_header1", &rule.FileHeaderRule{}, &lint.RuleConfig{
			Arguments: lint.Arguments{"foobar"},
		})
	}
}
