package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/rule"
)

func TestWaitGroupByValue(t *testing.T) {
	testRule(t, "waitgroup-by-value", &rule.WaitGroupByValueRule{})
}
