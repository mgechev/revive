package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestRedundantTestMainExit(t *testing.T) {
	testRule(t, "redundant_test_main_exit_test", &rule.RedundantTestMainExitRule{})
	testRule(t, "redundant_test_main_exit", &rule.RedundantTestMainExitRule{})
}
