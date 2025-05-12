package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestUnexportedReturn(t *testing.T) {
	testRule(t, "unexported_return", &rule.UnexportedReturnRule{})
	testRule(t, "unexported_return_test", &rule.UnexportedReturnRule{})
}
