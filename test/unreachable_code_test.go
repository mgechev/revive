package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestUnreachableCode(t *testing.T) {
	testRule(t, "unreachable_code", &rule.UnreachableCodeRule{})
}
