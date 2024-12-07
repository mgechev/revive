package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestAtomic(t *testing.T) {
	testRule(t, "atomic", &rule.AtomicRule{})
}
