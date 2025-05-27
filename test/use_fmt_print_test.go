package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestUseFmtPrint(t *testing.T) {
	testRule(t, "use_fmt_print", &rule.UseFmtPrintRule{})
}
