package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestUseErrorsNew(t *testing.T) {
	testRule(t, "use_errors_new", &rule.UseErrorsNewRule{})
}
