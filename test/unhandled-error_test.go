package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestUnhandledError(t *testing.T) {
	testRule(t, "unhandled-error", &rule.UnhandledErrorRule{})
}
