package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestIfReturn(t *testing.T) {
	testRule(t, "if_return", &rule.IfReturnRule{})
}
