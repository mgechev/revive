package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestImportsBlocklistOriginal(t *testing.T) {
	args := []any{"crypto/md5", "crypto/sha1"}

	testRule(t, "imports-blocklist-original", &rule.ImportsBlocklistRule{}, &lint.RuleConfig{
		Arguments: args,
	})
}

func TestImportsBlocklist(t *testing.T) {
	args := []any{"github.com/full/match", "wildcard/**/between", "wildcard/backward/**", "**/wildcard/forward", "full"}

	testRule(t, "imports-blocklist", &rule.ImportsBlocklistRule{}, &lint.RuleConfig{
		Arguments: args,
	})
}

func BenchmarkImportsBlocklist(b *testing.B) {
	args := []any{"github.com/full/match", "wildcard/**/between", "wildcard/backward/**", "**/wildcard/forward", "full"}
	var t *testing.T
	for i := 0; i <= b.N; i++ {
		testRule(t, "imports-blocklist", &rule.ImportsBlocklistRule{}, &lint.RuleConfig{
			Arguments: args,
		})
	}
}
