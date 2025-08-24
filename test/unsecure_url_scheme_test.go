package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestUnsecureURLScheme(t *testing.T) {
	testRule(t, "unsecure_url_scheme", &rule.UnsecureURLSchemeRule{}, &lint.RuleConfig{})
}
