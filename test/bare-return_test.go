package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/rule"
)

func TestBareReturn(t *testing.T) {
	testRule(t, "bare-return", &rule.BareReturnRule{})
}
