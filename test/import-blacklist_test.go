package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestImportsBlacklistOriginal(t *testing.T) {
	args := []interface{}{"crypto/md5", "crypto/sha1"}

	testRule(t, "imports-blacklist-original", &rule.ImportsBlacklistRule{}, &lint.RuleConfig{
		Arguments: args,
	})
}

func TestImportsBlacklist(t *testing.T) {
	args := []interface{}{"github.com/full/match", "wildcard/**/between", "wildcard/backward/**", "**/wildcard/forward", "full"}

	testRule(t, "imports-blacklist", &rule.ImportsBlacklistRule{}, &lint.RuleConfig{
		Arguments: args,
	})
}

func BenchmarkImportsBlacklist(b *testing.B) {
	args := []interface{}{"github.com/full/match", "wildcard/**/between", "wildcard/backward/**", "**/wildcard/forward", "full"}
	var t *testing.T
	for i := 0; i <= b.N; i++ {
		testRule(t, "imports-blacklist", &rule.ImportsBlacklistRule{}, &lint.RuleConfig{
			Arguments: args,
		})
	}
}
