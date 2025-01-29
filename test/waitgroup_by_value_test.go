package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestWaitGroupByValue(t *testing.T) {
	testRule(t, "waitgroup_by_value", &rule.WaitGroupByValueRule{})
}
