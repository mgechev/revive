package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestEnforceElse(t *testing.T) {
	testRule(t, "enforce_else", &rule.EnforceElseRule{})
}
