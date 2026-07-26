package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestImportsBlocklistOriginal(t *testing.T) {
	testRule(t, "imports_blocklist_original", &rule.ImportsBlocklistRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"crypto/md5", "crypto/sha1"},
	})
}

func TestImportsBlocklist(t *testing.T) {
	testRule(t, "imports_blocklist", &rule.ImportsBlocklistRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"github.com/full/match", "wildcard/**/between", "wildcard/backward/**", "**/wildcard/forward", "full"},
	})
}

func BenchmarkImportsBlocklist(b *testing.B) {
	for b.Loop() {
		testRule(b, "imports_blocklist", &rule.ImportsBlocklistRule{}, &lint.RuleConfig{
			Arguments: lint.Arguments{"github.com/full/match", "wildcard/**/between", "wildcard/backward/**", "**/wildcard/forward", "full"},
		})
	}
}
