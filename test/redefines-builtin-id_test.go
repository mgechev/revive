package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

// Tests RedefinesBuiltinID rule.
func TestRedefinesBuiltinID(t *testing.T) {
	testRule(t, "redefines-builtin-id", &rule.RedefinesBuiltinIDRule{})
}

func TestRedefinesBuiltinIDAfterGo1_21(t *testing.T) {
	testRule(t, "go1.21/redefines-builtin-id", &rule.RedefinesBuiltinIDRule{})
}
