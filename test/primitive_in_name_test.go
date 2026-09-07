package test_test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestPrimitiveInName(t *testing.T) {
	testRule(t, "primitive_in_name", &rule.PrimitiveInNameRule{})
}
