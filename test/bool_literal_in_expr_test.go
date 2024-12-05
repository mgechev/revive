package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestBoolLiteral(t *testing.T) {
	testRule(t, "bool_literal_in_expr", &rule.BoolLiteralRule{})
}
