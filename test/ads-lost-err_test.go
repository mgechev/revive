package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestADSLostErr(t *testing.T) {
	testRule(t, "ads-lost-err", &rule.ADSLostErrRule{})
}
