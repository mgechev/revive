package test_test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestUseErrorsNew(t *testing.T) {
	testRule(t, "use_errors_new", &rule.UseErrorsNewRule{})
	testRule(t, "go1.26/use_errors_new", &rule.UseErrorsNewRule{})
}
