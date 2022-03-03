package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/lint"
	"github.com/deepsourcelabs/revive/rule"
)

func TestImportsBlacklist(t *testing.T) {
	args := []interface{}{"crypto/md5", "crypto/sha1"}

	testRule(t, "imports-blacklist", &rule.ImportsBlacklistRule{}, &lint.RuleConfig{
		Arguments: args,
	})
}

func BenchmarkImportsBlacklist(b *testing.B) {
	args := []interface{}{"crypto/md5", "crypto/sha1"}
	var t *testing.T
	for i := 0; i <= b.N; i++ {
		testRule(t, "imports-blacklist", &rule.ImportsBlacklistRule{}, &lint.RuleConfig{
			Arguments: args,
		})
	}
}
