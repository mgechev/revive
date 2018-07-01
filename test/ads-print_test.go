package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestADSPrint(t *testing.T) {
	testRule(t, "ads-print", &rule.ADSPrintRule{})
}
