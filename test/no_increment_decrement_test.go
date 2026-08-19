package test_test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestNoIncrementDecrement(t *testing.T) {
	testRule(t, "no_increment_decrement", &rule.NoIncrementDecrementRule{})
}
