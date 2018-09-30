package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestWaitGroupByCopy(t *testing.T) {
	testRule(t, "waitgroup-by-copy", &rule.WaitGroupByCopyRule{})
}
