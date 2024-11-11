package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestBareReturn(t *testing.T) {
	testRule(t, "bare_return", &rule.BareReturnRule{})
}
