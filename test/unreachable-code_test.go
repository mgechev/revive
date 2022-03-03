package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/rule"
)

func TestUnreachableCode(t *testing.T) {
	testRule(t, "unreachable-code", &rule.UnreachableCodeRule{})
}
