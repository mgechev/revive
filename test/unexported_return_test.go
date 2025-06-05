package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestUnexportedReturn(t *testing.T) {
	testRule(t, "unexported_return_package_foo", &rule.UnexportedReturnRule{})
	testRule(t, "unexported_return_package_foo_test", &rule.UnexportedReturnRule{})
	testRule(t, "unexported_return_package_footest_test", &rule.UnexportedReturnRule{})
	testRule(t, "unexported_return_package_main", &rule.UnexportedReturnRule{})
	testRule(t, "unexported_return_package_main_test", &rule.UnexportedReturnRule{})
	testRule(t, "unexported_return_package_maintest_test", &rule.UnexportedReturnRule{})
}
