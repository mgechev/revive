package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestRedundantTestMainExit(t *testing.T) {
	testRule(t, "go1.15/redundant_test_main_exit_test", &rule.RedundantTestMainExitRule{})
	testRule(t, "go1.15/redundant_test_main_exit", &rule.RedundantTestMainExitRule{})
	testRule(t, "redundant_test_main_exit_test", &rule.RedundantTestMainExitRule{})
}
