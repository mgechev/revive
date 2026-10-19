package test_test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestUseSlicesConcat(t *testing.T) {
	testRule(t, "use_slices_concat", &rule.UseSlicesConcatRule{})
	testRule(t, "go1.22/use_slices_concat", &rule.UseSlicesConcatRule{})
}
