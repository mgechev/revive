package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestIdenticalSwitchConditions(t *testing.T) {
	testRule(t, "identical_switch_conditions", &rule.IdenticalSwitchConditionsRule{})
}
