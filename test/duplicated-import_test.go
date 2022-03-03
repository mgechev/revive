package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/rule"
)

func TestDuplicatedImports(t *testing.T) {
	testRule(t, "duplicated-imports", &rule.DuplicatedImportsRule{})
}
